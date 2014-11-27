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
	Disabled  bool `db"-"`
}

// WeeklyPicks creates a []Picks representing a user's picks for the given week
func weeklyPicks(db *gorp.DbMap, username string) (picks []*Picks) {
	year, week := CurrentWeek(db)

	sql := `SELECT picks.*
		FROM picks
		JOIN games ON games.id = picks.game_id
		JOIN weeks ON weeks.id = games.week_id
		JOIN years ON years.id = weeks.year_id
		JOIN users ON users.id = picks.user_id
		WHERE years.year = $1 AND weeks.week = $2 AND users.email = $3
		ORDER BY games.date ASC`

	_, err := db.Select(&picks, sql, year, week, username)
	if err != nil {
		log.Fatalf("GetWeeklyPicks: %s", err.Error())
	}

	return
}

func weeklySelectedPicks(db *gorp.DbMap, username string) (picks []*Picks) {
	year, week := CurrentWeek(db)

	sql := `SELECT picks.*
		FROM picks
		JOIN games ON games.id = picks.game_id
		JOIN weeks ON weeks.id = games.week_id
		JOIN years ON years.id = weeks.year_id
		JOIN users ON users.id = picks.user_id
		WHERE years.year = $1 AND weeks.week = $2 AND users.email = $3 AND picks.selection != 0 AND picks.points != 0
		ORDER BY games.date ASC`

	_, err := db.Select(&picks, sql, year, week, username)
	if err != nil {
		log.Fatalf("GetWeeklyPicks: %s", err.Error())
	}

	return
}

// FormPicks gathers the neccessary information needed render a user's pick-em form
func FormPicks(db *gorp.DbMap, username string, selectedOnly bool) []FormPick {
	var picks []*Picks
	if selectedOnly {
		picks = weeklySelectedPicks(db, username)
	} else {
		picks = weeklyPicks(db, username)
	}

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

		disabled := time.Now().After(g.Date.Add(time.Hour * 5))

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
		}
		formGames = append(formGames, f)
	}

	return formGames
}
