package testenv

import "github.com/jackc/pgx/v4"

var (
	DB        *pgx.Conn
	ServerURL string
)
