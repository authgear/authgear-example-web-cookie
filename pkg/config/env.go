package config

import (
	"github.com/authgear/authgear-server/pkg/util/log"
	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	ListenAddr       string    `envconfig:"LISTEN_ADDR" default:"0.0.0.0:18001"`
	LogLevel         log.Level `envconfig:"LOG_LEVEL" default:"info"`
	AuthgearEndpoint string    `envconfig:"AUTHGEAR_ENDPOINT" default:"http://accounts.portal.localhost:3100"`
	DefaultClientID  string    `envconfig:"DEFAULT_CLIENT_ID" default:"example"`
}

func LoadConfigFromEnv() (*EnvConfig, error) {
	var cfg EnvConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
