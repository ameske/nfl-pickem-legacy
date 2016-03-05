package main

import (
	"log"
	"net/http"

	"github.com/ameske/nfl-pickem/database"
)

func Standings(templatesDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, a := currentUser(r)
		year, week := yearWeek(r)

		s, err := database.Standings(year, week)
		if err != nil {
			log.Fatal(err)
		}

		err = Fetch(templatesDir, "standings.html").Execute(w, u, a, s)
		if err != nil {
			log.Println(err)
		}
	}
}
