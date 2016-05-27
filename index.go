package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ameske/nfl-pickem/database"
)

func Index(templatesDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		u, a := currentUser(r)

		// If we aren't in the offseason and a user is logged in,
		//redirect to the Standings Page
		_, _, err := database.CurrentWeek(time.Now())
		if err == nil && u != "" {
			CurrentStandings(templatesDir)(w, r)
			return
		} else if err != database.ErrOffseason && err != nil {
			log.Fatal(err)
		}

		welcome := struct {
			User string
		}{u}

		err = Fetch(templatesDir, "index.html").Execute(w, u, a, welcome)
		if err != nil {
			log.Println(err)
		}
	}
}
