package database

import (
	"log"

	"github.com/coopernurse/gorp"
)

type StandingsPage struct {
	User   string `db:"first_name"`
	Points int    `db:"points"`
	Behind int    `db:"-"`
}

// Standings returns the state of the pick-em pool as of the given week in the
// requested year
func Standings(db *gorp.DbMap, year, week int) []*StandingsPage {
	var sql string

	// If it's the first week, we of course cannot deduct the lowest week yet
	if week == 1 {
		sql = `SELECT temp.first_name, SUM(temp.total) AS points
		FROM 
		    (SELECT users.first_name, weeks.id, SUM(picks.points) AS total
		    FROM picks
		    JOIN games ON picks.game_id = games.id
		    JOIN weeks ON weeks.id = games.week_id
		    JOIN years ON weeks.year_id = years.id
		    JOIN users ON users.id = picks.user_id 
		    WHERE picks.correct = True AND years.year = $1 AND weeks.week = $2
		    GROUP BY weeks.id, users.first_name) temp
		GROUP BY temp.first_name ORDER BY points DESC`

	} else {
		sql = `SELECT temp.first_name, SUM(temp.total) - MIN(temp.total) AS points
		FROM 
		    (SELECT users.first_name, weeks.id, SUM(picks.points) AS total
		    FROM picks
		    JOIN games ON picks.game_id = games.id
		    JOIN weeks ON weeks.id = games.week_id
		    JOIN years ON weeks.year_id = years.id
		    JOIN users ON users.id = picks.user_id 
		    WHERE picks.correct = True AND years.year = $1 AND weeks.week <= $2
		    GROUP BY weeks.id, users.first_name) temp
		GROUP BY temp.first_name ORDER BY points DESC`
	}

	var standings []*StandingsPage
	_, err := db.Select(&standings, sql, year, week)
	if err != nil {
		log.Fatalf("standings: %s", err.Error())
	}

	max := standings[0].Points

	for _, s := range standings {
		s.Behind = max - s.Points
	}

	return standings
}
