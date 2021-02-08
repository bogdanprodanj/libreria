package config

import (
	"github.com/libreria/config/reader"
	"github.com/libreria/server/http"
	"github.com/libreria/storage/postgres"
)

type Config struct {
	LogLevel   string          `mapstructure:"log_level" default:"DEBUG"`
	HTTPServer http.Config     `mapstructure:"http_server"`
	Postgres   postgres.Config `mapstructure:"postgres"`
}

func New() (*Config, error) {
	cfg := new(Config)
	err := reader.Read(cfg)
	return cfg, err
}
