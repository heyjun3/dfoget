package test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"github.com/golang-jwt/jwt/v5"
	"github.com/heyjun3/dforget/backend/server"
	"github.com/uptrace/bun"
)

var cookie string
var PublicKey string

func init() {
	var private *rsa.PrivateKey
	var err error
	private, PublicKey, err = generateRSAKey()
	if err != nil {
		panic(err)
	}
	cookie = generateJWT(private)
}

func ResetModel(db *bun.DB) {
	models := []interface{}{
		server.MemoDM{},
	}
	for _, model := range models {
		db.NewTruncateTable().Model(&model).Exec(context.Background())
	}
}

func generateRSAKey() (private *rsa.PrivateKey, public string, err error) {
	bitSize := 4096
	private, err = generatePrivateKey(bitSize)
	if err != nil {
		panic(err)
	}
	public, err = generatePublicKey(&private.PublicKey)
	if err != nil {
		panic(err)
	}
	return private, public, err
}
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}
func generatePublicKey(pubKey *rsa.PublicKey) (string, error) {
	pub2, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		panic(err)
	}

	pubBlock := pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   pub2,
	}

	pubPem := pem.EncodeToMemory(&pubBlock)

	ks := strings.Split(string(pubPem), "\n")
	pubArr := []string{}
	for _, k := range ks {
		if strings.Contains(k, "-----") {
			continue
		}
		pubArr = append(pubArr, k)
	}
	pub := strings.Join(pubArr, "")

	return pub, nil
}

func generateJWT(privateKey *rsa.PrivateKey) string {
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"exp":                2724570077,
		"iat":                1724569777,
		"auth_time":          1724569777,
		"jti":                "c67ca50a-1b3c-41a1-9a5a-322caada8eff",
		"iss":                "http://localhost:8888/realms/myrealm",
		"aud":                "myclient",
		"sub":                "8dba0d75-f7be-482a-9267-6ed6d0a94ca4",
		"typ":                "ID",
		"azp":                "myclient",
		"sid":                "011ab9ac-3c06-403d-9863-4e4af172db62",
		"at_hash":            "7j9T7YCQbiVctBzTV0opmQ",
		"acr":                "1",
		"email_verified":     false,
		"name":               "hey jun",
		"preferred_username": "myuser",
		"given_name":         "hey",
		"family_name":        "jun",
		"email":              "myuser@gmai.com",
	})
	s, err := t.SignedString(privateKey)
	if err != nil {
		panic(err)
	}
	return s
}

func NewSetCookieInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			if req.Spec().IsClient {
				req.Header().Set("Cookie", fmt.Sprintf("%s=%s", server.AuthCookieName, cookie))
			}
			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
