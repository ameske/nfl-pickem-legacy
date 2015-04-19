package main

import (
	"net/http"

	"github.com/ameske/nfl-pickem/database"
)

func Index(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	year, week := database.CurrentWeek()
	s := database.Standings(year, week)
	Fetch("index.html").Execute(w, u, s)
}

func Standings(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	year, week := yearWeek(r)
	s := database.Standings(year, week)
	Fetch("standings.html").Execute(w, u, s)
}
