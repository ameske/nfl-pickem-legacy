package main

import (
	"encoding/json"
	"log"
	"os/exec"
	"strconv"

	"github.com/ameske/go_nfl/database"
	"github.com/codegangsta/cli"
)

type ResultsJson struct {
	Week      int    `json:"week"`
	Year      int    `json:"year"`
	Home      string `json:"home"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
}

func scores(c *cli.Context) {
	year, week := c.Int("year"), c.Int("week")
	if year == -1 || week == -1 {
		year, week = database.CurrentWeek(db)
	}

	db := database.NflDb()

	cmd := exec.Command("weeklyScores", strconv.Itoa(year), strconv.Itoa(week))
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Pipe: %s", err.Error())
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Start: %s", err.Error())
	}

	var results []ResultsJson
	err = json.NewDecoder(pipe).Decode(&results)
	if err != nil {
		log.Fatalf("Decode: %s", err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Wait: %s", err.Error())
	}

	for _, result := range results {
		// Lookup the year ID
		var yearId int64
		err = db.SelectOne(&yearId, "SELECT id FROM years WHERE year = $1", result.Year)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		// Lookup the week ID
		var weekId int64
		err = db.SelectOne(&weekId, "SELECT id FROM weeks WHERE year_id = $1 AND week = $2", yearId, result.Week)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		// Lookup the home team ID
		var teamId int64
		err = db.SelectOne(&teamId, "SELECT id FROM teams WHERE nickname = $1", result.Home)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		// Lookup the game ID based on the home team and week
		var game database.Games
		err = db.SelectOne(&game, "SELECT * FROM Games WHERE week_id = $1 and home_id = $2", weekId, teamId)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		// Update the scores for that game
		game.HomeScore = result.HomeScore
		game.AwayScore = result.AwayScore

		count, err := db.Update(&game)
		if err != nil {
			log.Fatalf("%s", err.Error())
		} else if count != 1 {
			log.Fatalf("More than one game was updated by the update statement.")
		}
	}
}
