package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	API      *apiConfig
	Postgres *postgresConfig
}

type apiConfig struct {
	Host string `envconfig:"API_HOST" default:"0.0.0.0"`
	Port string `envconfig:"API_PORT" default:"5000"`
}

func (a apiConfig) Address() string {
	return fmt.Sprintf("%s:%s", a.Host, a.Port)
}

type postgresConfig struct {
	DatabaseName     string `envconfig:"DB_NAME" default:"postgres"`
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

func LoadConfig() (*Config, error) {
	apiCfg, err := loadAPIConfig()
	if err != nil {
		return nil, err
	}

	postgresCfg, err := loadPostgresConfig()
	if err != nil {
		return nil, err
	}

	return &Config{apiCfg, postgresCfg}, nil
}

func loadAPIConfig() (*apiConfig, error) {
	var apiCfg apiConfig

	err := envconfig.Process("API", &apiCfg)
	if err != nil {
		return nil, fmt.Errorf("could not get api config: %w", err)
	}

	return &apiCfg, nil
}

func loadPostgresConfig() (*postgresConfig, error) {
	var postgresCfg postgresConfig

	err := envconfig.Process("DB", &postgresCfg)
	if err != nil {
		return nil, fmt.Errorf("could not get postgres config: %w", err)
	}

	return &postgresCfg, nil
}
