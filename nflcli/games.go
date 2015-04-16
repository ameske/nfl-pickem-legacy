package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/ameske/nfl-pickem/database"
	"github.com/codegangsta/cli"
	_ "github.com/lib/pq"
)

type GameJson struct {
	Week       int       `json:"week"`
	Home       string    `json:"home"`
	Away       string    `json:"away"`
	DateString string    `json:"date"`
	Date       time.Time `json:"-"`
	Year       int       `json:"year"`
}

func schedule(c *cli.Context) {
	games := make([]*GameJson, 0)

	// Open the 2014 schedule, and parse the json into our golang struct
	bytes, err := ioutil.ReadFile("json/2014/2014-Schedule.json")
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = json.Unmarshal(bytes, &games)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Convert the time string into a time.Time for postgres
	for _, game := range games {
		estDateString := game.DateString + " EST"
		game.Date, err = time.Parse("2006-01-02T15:04:05 MST", estDateString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	// Construct a mapping of team nicknames to ID's
	var teams []database.Teams
	_, err = db.Select(&teams, "SELECT * FROM teams")
	if err != nil {
		log.Fatalf(err.Error())
	}
	teamsMap := make(map[string]int64)
	for _, team := range teams {
		teamsMap[team.Nickname] = team.Id
	}

	// Now, construct a game row and add it into postgres
	for _, game := range games {
		// Note: Manually update year id!
		weekId, err := db.SelectInt("SELECT id FROM weeks WHERE week = $1 AND year_id = 1", game.Week)
		temp := database.Games{
			WeekId:    weekId,
			HomeId:    teamsMap[game.Home],
			AwayId:    teamsMap[game.Away],
			Date:      game.Date,
			HomeScore: -1,
			AwayScore: -1,
		}
		err = db.Insert(&temp)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}
