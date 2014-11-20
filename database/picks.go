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

//FormPick contains the information required to populate a picks HTML Form
type FormPick struct {
	Id        int64
	Time      time.Time
	Away      string
	AwayNick  string
	AwayId    int64
	Home      string
	HomeNick  string
	HomeId    int64
	Selection int
	Points    int
	Disabled  bool
	Graded    bool
	Correct   bool
}

// WeeklyPicks creates a []Picks representing a user's picks for the given week
func WeeklyPicks(db *gorp.DbMap, userId int64, year int, week int) (picks []*Picks) {
	weekId := WeekId(db, year, week)

	_, err := db.Select(&picks, "SELECT picks.* FROM picks join games ON picks.game_id = games.id WHERE games.week_id = $1 AND picks.user_id = $2 ORDER BY games.date ASC", weekId, userId)
	if err != nil {
		log.Fatalf("GetWeeklyPicks: %s", err.Error())
	}

	return
}

// FormPicks gathers the neccessary information needed render a user's pick-em form
func FormPicks(db *gorp.DbMap, username string, year int, week int) []FormPick {
	userId := UserId(db, username)

	picks := WeeklyPicks(db, userId, year, week)
	formGames := make([]FormPick, 0)
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

		disabled, graded := false, false
		if time.Now().After(g.Date) {
			disabled = true
		}
		if time.Since(g.Date) > time.Duration(48)*time.Hour && g.HomeScore != -1 && g.AwayScore != -1 {
			graded = true
		}

		// Construct the FormPick
		f := FormPick{
			Id:        p.Id,
			Time:      g.Date,
			Away:      a.City,
			AwayNick:  a.Nickname,
			AwayId:    a.Id,
			Home:      h.City,
			HomeNick:  h.Nickname,
			HomeId:    h.Id,
			Selection: p.Selection,
			Points:    p.Points,
			Disabled:  disabled,
			Graded:    graded,
			Correct:   p.Correct,
		}
		formGames = append(formGames, f)
	}

	return formGames
}
