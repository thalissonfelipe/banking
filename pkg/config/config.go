package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	API      *apiConfig
	Postgres *postgresConfig
}

type apiConfig struct {
	Host string `envconfig:"API_HOST" default:"0.0.0.0"`
	Port string `envconfig:"API_PORT" default:"5000"`
}

type postgresConfig struct {
	DatabaseName     string `envconfig:"DB_NAME" default:"banking_db"`
	DatabaseUser     string `envconfig:"DB_USER" default:"postgres"`
	DatabasePassword string `envconfig:"DB_PASSWORD" default:"postgres"`
	DatabaseHost     string `envconfig:"DB_HOST" default:"localhost"`
	DatabasePort     string `envconfig:"DB_PORT" default:"5432"`
	SSLMode          string `envconfig:"DB_SSL_MODE" default:"disable"`
}

func (p postgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		p.DatabaseUser,
		p.DatabasePassword,
		p.DatabaseHost,
		p.DatabasePort,
		p.DatabaseName,
		p.SSLMode,
	)
}

func LoadConfig() (*config, error) {
	apiCfg, err := loadApiConfig()
	if err != nil {
		return nil, err
	}

	postgresCfg, err := loadPostgresConfig()
	if err != nil {
		return nil, err
	}

	return &config{apiCfg, postgresCfg}, nil
}

func loadApiConfig() (*apiConfig, error) {
	var apiCfg apiConfig
	err := envconfig.Process("API", &apiCfg)
	if err != nil {
		return nil, err
	}

	return &apiCfg, nil
}

func loadPostgresConfig() (*postgresConfig, error) {
	var postgresCfg postgresConfig
	err := envconfig.Process("DB", &postgresCfg)
	if err != nil {
		return nil, err
	}

	return &postgresCfg, nil
}
