package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ameske/nfl-pickem/database"
)

func CurrentStandings(templatesDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, a := currentUser(r)

		year, week, err := database.CurrentWeek(time.Now())
		if err != nil {
			log.Fatal(err)
		}

		s, err := database.Standings(year, week)
		if err != nil && err != database.ErrNoStandings {
			log.Fatal(err)
		}

		err = Fetch(templatesDir, "standings.html").Execute(w, u, a, s)
		if err != nil {
			log.Println(err)
		}
	}
}
