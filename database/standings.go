package database

import "errors"

var (
	ErrNoStandings = errors.New("no standings")
)

type StandingsPage struct {
	User   string
	Points int
	Behind int
}

// Standings returns the state of the pick-em pool as of the given week in the
// requested year
func Standings(year, week int) ([]StandingsPage, error) {
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
		    WHERE picks.correct = 1 AND years.year = ?1 AND weeks.week = ?2
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
		    WHERE picks.correct = 1 AND years.year = ?1 AND weeks.week <= ?2
		    GROUP BY weeks.id, users.first_name) temp
		GROUP BY temp.first_name ORDER BY points DESC`
	}

	standings := make([]StandingsPage, 0)

	rows, err := db.Query(sql, year, week)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		tmp := StandingsPage{}
		err := rows.Scan(&tmp.User, &tmp.Points)
		if err != nil {
			return nil, err
		}
		standings = append(standings, tmp)
	}

	rows.Close()

	// Return just the users and 0 if there's nothing to pull out just yet
	if len(standings) == 0 {
		return emptyStandings()
	}

	max := standings[0].Points

	for _, s := range standings {
		s.Behind = max - s.Points
	}

	return standings, nil
}

func emptyStandings() ([]StandingsPage, error) {
	standings := make([]StandingsPage, 0)

	users, err := UserFirstNames()
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		tmp := StandingsPage{}
		tmp.User = u

		standings = append(standings, tmp)
	}

	return standings, nil
}
