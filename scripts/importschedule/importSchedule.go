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

var (
	urls = map[int]string{1: "https://web.archive.org/web/20160708083559/http://www.nfl.com/schedules/2016/REG1",
		2:  "https://web.archive.org/web/20160708100600/http://www.nfl.com/schedules/2016/REG2",
		3:  "https://web.archive.org/web/20160708100603/http://www.nfl.com/schedules/2016/REG3",
		4:  "https://web.archive.org/web/20160708100606/http://www.nfl.com/schedules/2016/REG4",
		5:  "https://web.archive.org/web/20160708100609/http://www.nfl.com/schedules/2016/REG5",
		6:  "https://web.archive.org/web/20160708100612/http://www.nfl.com/schedules/2016/REG6",
		7:  "https://web.archive.org/web/20160708100615/http://www.nfl.com/schedules/2016/REG7",
		8:  "https://web.archive.org/web/20160708100618/http://www.nfl.com/schedules/2016/REG8",
		9:  "https://web.archive.org/web/20160708100621/http://www.nfl.com/schedules/2016/REG9",
		10: "https://web.archive.org/web/20160708100624/http://www.nfl.com/schedules/2016/REG10",
		11: "https://web.archive.org/web/20160708100627/http://www.nfl.com/schedules/2016/REG11",
		12: "https://web.archive.org/web/20160708100630/http://www.nfl.com/schedules/2016/REG12",
		13: "https://web.archive.org/web/20160708100634/http://www.nfl.com/schedules/2016/REG13",
		14: "https://web.archive.org/web/20160708100637/http://www.nfl.com/schedules/2016/REG14",
		15: "https://web.archive.org/web/20160708100639/http://www.nfl.com/schedules/2016/REG15",
		16: "https://web.archive.org/web/20160708100642/http://www.nfl.com/schedules/2016/REG16",
		17: "https://web.archive.org/web/20160708100646/http://www.nfl.com/schedules/2016/REG17"}
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

	for week := 1; week <= 17; week++ {
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
			if week == 17 && *splitYear {
				err = database.AddGame(g.Date, g.Home, g.Away, true)
			} else {
				err = database.AddGame(g.Date, g.Home, g.Away, false)
			}
			if err != nil {
				fmt.Printf("Date: %v\tHome: %s\tAway: %s\t", g.Date, g.Home, g.Away)
				log.Fatal(err)
			}
		}
	}
}

func getScheduleHTML(year, week int) ([]byte, error) {
	url := urls[week]

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
