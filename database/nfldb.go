package database

import (
	"database/sql"
	"log"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)

func NflDb() *gorp.DbMap {
	db, err := sql.Open("postgres", "host=/run/postgresql user=nfl database=nfl_app sslmode=disable")
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf(err.Error())
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(Users{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(Pvs{}, "pvs").SetKeys(true, "Id")
	dbmap.AddTableWithName(Teams{}, "teams").SetKeys(true, "Id")
	dbmap.AddTableWithName(Years{}, "years").SetKeys(true, "Id")
	dbmap.AddTableWithName(Weeks{}, "weeks").SetKeys(true, "Id")
	dbmap.AddTableWithName(Games{}, "games").SetKeys(true, "Id")
	dbmap.AddTableWithName(Picks{}, "picks").SetKeys(true, "Id")

	return dbmap

}
