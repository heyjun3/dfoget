//go:build wireinject

package presentation

import (
	"net/http"

	"github.com/google/wire"

	cfg "github.com/heyjun3/dforget/backend/config"
	"github.com/heyjun3/dforget/backend/presentation/auth"
)

func InitializeOIDCHandler(conf cfg.Config) *auth.OIDCHandler {
	wire.Build(
		provideHttpClient,
		auth.NewOIDCHandler,
	)
	return nil
}

func provideHttpClient() auth.HttpClient {
	return &http.Client{}
}
