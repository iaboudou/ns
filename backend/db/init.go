package db

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
)

func InitDB() (*sql.DB, error) {
	m, err := migrate.New(
		"file://pkg/db/migrations/sqlite/",
		"sqlite3://pkg/db/db.db",
	)
	if err != nil {
		return nil, err
	}

	// Run migrations
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	db, err := sql.Open("sqlite3", "./db/db.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA journal_mode = WAL;")
	if err != nil {
		return nil, err
	}
	return db, nil
}
