package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)

type Users struct {
	Id        int64     `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Admin     bool      `db:"admin"`
	LastLogin time.Time `db:"last_login"`
	Password  string    `db:"password"`
}

type Pvs struct {
	Id    int64  `db:"id"`
	Type  string `db:"type"`
	Seven int    `db:"seven"`
	Five  int    `db:"five"`
	Three int    `db:"three"`
	One   int    `db:"one"`
}

type Teams struct {
	Id       int64  `db:"id"`
	City     string `db:"city"`
	Nickname string `db:"nickname"`
	Stadium  string `db:"stadium"`
}

type Years struct {
	Id   int64 `db:"id"`
	Year int   `db:"year"`
}

type Weeks struct {
	Id     int64 `db:"id"`
	YearId int64 `db:"year_id"`
	PvsId  int64 `db:"pvs_id"`
	Week   int   `db:"week"`
}

type Games struct {
	Id        int64     `db:"id"`
	WeekId    int64     `db:"week_id"`
	Date      time.Time `db:"date"`
	HomeId    int64     `db:"home_id"`
	AwayId    int64     `db:"away_id"`
	HomeScore int       `db:"home_score"`
	AwayScore int       `db:"away_score"`
}

type Picks struct {
	Id        int64 `db:"id"`
	UserId    int64 `db:"user_id"`
	GameId    int64 `db:"game_id"`
	Selection int   `db:"selection"`
	Points    int   `db:"points"`
	Correct   bool  `db:"correct"`
}

func NflDb() *gorp.DbMap {
	db, err := sql.Open("postgres", "user=nfl database=nfl_app sslmode=disable")
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
