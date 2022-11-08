package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"remi/pkg/cmsql/cmsql"
)

type Postgres = cmsql.ConfigPostgres

func LoadPostgres() (*Postgres, error) {
	portStr := Coalesce(os.Getenv("POSTGRES_PORT"), "5432")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("POSTGRES_PORT must be number")
	}

	return &Postgres{
		Protocol: Coalesce(os.Getenv("POSTGRES_PROTOCOL"), "postgres"),
		Host:     Coalesce(os.Getenv("POSTGRES_HOST"), "postgres"),
		Port:     port,
		Username: Coalesce(os.Getenv("POSTGRES_USERNAME"), "postgres"),
		Password: Coalesce(os.Getenv("POSTGRES_PASSWORD"), "postgres"),
		Database: Coalesce(os.Getenv("POSTGRES_DATABASE"), "remi"),
		SSLMode:  Coalesce(os.Getenv("POSTGRES_SSL_MODE"), ""),
	}, nil
}

// HTTP ...
type HTTP struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Address ...
func (c HTTP) Address() string {
	if c.Port == 0 {
		log.Panic("Missing HTTP port")
	}
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

type Config struct {
	*Postgres `yaml:"postgres"`
	JWTSecret string `yaml:"jwt_secret"`
	HTTP      HTTP   `yaml:"http"`
	URL       string `yaml:"url"`
}

func Load() (cfg *Config, err error) {
	postgresCfg, err := LoadPostgres()
	if err != nil {
		return nil, err
	}
	jwtSecret := Coalesce(os.Getenv("JWT_SECRET"), "secret")
	httpPort, err := strconv.Atoi(Coalesce(os.Getenv("HTTP_PORT"), "8080"))
	if err != nil {
		return nil, fmt.Errorf("HTTP_PORT must be number")
	}
	url := os.Getenv("URL")

	return &Config{
		Postgres:  postgresCfg,
		JWTSecret: jwtSecret,
		HTTP: HTTP{
			Host: "",
			Port: httpPort,
		},
		URL: url,
	}, nil
}

func Coalesce(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
