package database

import (
	"log"
	"time"
)

type Game struct {
	Date             time.Time
	Home             string
	HomeNickname     string
	HomeAbbreviation string
	HomeScore        int
	Away             string
	AwayNickname     string
	AwayAbbreviation string
	AwayScore        int
}

func WeeklyGames(year, week int) ([]Game, error) {
	sql := `SELECT games.date, hometeam.city, hometeam.nickname, hometeam.abbreviation, games.home_score, awayteam.city, awayteam.nickname, awayteam.abbreviation, games.away_score
		FROM games
		JOIN weeks ON weeks.id = games.week_id
		JOIN years ON years.id = weeks.year_id
		JOIN teams AS hometeam ON games.home_id = hometeam.id
		JOIN teams AS awayteam ON games.away_id = awayteam.id
		WHERE year = ?1 AND week = ?2 ORDER BY games.date ASC, games.id ASC`
	rows, err := db.Query(sql, year, week)
	if err != nil {
		log.Println("weekly games")
		return nil, err
	}

	var games []Game
	for rows.Next() {
		tmp := Game{}
		var d int64
		err := rows.Scan(&d, &tmp.Home, &tmp.HomeNickname, &tmp.HomeAbbreviation, &tmp.HomeScore, &tmp.Away, &tmp.AwayNickname, &tmp.AwayAbbreviation, &tmp.AwayScore)
		if err != nil {
			return nil, err
		}
		tmp.Date = time.Unix(d, 0)
		games = append(games, tmp)
	}

	rows.Close()

	return games, nil
}

func UpdateScore(week int, year int, homeTeam string, homeScore int, awayScore int) error {
	//sqlite3 makes this hard on us....so we have to do this in a couple of steps
	sql := `SELECT games.id FROM games
		JOIN weeks ON games.week_id = weeks.id
		JOIN years ON weeks.year_id = years.id
		JOIN teams ON games.home_id = teams.id
		WHERE weeks.week = ?1 AND years.year = ?2 AND teams.nickname = ?3`

	var gameId int64
	err := db.QueryRow(sql, week, year, homeTeam).Scan(&gameId)
	if err != nil {
		return err
	}

	_, err = db.Exec(`UPDATE games
			  SET home_score = ?2, away_score = ?3
			  WHERE id = ?1`, gameId, homeScore, awayScore)

	return err
}

func AddGame(date time.Time, homeTeam string, awayTeam string, wk17splitYear bool) error {
	_, week, err := CurrentWeek(date)
	if err != nil {
		return err
	}

	if wk17splitYear {
		_, err = db.Exec(`INSERT INTO games(week_id, date, home_id, away_id)
			 VALUES((SELECT weeks.id FROM weeks JOIN years ON weeks.year_id = years.id WHERE years.year = ?1 AND weeks.week = ?2), ?3, (SELECT id FROM teams WHERE nickname = ?4), (SELECT id FROM teams WHERE nickname = ?5))`, date.Year()-1, week, date.Unix(), homeTeam, awayTeam)
	} else {
		_, err = db.Exec(`INSERT INTO games(week_id, date, home_id, away_id)
			 VALUES((SELECT weeks.id FROM weeks JOIN years ON weeks.year_id = years.id WHERE years.year = ?1 AND weeks.week = ?2), ?3, (SELECT id FROM teams WHERE nickname = ?4), (SELECT id FROM teams WHERE nickname = ?5))`, date.Year(), week, date.Unix(), homeTeam, awayTeam)

	}

	return err
}
