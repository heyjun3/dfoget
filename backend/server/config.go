package server

import (
	"fmt"
	"os"
)

type Config struct {
	db          DBConfig
	oidc        OIDCConfig
	frontEndURL string
}

type DBConfig struct {
	user     string
	password string
	host     string
	port     string
	name     string
}

type OIDCConfig struct {
	redirectUri  string
	clientId     string
	clientSecret string
	tokenUrl     string
}

type ConfigOption func(Config) Config

func NewConfig(opts ...ConfigOption) Config {
	conf := Config{
		db: DBConfig{
			user:     os.Getenv("DB_USER"),
			password: os.Getenv("DB_PASSWORD"),
			host:     os.Getenv("DB_HOST"),
			port:     os.Getenv("DB_PORT"),
			name:     os.Getenv("DB_NAME"),
		},
		oidc: OIDCConfig{
			redirectUri:  os.Getenv("OIDC_REDIRECT_URI"),
			clientId:     os.Getenv("OIDC_CLIENT_ID"),
			clientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
			tokenUrl:     os.Getenv("OIDC_TOKEN_URL"),
		},
		frontEndURL: os.Getenv("FRONTEND_URL"),
	}
	for _, opt := range opts {
		conf = opt(conf)
	}
	return conf
}

func (c Config) DBDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.db.user, c.db.password, c.db.host, c.db.port, c.db.name)
}

func WithDBName(name string) func(Config) Config {
	return func(conf Config) Config {
		conf.db.name = name
		return conf
	}
}
