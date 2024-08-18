package server

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/heyjun3/dforget/backend/gen/api/memo/v1/memov1connect"
)

func New(conf Config) *http.ServeMux {
	mux := http.NewServeMux()
	db := InitDBConn(conf)
	memo := initializeMemoHandler(db)
	path, handler := memov1connect.NewMemoServiceHandler(memo)
	mux.Handle(path, handler)
	mux.HandleFunc("GET /oidc", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request to oidc")
		log.Println(r.URL.Query())
		code := r.URL.Query().Get("code")
		if code == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		formData := url.Values{
			"grant_type": {"authorization_code"},
			"code": {code},
			"redirect_uri": {"http://localhost:8080/oidc"},
			"client_id": {"myclient"},
			"client_secret": {"aOtFO4E9FmCB3UXq47GNg5OH8ABYM3MR"},
		}
		req, err := http.NewRequest("POST", "http://keycloak:8080/realms/myrealm/protocol/openid-connect/token", strings.NewReader(formData.Encode()))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print(err)
			return
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		log.Println(string(body))
		w.WriteHeader(http.StatusOK)
	})
	return mux
}
