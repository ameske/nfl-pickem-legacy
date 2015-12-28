package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var here = time.Now().Location()

type Games struct {
	Games []interface{}
}

type Game struct {
	Away     string
	Home     string
	Week     int
	Day      int
	Month    int
	Year     int
	Time     string
	Meridiem string
	Type     string
}

func gameMapToStruct(m map[string]interface{}) Game {
	var g Game

	g.Away = m["away"].(string)
	g.Home = m["home"].(string)
	g.Week = int(m["week"].(float64))
	g.Day = int(m["day"].(float64))
	g.Month = int(m["month"].(float64))
	g.Year = int(m["year"].(float64))
	g.Time = m["time"].(string)
	g.Type = m["season_type"].(string)
	meridiem, _ := m["meridiem"]
	if meridiem != nil {
		g.Meridiem = meridiem.(string)
	}

	return g
}

func main() {
	bytes, err := ioutil.ReadFile("schedule.json")
	if err != nil {
		log.Fatal(err)
	}

	var rawGames Games
	err = json.Unmarshal(bytes, &rawGames)
	if err != nil {
		log.Fatal(err)
	}

	for _, rawGame := range rawGames.Games {
		rawGameParts, ok := rawGame.([]interface{})
		if !ok {
			log.Fatal("couldn't split into []interface")
		}

		if len(rawGameParts) != 2 {
			log.Fatal("malformed raw game")
		}
		rawGameMap, ok := rawGameParts[1].(map[string]interface{})

		game := gameMapToStruct(rawGameMap)
		parseAndPrintGame(game)
	}

}

func parseAndPrintGame(g Game) {
	if g.Year != 2015 || g.Type != "REG" {
		return
	}

	timeParts := strings.Split(g.Time, ":")

	hour, err := strconv.ParseInt(timeParts[0], 10, 32)
	if err != nil {
		log.Fatal(err)
	}

	if g.Meridiem == "PM" {
		hour += 12
	}

	minute, err := strconv.ParseInt(timeParts[1], 10, 32)
	if err != nil {
		log.Fatal(err)
	}

	gameTime := time.Date(g.Year, time.Month(g.Month), g.Day, int(hour), int(minute), 0, 0, here)

	fmt.Fprintf(os.Stdout, "INSERT INTO games(week_id, date, home_id, away_id) VALUES((SELECT weeks.id FROM weeks JOIN years ON years.id = weeks.year_id WHERE weeks.week = %d AND years.year = %d), %d, (SELECT id FROM teams WHERE abbreviation = '%s'), (SELECT id FROM teams WHERE abbreviation = '%s'));\n", g.Week, g.Year, gameTime.Unix(), g.Home, g.Away)
}
