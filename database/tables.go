package database

import (
	"time"

	_ "github.com/lib/pq"
)

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
