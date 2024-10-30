package db

import (
	"database/sql"

	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed sql/schema.sql
var schema string

func OpenDB() (*sql.DB, error) {
	conn, err := sql.Open("sqlite", "db.sqlite")
	if err != nil {
		return nil, err
	}

	if _, err := conn.Exec(schema); err != nil {
		return nil, err
	}

	return conn, nil
}
