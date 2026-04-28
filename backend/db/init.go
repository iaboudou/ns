package db

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/mattn/go-sqlite3"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
)

func InitDB() (*sql.DB, error) {
	m, err := migrate.New("file://db/migration", "sqlite3://db/db.db")
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
