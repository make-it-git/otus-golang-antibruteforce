package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Listen        string `env:"LISTEN"`
	RedisDSN      string `env:"REDIS_DSN" envDefault:"127.0.0.1:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`

	LimitLogin    int64 `env:"LIMIT_LOGIN"`
	LimitPassword int64 `env:"LIMIT_PASSWORD"`
	LimitIP       int64 `env:"LIMIT_IP"`

	TTL int `env:"BUCKET_TTL" envDefault:"60"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
