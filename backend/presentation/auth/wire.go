//go:build wireinject

package auth

import (
	"net/http"

	"github.com/google/wire"

	cfg "github.com/heyjun3/dforget/backend/config"
)

func InitializeOIDCHandler(conf cfg.Config) *OIDCHandler {
	wire.Build(
		provideHttpClient,
		NewOIDCHandler,
	)
	return nil
}

func provideHttpClient() HttpClient {
	return &http.Client{}
}
