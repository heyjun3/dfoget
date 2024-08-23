package server

import (
	"log/slog"
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
		ctx := r.Context()
		slog.InfoContext(ctx, "recieve oidc redirect")
		code := r.URL.Query().Get("code")
		if code == "" {
			slog.ErrorContext(ctx, "code is empty")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		formData := url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {code},
			"redirect_uri":  {conf.oidc.redirectUri},
			"client_id":     {conf.oidc.clientId},
			"client_secret": {conf.oidc.clientSecret},
		}
		req, err := http.NewRequest("POST", conf.oidc.tokenUrl, strings.NewReader(formData.Encode()))
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		slog.InfoContext(ctx, "oidc verified")
		http.Redirect(w, r, "http://localhost:3000", http.StatusTemporaryRedirect)
	})
	return mux
}
