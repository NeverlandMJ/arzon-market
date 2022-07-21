package postgres

import (
	"database/sql"
	"fmt"
)

const (
	pgHost = "localhost"
	pgPort = 5432
	pgUser = "sunbula"
	pgPass = "2307"
	pdDB   = "test"
)

func Connect() (*sql.DB, error) {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			pgHost, pgPort, pgUser, pgPass, pdDB,
		),
	)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}