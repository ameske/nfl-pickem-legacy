package database

import (
	"log"
	"math/rand"
	"time"
)

type Games struct {
	Id        int64 `db:"id"`
	WeekId    int64 `db:"week_id"`
	Date      int64 `db:"date"`
	HomeId    int64 `db:"home_id"`
	AwayId    int64 `db:"away_id"`
	HomeScore int   `db:"home_score"`
	AwayScore int   `db:"away_score"`
}

func GamesBySeason(year int) []Games {
	var games []Games
	_, err := db.Select(&games, `SELECT games.*
				 FROM games 
				 JOIN weeks ON weeks.id = games.week_id
				 JOIN years ON years.id = weeks.year_id
				 WHERE years.year = ?1`, year)
	if err != nil {
		log.Fatalf("GamesBySeason: %s", err.Error())
	}

	return games
}

func WeeklyGames(year, week int) []Games {
	var games []Games
	_, err := db.Select(&games, "SELECT games.* FROM games JOIN weeks ON weeks.id = games.week_id JOIN years ON years.id = weeks.year_id WHERE year = ?1 AND week = ?2 ORDER BY date ASC, games.id ASC", year, week)
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

func UpdateScores(week int, year int, homeTeam string, homeScore int, awayScore int) error {
	// Lookup the year ID
	var yearId int64
	err := db.SelectOne(&yearId, "SELECT id FROM years WHERE year = ?1", year)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// Lookup the week ID
	var weekId int64
	err = db.SelectOne(&weekId, "SELECT id FROM weeks WHERE year_id = ?1 AND week = ?2", yearId, week)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// Lookup the home team ID
	var teamId int64
	err = db.SelectOne(&teamId, "SELECT id FROM teams WHERE nickname = ?1", homeTeam)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// Lookup the game ID based on the home team and week
	var game Games
	err = db.SelectOne(&game, "SELECT * FROM Games WHERE week_id = ?1 and home_id = ?2", weekId, teamId)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// Update the scores for that game
	game.HomeScore = homeScore
	game.AwayScore = awayScore

	count, err := db.Update(&game)
	if err != nil {
		log.Fatalf("%s", err.Error())
	} else if count != 1 {
		log.Fatalf("More than one game was updated by the update statement.")
	}

	return nil
}

func AddGame(week int, year int, homeTeam string, awayTeam string, date int64) error {
	yearID, err := db.SelectInt("SELECT id FROM years WHERE year = ?1", year)
	if err != nil {
		return err
	}

	weekId, err := db.SelectInt("SELECT id FROM weeks WHERE week = ?1 AND year_id = 1", week, yearID)
	if err != nil {
		return err
	}

	homeID, err := db.SelectInt("SELECT id from teams WHERE nickname = ?1", homeTeam)
	if err != nil {
		return err
	}

	awayID, err := db.SelectInt("SELECT id from teams WHERE nickname = ?1", awayTeam)
	if err != nil {
		return err
	}

	temp := Games{
		WeekId:    weekId,
		HomeId:    homeID,
		AwayId:    awayID,
		Date:      date,
		HomeScore: -1,
		AwayScore: -1,
	}

	return db.Insert(&temp)
}

func CreateRandomGames(week int, year int) error {
	weekID, err := weekID(week, year)
	if err != nil {
		return err
	}

	options := rand.Perm(32)

	for i := 0; i < 32; i += 2 {
		g := Games{
			WeekId:    weekID,
			HomeId:    int64(options[i]) + 1,
			AwayId:    int64(options[i+1]) + 1,
			Date:      nextSunday(),
			HomeScore: -1,
			AwayScore: -1,
		}

		err := db.Insert(&g)
		if err != nil {
			return err
		}
	}

	return nil
}

func nextSunday() int64 {
	t := time.Now()

	return t.AddDate(0, 0, 7-int(t.Weekday())).Unix()
}
