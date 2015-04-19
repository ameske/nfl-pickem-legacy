package database

import (
	"database/sql"
	"fmt"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)

var db *gorp.DbMap

func SetDefaultDb(conn string) error {
	fullConnString := fmt.Sprintf("%s user=nfl database=nfl_app sslmode=disable", conn)
	dbConn, err := sql.Open("postgres", fullConnString)
	if err != nil {
		return err
	}

	err = dbConn.Ping()
	if err != nil {
		return err
	}

	db = &gorp.DbMap{Db: dbConn, Dialect: gorp.PostgresDialect{}}
	db.AddTableWithName(Users{}, "users").SetKeys(true, "Id")
	db.AddTableWithName(Pvs{}, "pvs").SetKeys(true, "Id")
	db.AddTableWithName(Teams{}, "teams").SetKeys(true, "Id")
	db.AddTableWithName(Years{}, "years").SetKeys(true, "Id")
	db.AddTableWithName(Weeks{}, "weeks").SetKeys(true, "Id")
	db.AddTableWithName(Games{}, "games").SetKeys(true, "Id")
	db.AddTableWithName(Picks{}, "picks").SetKeys(true, "Id")

	return err
}
