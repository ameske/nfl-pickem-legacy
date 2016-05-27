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

func WeekByWeekStandings(templatesDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, a := currentUser(r)

		year, week := yearWeek(r)

		weeks := make([]int, 0, week)
		for i := 1; i <= week; i++ {
			weeks = append(weeks, i)
		}

		s, err := database.WeekByWeekStandings(year, week)
		if err != nil {
			log.Fatal(err)
		}

		tmp := struct {
			CurrentWeek int
			Weeks       []int
			Users       []database.WeekByWeekStandingsPage
		}{CurrentWeek: week, Weeks: weeks, Users: s}

		err = Fetch(templatesDir, "weeklyStandings.html").Execute(w, u, a, tmp)
		if err != nil {
			log.Println(err)
		}
	}
}
