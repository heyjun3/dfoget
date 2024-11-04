package main

import (
	// "bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"log/slog"
	"slices"
	"strings"

	"github.com/cockroachdb/swiss"
	"github.com/golang-jwt/jwt/v5"
)

func fn(str []string, ns []string) ([]string, []string) {
	if len(str) == 0 {
		return str, ns
	}
	if str[0] != "AAA" {
		ns = append(ns, str[0])
	}
	nstr := slices.Delete(str, 0, 1)

	return fn(nstr, ns)
}
func main() {
	s := []string{"AAA", "BBB", "CCC"}

	s1, s2 := fn(s, []string{})
	fmt.Println("s1", s1)
	fmt.Println("s2", s2)
}

func SwissTable() {
	m := swiss.New[string, string](10)
	m.Put("test", "test")
	v, _ := m.Get("test")
	print(v, "\n")
}

func JwtVerify() {
	private, public, err := GenerateRSAKey()
	if err != nil {
		panic(err)
	}
	token := GenerateJWT(private)
	VerifyJWT(token, string(public))
}

func GenerateJWT(privateKey *rsa.PrivateKey) string {
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

func GenerateRSAKey() (privateKey *rsa.PrivateKey, publicKey []byte, err error) {
	bitSize := 4096
	privateKey, err = GeneratePrivateKey(bitSize)
	if err != nil {
		panic(err)
	}
	publicKey, err = GeneratePublicKey(&privateKey.PublicKey)
	if err != nil {
		panic(err)
	}
	return privateKey, publicKey, err
}

func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}
	slog.Info("generate private key")
	return privateKey, nil
}

func GeneratePublicKey(pubKey *rsa.PublicKey) ([]byte, error) {
	pub2, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		panic(err)
	}
	// var b bytes.Buffer

	// b64 := base64.NewEncoder(base64.StdEncoding, &b)
	// if _, err := b64.Write(pub2); err != nil {
	// 	panic(err)
	// }
	// defer b64.Close()
	// slog.Info("generate public key")
	// slog.Info(b.String())
	// return b.Bytes(), nil

	pubBlock := pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   pub2,
	}

	pubPem := pem.EncodeToMemory(&pubBlock)

	// public, err := ssh.NewPublicKey(pubKey)
	// if err != nil {
	// 	return nil, err
	// }
	// pubKeyBytes := ssh.MarshalAuthorizedKey(public)

	slog.Info("generate public key")
	ks := strings.Split(string(pubPem), "\n")
	pubArr := []string{}
	for _, k := range ks {
		if strings.Contains(k, "-----") {
			continue
		}
		pubArr = append(pubArr, k)
	}
	pub := strings.Join(pubArr, "")

	slog.Info(pub)
	return []byte(pub), nil
}

