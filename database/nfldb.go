package database

import (
	"database/sql"

	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
)

var db *gorp.DbMap

func SetDefaultDb(conn string) error {
	dbConn, err := sql.Open("sqlite3", conn)
	if err != nil {
		return err
	}

	err = dbConn.Ping()
	if err != nil {
		return err
	}

	db = &gorp.DbMap{Db: dbConn, Dialect: gorp.SqliteDialect{}}
	db.AddTableWithName(Users{}, "users").SetKeys(true, "Id")
	db.AddTableWithName(Pvs{}, "pvs").SetKeys(true, "Id")
	db.AddTableWithName(Teams{}, "teams").SetKeys(true, "Id")
	db.AddTableWithName(Years{}, "years").SetKeys(true, "Id")
	db.AddTableWithName(Weeks{}, "weeks").SetKeys(true, "Id")
	db.AddTableWithName(Games{}, "games").SetKeys(true, "Id")
	db.AddTableWithName(Picks{}, "picks").SetKeys(true, "Id")

	return err
}
