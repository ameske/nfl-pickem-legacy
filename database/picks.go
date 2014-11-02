package database

import (
	"log"
	"time"

	"github.com/coopernurse/gorp"
)

// Picks is a struct mapping to the corresponding postgres table
type Picks struct {
	Id        int64 `db:"id"`
	UserId    int64 `db:"user_id"`
	GameId    int64 `db:"game_id"`
	Selection int   `db:"selection"`
	Points    int   `db:"points"`
	Correct   bool  `db:"correct"`
}

// FormGame contains the information required to populate a picks HTML form
type FormGame struct {
	Id       int64
	Time     time.Time
	Away     string
	AwayNick string
	AwayId   int64
	Home     string
	HomeNick string
	HomeId   int64
}

// GetWeeklyPicks constructs an array of FormPicks for a given user and week which will be used to construct the picks html
func GetWeeklyPicks(db *gorp.DbMap, userId int64, year int, week int) []FormGame {
	weekId := WeekId(db, year, week)

	var picks []Picks
	_, err := db.Select(&picks, "SELECT picks.* FROM picks join games ON picks.game_id = games.id WHERE games.week_id = $1 AND picks.user_id = $2", weekId, userId)
	if err != nil {
		log.Fatalf("GetWeeklyPicks: %s", err.Error())
	}

	formGames := make([]FormGame, 0)
	for _, p := range picks {
		// Lookup the game information
		var g Games
		err := db.SelectOne(&g, "SELECT * FROM games WHERE id = $1", p.GameId)
		if err != nil {
			log.Fatalf("GetWeeklyPicks: %s", err.Error())
		}

		// Lookup the team information
		var h, a Teams
		err = db.SelectOne(&h, "SELECT * FROM teams WHERE id = $1", g.HomeId)
		if err != nil {
			log.Fatalf("GetWeeklyPicks: %s", err.Error())
		}
		err = db.SelectOne(&a, "SELECT * FROM teams WHERE id = $1", g.AwayId)
		if err != nil {
			log.Fatalf("GetWeeklyPicks: %s", err.Error())
		}

		// Construct the FormGame
		f := FormGame{
			Id:       p.Id,
			Time:     g.Date,
			Away:     a.City,
			AwayNick: a.Nickname,
			AwayId:   a.Id,
			Home:     h.City,
			HomeNick: h.Nickname,
			HomeId:   h.Id,
		}
		formGames = append(formGames, f)
	}

	return formGames
}