func CheckVerifyJWT() {
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJkek1fbDRBM2hZZ1dpcDRUNW93M210QlNhVEhVelBPOXVwVjJKYWwzdUZBIn0.eyJleHAiOjE3Mjc4NzY4MDIsImlhdCI6MTcyNzg3NjUwMiwiYXV0aF90aW1lIjoxNzI3ODc2Mzc4LCJqdGkiOiJjNmNmOTlhZC1jMzFmLTQ5NzQtYWJlOC01MTQ3N2JhYmI5NzMiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0Ojg4ODgvcmVhbG1zL215cmVhbG0iLCJhdWQiOiJteWNsaWVudCIsInN1YiI6IjhkYmEwZDc1LWY3YmUtNDgyYS05MjY3LTZlZDZkMGE5NGNhNCIsInR5cCI6IklEIiwiYXpwIjoibXljbGllbnQiLCJzaWQiOiI3Y2Y0ZjQzNy1mYzY1LTRhNjktOTZiNC1mYjhhOTAzNDg4NjQiLCJhdF9oYXNoIjoiLXIwUm04bXpub0s2Qlh0bDdqM1hQUSIsImFjciI6IjAiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsIm5hbWUiOiJoZXkganVuIiwicHJlZmVycmVkX3VzZXJuYW1lIjoibXl1c2VyIiwiZ2l2ZW5fbmFtZSI6ImhleSIsImZhbWlseV9uYW1lIjoianVuIiwiZW1haWwiOiJteXVzZXJAZ21haS5jb20ifQ.AKb5MLk_3AWAG75Ft7rekAHNz4emiNKaE4NTUSPVEjd6cA3xgS_0L3pK-2HAOtRhubsRvSOknZ9aYPSCbC16_ZZIJAUOTlpMvzLygg5fixdERflHdJdmzRfmNHJqbTd3GSUBH5pM1D_4n7dMcI3Q6SEWxZn0SF2JL8_ru6XS_nWmiRTmlxaQVViRgcBf9OkFYyQK2es-okdxC30VLjJuP9bdniZN7MxDWa9qzLPgItlGgwOcAmSvvMrZs250LBUHzwRuY-9fIXyTuEo0wsSv0yP-Udq-HDSYF1iwNlA-r0gNO5nkfkZ-XOU8pu1VIuMozcDqeRy8bXtvhqSL7m77zQ"
	pubkey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnyKst0eKVrUHWQ+b0aci9TAC6aPP6LYxNeICknWiQrozx0KWnX+Bvqm75q/SM3D1WbouMyYAy2JaAcZDqhLb5z4Nx3ZqzHzAx0QlCWJ0pYkYolyXsTVgKrp1gxHYHc2jKg+UOcNM624QK2ApBrBN4IK80Vi0dgdgbSLo3tBWpm/ZTCj7j468lSlZs+JjBFP8na8NFsZahd6hE+V37foPYBZxODeMBemnkEr6eEZ5EJK0gYlD/4bdoK52u0jGLPITLtlwOiJasTG7rVjczkoePylMANk3mHjZV/lL+raPQMGrGdvHNiOYXTlNAe1J7aELwrQIgL7GK290iCdrT0fIaQIDAQAB"
	VerifyJWT(tokenString, pubkey)
}

func VerifyJWT(jwtString string, pubKey string) {

	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signin method: %s", token.Header["alg"])
		}
		// pkey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pubKey))
		// if err != nil {
		// 	panic(err)
		// }

		var d []byte
		d = make([]byte, base64.StdEncoding.EncodedLen(len([]byte(pubKey))))
		n, _ := base64.StdEncoding.Decode(d, []byte(pubKey))
		d = d[:n]

		var parsedKey interface{}
		var err error
		if parsedKey, err = x509.ParsePKIXPublicKey(d); err != nil {
			return nil, err
		}

		var pkey *rsa.PublicKey
		var ok bool
		if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
			return nil, jwt.ErrNotRSAPublicKey
		}

		return pkey, nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims.GetAudience())
		fmt.Println(claims.GetExpirationTime())
		fmt.Println(claims.GetIssuedAt())
		fmt.Println(claims.GetIssuer())
		fmt.Println(claims.GetNotBefore())
		fmt.Println(claims.GetSubject())
	} else {
		fmt.Println(err)
	}
}

type User struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Age      *int    `json:"age,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
	MidName  *string `json:"mid_name"`
}

func Ptr[T any](v T) *T {
	return &v
}

func CheckOmitEmpty() {
	user := User{
		ID:       "test",
		Name:     "test_name",
		Age:      Ptr(30),
		IsActive: Ptr(true),
	}
	buf, _ := json.Marshal(user)
	fmt.Println(string(buf))

	user = User{
		ID:   "test",
		Name: "test_name",
	}
	buf, _ = json.Marshal(user)
	fmt.Println(string(buf))

	user = User{
		ID:       "test",
		Name:     "test_name",
		Age:      nil,
		IsActive: nil,
	}
	buf, _ = json.Marshal(user)
	fmt.Println(string(buf))

}
