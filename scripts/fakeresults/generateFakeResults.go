package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ameske/nfl-pickem/database"
)

func main() {
	year := flag.Int("year", -1, "year")
	week := flag.Int("week", -1, "week")
	thur := flag.Bool("thur", false, "generate result for thursday night game")
	sun1 := flag.Bool("sun1", false, "generate results for sunday 1:00 games")
	sun4 := flag.Bool("sun4", false, "generate results for sunday 4:00 games")
	sun8 := flag.Bool("sun8", false, "generate results for sunday night game")
	mon := flag.Bool("mon", false, "generate results for monday night game")

	flag.Parse()

	if *year == -1 || *week == -1 {
		log.Fatal("year and week required")
	}

	err := database.SetDefaultDb("nfl-test.db")
	if err != nil {
		log.Fatal(err)
	}

	games, err := database.WeeklyGames(*year, *week)
	if err != nil {
		log.Fatal(err)
	}

	results, err := os.Create("results.json")
	if err != nil {
		log.Fatal(err)
	}
	defer results.Close()

	enc := json.NewEncoder(results)

	rand.Seed(time.Now().Unix())

	rj := make([]ResultsJson, 0)
	for _, g := range games {
		if *thur && g.Date.Weekday() == time.Thursday {
			rj = append(rj, generateRandomResult(*year, *week, g.HomeNickname))
		}

		if *sun1 && g.Date.Weekday() == time.Sunday && g.Date.Hour() == 13 {
			rj = append(rj, generateRandomResult(*year, *week, g.HomeNickname))
		}

		if *sun4 && g.Date.Weekday() == time.Sunday && g.Date.Hour() == 16 {
			rj = append(rj, generateRandomResult(*year, *week, g.HomeNickname))
		}

		if *sun8 && g.Date.Weekday() == time.Sunday && g.Date.Hour() == 20 {
			rj = append(rj, generateRandomResult(*year, *week, g.HomeNickname))
		}

		if *mon && g.Date.Weekday() == time.Monday {
			rj = append(rj, generateRandomResult(*year, *week, g.HomeNickname))
		}
	}

	err = enc.Encode(&rj)
	if err != nil {
		log.Fatal(err)
	}
}

type ResultsJson struct {
	Week      int    `json:"week"`
	Year      int    `json:"year"`
	Home      string `json:"home"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
}

func generateRandomResult(year int, week int, home string) ResultsJson {
	hscore, ascore := generateRandomScore()

	return ResultsJson{
		Week:      week,
		Year:      year,
		Home:      home,
		HomeScore: hscore,
		AwayScore: ascore,
	}
}

func generateRandomScore() (home, away int) {
	home = rand.Intn(64)
	away = rand.Intn(64)

	return home, away
}
