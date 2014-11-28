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
	var games []Games
	_, err := db.Select(&games, "SELECT games.* FROM games JOIN weeks ON weeks.id = games.week_id JOIN years ON years.id = weeks.year_id WHERE year = $1 AND week = $2 ORDER BY date ASC, games.id ASC", year, week)
	if err != nil {
		log.Fatalf("WeeklyGames: %s", err.Error())
	}

	return games
}

func GamesMap(games []Games) map[int64]Games {
	gm := make(map[int64]Games)
	for _, g := range games {
		g := g
		gm[g.Id] = g
	}

	return gm
}
