package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ameske/go_nfl/database"
)

type ResultsJson struct {
	Week      int    `json:"week"`
	Year      int    `json:"year"`
	Home      string `json:"home"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
}

var (
	week = flag.Int("week", -1, "Specify which week's results to load")
)

func main() {
	flag.Parse()
	if *week == -1 {
		log.Fatalf("--week is required.")
	}

	db := database.NflDb()

	resultBytes, err := ioutil.ReadFile(fmt.Sprintf("../json/2014/2014-Week%d-Results.json", *week))
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	results := make([]*ResultsJson, 0)
	err = json.Unmarshal(resultBytes, &results)
	if err != nil {
		log.Fatalf("%s", err.Error())
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
