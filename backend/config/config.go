package server

import (
	"fmt"
	"os"
)

type Config struct {
	db          DBConfig
	OIDC        OIDCConfig
	FrontEndURL string
}

type DBConfig struct {
	user     string
	password string
	host     string
	port     string
	name     string
}

type OIDCConfig struct {
	RedirectUri  string
	ClientId     string
	ClientSecret string
	TokenUrl     string
	Pubkey       string
}

type ConfigOption func(Config) Config

func NewConfig(opts ...ConfigOption) Config {
	pubkey := os.Getenv("OIDC_PUBLIC_KEY")
	if pubkey == "" {
		panic("don't set oidc public key")
	}
	conf := Config{
		db: DBConfig{
			user:     os.Getenv("DB_USER"),
			password: os.Getenv("DB_PASSWORD"),
			host:     os.Getenv("DB_HOST"),
			port:     os.Getenv("DB_PORT"),
			name:     os.Getenv("DB_NAME"),
		},
		OIDC: OIDCConfig{
			RedirectUri:  os.Getenv("OIDC_REDIRECT_URI"),
			ClientId:     os.Getenv("OIDC_CLIENT_ID"),
			ClientSecret: os.Getenv("OIDC_CLIENT_SECRET"),
			TokenUrl:     os.Getenv("OIDC_TOKEN_URL"),
			Pubkey:       pubkey,
		},
		FrontEndURL: os.Getenv("FRONTEND_URL"),
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
func WithPubKey(pubKey string) func(Config) Config {
	return func(conf Config) Config {
		conf.OIDC.Pubkey = pubKey
		return conf
	}
}
