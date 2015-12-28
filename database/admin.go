package database

import "log"

type AdminPickRow struct {
	Home  string
	Away  string
	Picks []UserPick
}

type UserPick struct {
	Home      string
	Away      string
	Id        int64
	Selection int
	Points    int
}

func homeTeams(year, week int) []string {
	sql := `SELECT teams.abbreviation
	  FROM teams
	  JOIN games ON games.home_id = teams.id
	  JOIN weeks ON games.week_id = weeks.id
	  JOIN years ON weeks.year_id = years.id
	  WHERE years.year = ?1 AND weeks.week = ?2
	  ORDER BY games.date ASC, games.id ASC`

	rows, err := db.Query(sql, year, week)
	if err != nil {
		log.Fatal(err)
	}

	teams := make([]string, 0)
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			log.Fatal(err)
		}
		teams = append(teams, tmp)
	}

	rows.Close()

	return teams
}

func awayTeams(year, week int) []string {
	sql := `SELECT teams.abbreviation
	  FROM teams
	  JOIN games ON games.away_id = teams.id
	  JOIN weeks ON games.week_id = weeks.id
	  JOIN years ON weeks.year_id = years.id
	  WHERE years.year = ?1 AND weeks.week = ?2
	  ORDER BY games.date ASC, games.id ASC`

	rows, err := db.Query(sql, year, week)
	if err != nil {
		log.Fatal(err)
	}

	teams := make([]string, 0)
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			log.Fatal(err)
		}
		teams = append(teams, tmp)
	}

	rows.Close()

	return teams
}

func AdminForm(year, week int) ([]string, []AdminPickRow) {
	home := homeTeams(year, week)
	away := awayTeams(year, week)

	if len(home) != len(away) {
		log.Fatal("len(home) != len(away)")
	}

	sql := `SELECT picks.id, picks.selection, picks.points
	       FROM picks
	       JOIN games ON games.id = picks.game_id
	       JOIN teams ON teams.id = games.home_id
	       JOIN users ON users.id = picks.user_id
	       JOIN weeks ON weeks.id = games.week_id
	       JOIN years ON years.id = weeks.year_id
	       WHERE years.year = ?1 AND weeks.week = ?2 AND teams.abbreviation = ?3
	       ORDER BY games.date ASC, games.id ASC, users.id ASC`

	formRows := make([]AdminPickRow, 0, len(home))

	for i := 0; i < len(home); i++ {
		apr := AdminPickRow{}
		apr.Home = home[i]
		apr.Away = away[i]

		rows, err := db.Query(sql, year, week, home[i])
		if err != nil {
			log.Fatal(err)
		}

		apr.Picks = make([]UserPick, 0)
		for rows.Next() {
			tmp := UserPick{}
			err := rows.Scan(&tmp.Id, &tmp.Selection, &tmp.Points)
			if err != nil {
				log.Fatal(err)
			}
			tmp.Home = home[i]
			tmp.Away = away[i]
			apr.Picks = append(apr.Picks, tmp)
		}
		rows.Close()

		formRows = append(formRows, apr)
	}

	users := usernames()

	return users, formRows
}

func usernames() []string {
	rows, err := db.Query("SELECT first_name FROM users ORDER BY id ASC")
	if err != nil {
		log.Fatal(err)
	}

	usernames := make([]string, 0)
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			log.Fatal(err)
		}

		usernames = append(usernames, tmp)
	}

	rows.Close()

	return usernames
}
