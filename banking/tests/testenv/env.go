package testenv

import "github.com/jackc/pgx/v4"

// Global variables.
var (
	DB        *pgx.Conn
	ServerURL string
)
