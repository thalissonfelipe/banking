package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	API      APIConfig
	GRPC     GRPCConfig
	Postgres PostgresConfig
}

type APIConfig struct {
	Host string `envconfig:"API_HOST" default:"0.0.0.0"`
	Port string `envconfig:"API_PORT" default:"5000"`
}

func (a APIConfig) Address() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

type GRPCConfig struct {
	Address string `envconfig:"GRPC_ADDRESS" default:"0.0.0.0:9000"`
}

type PostgresConfig struct {
	DatabaseName     string `envconfig:"DB_NAME" default:"postgres"`
	DatabaseUser     string `envconfig:"DB_USER" default:"postgres"`
	DatabasePassword string `envconfig:"DB_PASSWORD" default:"postgres"`
	DatabaseHost     string `envconfig:"DB_HOST" default:"localhost"`
	DatabasePort     string `envconfig:"DB_PORT" default:"5432"`
	SSLMode          string `envconfig:"DB_SSL_MODE" default:"disable"`
}

func (p PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.DatabaseUser,
		p.DatabasePassword,
		p.DatabaseHost,
		p.DatabasePort,
		p.DatabaseName,
		p.SSLMode,
	)
}

func LoadConfig() (*Config, error) {
	const noPrefix = ""

	var config Config

	err := envconfig.Process(noPrefix, &config)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	return &config, nil
}
