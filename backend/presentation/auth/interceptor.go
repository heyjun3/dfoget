package auth

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v5"

	"github.com/heyjun3/dforget/backend/lib"
	cfg "github.com/heyjun3/dforget/backend/config"
)

const (
	AuthCookieName = "dforget"
)

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
func verifyJWT(ctx context.Context, cookie *http.Cookie, pubkey *rsa.PublicKey) (*claim, error) {
	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signin method: %s", token.Header)
		}
		return pubkey, nil
	})
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), "msg", "jwt parse error")
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		slog.ErrorContext(ctx, "claims error")
		return nil, err
	}

	sub, err := claims.GetSubject()
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), "msg", "get subject error")
		return nil, err
	}
	return &claim{
		sub: sub,
	}, nil
}

type claim struct {
	sub string
}

type authInterceptor struct {
	publicKey *rsa.PublicKey
}

func NewAuthInterceptorV2(conf cfg.Config) *authInterceptor {
	publicKey, err := NewPublicKey(conf.Oidc.Pubkey)
	if err != nil {
		slog.ErrorContext(context.Background(), "failed load pub key error")
		panic(err)
	}
	return &authInterceptor{
		publicKey: publicKey,
	}
}

func (i *authInterceptor) WrapUnary(
	next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(
		ctx context.Context,
		req connect.AnyRequest,
	) (connect.AnyResponse, error) {
		if req.Spec().IsClient {
			return next(ctx, req)
		}
		r := http.Request{Header: req.Header()}
		cookie, err := r.Cookie(AuthCookieName)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				fmt.Errorf("expect cookie name is %s", AuthCookieName),
			)
		}
		claim, err := verifyJWT(ctx, cookie, i.publicKey)
		if err != nil {
			return nil, connect.NewError(
				connect.CodeUnauthenticated,
				fmt.Errorf("invalid access token"),
			)
		}
		ctx = lib.SetSubKey(ctx, claim.sub)
		return next(ctx, req)
	})
}

func (i *authInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return connect.StreamingClientFunc(func(
		ctx context.Context,
		spec connect.Spec,
	) connect.StreamingClientConn {
		return next(ctx, spec)
	})
}

func (i *authInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return connect.StreamingHandlerFunc(func(
		ctx context.Context,
		conn connect.StreamingHandlerConn,
	) error {
		header := conn.RequestHeader()
		req := http.Request{Header: header}
		cookie, err := req.Cookie(AuthCookieName)
		if err != nil {
			return connect.NewError(
				connect.CodeUnauthenticated,
				fmt.Errorf("expect cookie name is %s", AuthCookieName),
			)
		}
		claim, err := verifyJWT(ctx, cookie, i.publicKey)
		if err != nil {
			return connect.NewError(
				connect.CodeUnauthenticated,
				fmt.Errorf("invalid access token"),
			)
		}
		ctx = lib.SetSubKey(ctx, claim.sub)
		return next(ctx, conn)
	})
}

func NewAuthInterceptor(conf cfg.Config) connect.UnaryInterceptorFunc {
	publicKey, err := NewPublicKey(conf.oidc.pubkey)
	if err != nil {
		panic(err)
	}
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			r := http.Request{Header: req.Header()}
			cookie, err := r.Cookie(AuthCookieName)
			if err != nil {
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					fmt.Errorf("expect cookie name is %s", AuthCookieName),
				)
			}
			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signin method: %s", token.Header)
				}
				return publicKey, nil
			})
			if err != nil {
				slog.ErrorContext(ctx, err.Error(), "msg", "jwt parse error")
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
				slog.ErrorContext(ctx, err.Error(), "msg", "get subject error")
				return nil, connect.NewError(
					connect.CodeUnauthenticated,
					fmt.Errorf("token isn't valid"),
				)
			}
			ctx = lib.SetSubKey(ctx, sub)
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
