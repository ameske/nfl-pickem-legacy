package database

import (
	"log"
	"time"
)

type Games struct {
	Id        int64
	WeekId    int64
	Date      int64
	HomeId    int64
	AwayId    int64
	HomeScore int
	AwayScore int
}

func gamePickInfo(id int64) (awayId int64, homeId int64, time int64) {
	row := db.QueryRow("SELECT away_id, home_id, date FROM games WHERE id = ?1", id)
	err := row.Scan(&awayId, &homeId, &time)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func GamesBySeason(year int) []Games {
	games := make([]Games, 0)

	rows, err := db.Query(`SELECT games.id, games.week_id, games.date, games.home_id, games.away_id, games.home_score, games.away_score
			       FROM games 
			       JOIN weeks ON weeks.id = games.week_id
			       JOIN years ON years.id = weeks.year_id
			       WHERE years.year = ?1`, year)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		tmp := Games{}
		err := rows.Scan(&tmp.Id, &tmp.WeekId, &tmp.Date, &tmp.HomeId, &tmp.AwayId, &tmp.HomeScore, &tmp.AwayScore)
		if err != nil {
			log.Fatal(err)
		}
	}

	rows.Close()

	return games
}

func WeeklyGames(year, week int) []Games {
	var games []Games

	sql := `SELECT games.id, games.week_id, games.date, games.home_id, games.away_id, games.home_score, games.away_score
		FROM games
		JOIN weeks ON weeks.id = games.week_id
		JOIN years ON years.id = weeks.year_id
		WHERE year = ?1 AND week = ?2 ORDER BY date ASC, games.id ASC`

	rows, err := db.Query(sql, year, week)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		tmp := Games{}
		err := rows.Scan(&tmp.Id, &tmp.WeekId, &tmp.Date, &tmp.HomeId, &tmp.AwayId, &tmp.HomeScore, &tmp.AwayScore)
		if err != nil {
			log.Fatal(err)
		}
	}

	rows.Close()

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

func UpdateScores(week int, year int, homeTeam string, homeScore int, awayScore int) error {
	// Lookup the week ID
	weekId, err := weekID(week, year)
	if err != nil {
		log.Fatal(err)
	}

	// Lookup the home team ID
	var teamId int64
	row := db.QueryRow("SELECT id FROM teams WHERE nickname = ?1", homeTeam)
	err = row.Scan(&teamId)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`UPDATE games
			      SET home_score = ?1, away_score = ?2
			      WHERE week_id = ?3 AND home_id = ?4`, homeScore, awayScore, weekId, teamId)

	return err
}

func nextSunday() int64 {
	t := time.Now()

	return t.AddDate(0, 0, 7-int(t.Weekday())).Unix()
}
