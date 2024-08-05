package server

import (
	"fmt"
	"os"
)

type Config struct {
	db DBConfig
}

type DBConfig struct {
	user     string
	password string
	host     string
	port     string
	name     string
}

func NewConfig() Config {
	return Config{
		db: DBConfig{
			user:     os.Getenv("DB_USER"),
			password: os.Getenv("DB_PASSWORD"),
			host:     os.Getenv("DB_HOST"),
			port:     os.Getenv("DB_PORT"),
			name:     os.Getenv("DB_NAME"),
		},
	}
}

func (c Config) DBDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.db.user, c.db.password, c.db.host, c.db.port, c.db.name)
}
func (c Config) TESTDBDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/test?sslmode=disable",
		c.db.user, c.db.password, c.db.host, c.db.port)
}
