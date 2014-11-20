package database

import (
	"log"
	"time"

	"github.com/coopernurse/gorp"
)

type Games struct {
	Id        int64     `db:"id"`
	WeekId    int64     `db:"week_id"`
	Date      time.Time `db:"date"`
	HomeId    int64     `db:"home_id"`
	AwayId    int64     `db:"away_id"`
	HomeScore int       `db:"home_score"`
	AwayScore int       `db:"away_score"`
}

func WeeklyGames(db *gorp.DbMap, year, week int) []Games {
	weekId := WeekId(db, year, week)

	var games []Games
	_, err := db.Select(&games, "SELECT * FROM games WHERE week_id = $1 ORDER BY date ASC", weekId)
	if err != nil {
		log.Fatalf("WeeklyGames: %s", err.Error())
	}

	return games
}
