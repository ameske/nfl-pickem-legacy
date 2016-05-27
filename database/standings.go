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

type WeekByWeekStandingsPage struct {
	Name          string
	Scores        []int
	Total         int
	AdjustedTotal int
}

type rawWeekByWeek struct {
	Name   string
	Week   int
	Points int
}

// WeekByWeekStandings gathers information to show a standings page showing
// weekly results in addition to the total.
func WeekByWeekStandings(year, week int) ([]WeekByWeekStandingsPage, error) {
	sql := `SELECT users.first_name, weeks.week, SUM(picks.points)
		FROM picks
		JOIN games ON games.id = picks.game_id
		JOIN weeks ON weeks.id = games.week_id
		JOIN years ON years.id = weeks.year_id
		JOIN users ON users.id = picks.user_id
		WHERE picks.correct = 1 AND years.year = ?1 AND weeks.week <= ?2
		GROUP BY weeks.week, users.first_name;`

	rows, err := db.Query(sql, year, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// This gives us a list of (User, Week, Points) for every user and every week
	rawData := make([]rawWeekByWeek, 0)
	for rows.Next() {
		tmp := rawWeekByWeek{}
		err := rows.Scan(&tmp.Name, &tmp.Week, &tmp.Points)
		if err != nil {
			return nil, err
		}
		rawData = append(rawData, tmp)
	}

	// Now we need to condense it into the format of WeekByWeekStandingsPage
	standingsMap := make(map[string]*WeekByWeekStandingsPage, 0)
	for _, r := range rawData {
		wbwsp, ok := standingsMap[r.Name]
		// This is our first time seeing this user
		if !ok {
			wbwsp = &WeekByWeekStandingsPage{Scores: make([]int, 0)}
			wbwsp.Name = r.Name
			standingsMap[r.Name] = wbwsp
		}
		wbwsp.Scores = append(wbwsp.Scores, r.Points)
	}

	// Now, order the week by week standings page from 1st to last
	final := make([]WeekByWeekStandingsPage, 0)
	for _, v := range standingsMap {
		min := 9999
		for _, s := range v.Scores {
			if s <= min {
				min = s
			}
			v.Total += s
		}
		v.AdjustedTotal = v.Total - min
		final = append(final, *v)
	}

	for i := 0; i < len(final); i++ {
		adjT := final[i].AdjustedTotal
		for j := i + 1; j < len(final); j++ {
			if adjT < final[j].AdjustedTotal {
				tmp := final[i]
				final[i] = final[j]
				final[j] = tmp
			}
		}
	}

	return final, nil
}
