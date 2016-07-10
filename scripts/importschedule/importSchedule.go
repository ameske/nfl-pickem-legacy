package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ameske/nfl-pickem/database"
	"github.com/ameske/nfl-pickem/schedule"
)

func main() {
	year := flag.Int("year", -1, "year")
	db := flag.String("db", "", "database file")
	splitYear := flag.Bool("split", false, "if week 17 is in the following calendar year")
	flag.Parse()

	if *year == -1 {
		log.Fatal("must provide year")
	}

	if *db == "" {
		log.Fatal("must provide db")
	}

	err := database.SetDefaultDb(*db)
	if err != nil {
		log.Fatal(err)
	}

	for week := 17; week <= 17; week++ {
		log.Println("Importing week", week)

		b, err := getScheduleHTML(*year, week)
		if err != nil {
			log.Fatal(err)
		}

		// Week 17 bleeds over to January
		if week == 17 && *splitYear {
			*year += 1
		}

		p := schedule.NewParser(*year, bytes.NewBuffer(b))

		games, err := p.Parse()
		if err != nil {
			log.Fatal(err)
		}

		for _, g := range games {
			err := database.AddGame(g.Date, g.Home, g.Away)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func getScheduleHTML(year, week int) ([]byte, error) {
	url := fmt.Sprintf("http://www.nfl.com/schedules/%d/REG%d", year, week)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
