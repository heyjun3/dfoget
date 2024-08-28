package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	CheckJWTVerify()
}

func CheckJWTVerify() {
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJkek1fbDRBM2hZZ1dpcDRUNW93M210QlNhVEhVelBPOXVwVjJKYWwzdUZBIn0.eyJleHAiOjE3MjQ1NzAwNzcsImlhdCI6MTcyNDU2OTc3NywiYXV0aF90aW1lIjoxNzI0NTY5Nzc3LCJqdGkiOiJjNjdjYTUwYS0xYjNjLTQxYTEtOWE1YS0zMjJjYWFkYThlZmYiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0Ojg4ODgvcmVhbG1zL215cmVhbG0iLCJhdWQiOiJteWNsaWVudCIsInN1YiI6IjhkYmEwZDc1LWY3YmUtNDgyYS05MjY3LTZlZDZkMGE5NGNhNCIsInR5cCI6IklEIiwiYXpwIjoibXljbGllbnQiLCJzaWQiOiIwMTFhYjlhYy0zYzA2LTQwM2QtOTg2My00ZTRhZjE3MmRiNjIiLCJhdF9oYXNoIjoiN2o5VDdZQ1FiaVZjdEJ6VFYwb3BtUSIsImFjciI6IjEiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsIm5hbWUiOiJoZXkganVuIiwicHJlZmVycmVkX3VzZXJuYW1lIjoibXl1c2VyIiwiZ2l2ZW5fbmFtZSI6ImhleSIsImZhbWlseV9uYW1lIjoianVuIiwiZW1haWwiOiJteXVzZXJAZ21haS5jb20ifQ.k4VC1r_kwDEwSEdGJWwMNti_PzMtH3_EHDIsXU7AasA8tk-5ElybRs3fwr7rpXNjmCRO_X9_iX2Dv-AfzCSZATkUV4wpwT7h5bYYrf9V1deLOGHcs8FEnFnIx6JfdVlxkUhQ1ednJMx7-usDjhU_U8FIh4vWo77b6_3CishpYppqsRxgxBrHuJrxY9E_bm6bNsxRyIFDS0K5ixyNUk-kws5P1GItJNvqqVAwWbaO9XPseEQ48tugw-rIy3Abno4nxYN3628RANmZBW3DBXGqNv4rvElwMA7aGcgCVvovLQZhQqXe7m5dqXjViKQSOCZ3fAWghlazTm70wIg72juuWg"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signin method: %s", token.Header["alg"])
		}
		k := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnyKst0eKVrUHWQ+b0aci9TAC6aPP6LYxNeICknWiQrozx0KWnX+Bvqm75q/SM3D1WbouMyYAy2JaAcZDqhLb5z4Nx3ZqzHzAx0QlCWJ0pYkYolyXsTVgKrp1gxHYHc2jKg+UOcNM624QK2ApBrBN4IK80Vi0dgdgbSLo3tBWpm/ZTCj7j468lSlZs+JjBFP8na8NFsZahd6hE+V37foPYBZxODeMBemnkEr6eEZ5EJK0gYlD/4bdoK52u0jGLPITLtlwOiJasTG7rVjczkoePylMANk3mHjZV/lL+raPQMGrGdvHNiOYXTlNAe1J7aELwrQIgL7GK290iCdrT0fIaQIDAQAB"
		var d []byte
		d = make([]byte, base64.StdEncoding.EncodedLen(len([]byte(k))))
		n, _ := base64.StdEncoding.Decode(d, []byte(k))
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
		fmt.Println(claims)
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
