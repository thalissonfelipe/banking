package config

import (
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	API      *apiConfig
	Postgres *pgx.ConnConfig
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

func loadPostgresConfig() (*pgx.ConnConfig, error) {
	var postgresCfg postgresConfig
	err := envconfig.Process("DB", &postgresCfg)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		postgresCfg.DatabaseUser,
		postgresCfg.DatabasePassword,
		postgresCfg.DatabaseHost,
		postgresCfg.DatabasePort,
		postgresCfg.DatabaseName,
	)
	cfg, err := pgx.ParseConfig(uri)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
