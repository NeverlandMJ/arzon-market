package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/NeverlandMJ/arzon-market/config"
)

func Connect(cfg config.Config) (*sql.DB, error) {
	//protocol: //login:password@host:port/yourDatabase'sName
	// dbURL := "postgres://sunbula:2307@localhost:5432/test"

	db, err := sql.Open(
		"pgx",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s",
			cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB,
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
