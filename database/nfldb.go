package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func SetDefaultDb(conn string) error {
	db, err := sql.Open("sqlite3", conn)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	return err
}
