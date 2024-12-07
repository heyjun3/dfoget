package auth

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	cfg "github.com/heyjun3/dforget/backend/config"
)

type OIDCHandler struct {
	conf       cfg.Config
	httpClient HttpClient
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewOIDCHandler(conf cfg.Config, client HttpClient) *OIDCHandler {
	return &OIDCHandler{
		conf:       conf,
		httpClient: client,
	}
}

func (h OIDCHandler) RecieveRedirect(w http.ResponseWriter, r *http.Request) {
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
		"redirect_uri":  {h.conf.OIDC.RedirectUri},
		"client_id":     {h.conf.OIDC.ClientId},
		"client_secret": {h.conf.OIDC.ClientSecret},
	}
	req, err := http.NewRequest("POST", h.conf.OIDC.TokenUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err := h.httpClient.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var token OIDCToken
	if err := json.Unmarshal(buf, &token); err != nil {
		slog.ErrorContext(ctx, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{
		Name:     AuthCookieName,
		Value:    token.IdToken,
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	slog.InfoContext(ctx, "oidc verified")
	http.Redirect(w, r, h.conf.FrontEndURL, http.StatusTemporaryRedirect)
}

type OIDCToken struct {
	IdToken string `json:"id_token"`
}
