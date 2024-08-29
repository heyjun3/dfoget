package server

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v5"
)

var publicKey *rsa.PublicKey

func init() {
	var err error
	conf := NewConfig()
	publicKey, err = NewPublicKey(conf.oidc.pubkey)
	if err != nil {
		panic(err)
	}
}

const (
	AuthCookieName = "dforget"
)

type subKey struct{}

func SetSubKey(ctx context.Context, sub string) context.Context {
	return context.WithValue(ctx, subKey{}, sub)
}
func GetSubValue(ctx context.Context) (string, error) {
	val, ok := ctx.Value(subKey{}).(string)
	if !ok {
		return "", fmt.Errorf("no set subject value")
	}
	return val, nil
}

func NewPublicKey(pubkey string) (*rsa.PublicKey, error) {
	var bytes []byte
	bytes = make([]byte, base64.StdEncoding.EncodedLen(len([]byte(pubkey))))
	n, _ := base64.StdEncoding.Decode(bytes, []byte(pubkey))
	bytes = bytes[:n]

	var parsedKey interface{}
	var err error
	if parsedKey, err = x509.ParsePKIXPublicKey(bytes); err != nil {
		return nil, err
	}

	var pub *rsa.PublicKey
	var ok bool
	if pub, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, jwt.ErrNotRSAPublicKey
	}

	return pub, nil
}

func NewAuthInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			cookie := req.Header().Get("Cookie")
			name, val, _ := strings.Cut(cookie, "=")
			if name != AuthCookieName {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					fmt.Errorf("expect cookie name is %s", AuthCookieName),
				)
			}
			token, err := jwt.Parse(val, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signin method: %s", token.Header)
				}
				return publicKey, nil
			})
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					fmt.Errorf("token isn't valid"),
				)
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				slog.ErrorContext(ctx, "claims error")
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					fmt.Errorf("token isn't valid"),
				)
			}

			sub, err := claims.GetSubject()
			if err != nil {
				slog.ErrorContext(ctx, err.Error())
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					fmt.Errorf("token isn't valid"),
				)
			}
			ctx = SetSubKey(ctx, sub)
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
