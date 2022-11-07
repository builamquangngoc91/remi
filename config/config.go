package config

import (
	cc "remi/pkg/config"
)

type Config struct {
	cc.Postgres `yaml:"postgres"`
	JWTSecret   string  `yaml:"jwt_secret"`
	HTTP        cc.HTTP `yaml:"http"`
}

func Default() Config {
	return Config{
		Postgres:  cc.DefaultPostgres(),
		JWTSecret: "secret",
		HTTP: cc.HTTP{
			Host: "",
			Port: 8080,
		},
	}
}

func Load() (cfg Config, err error) {
	err = cc.LoadWithDefault(&cfg, Default())
	return cfg, err
}
