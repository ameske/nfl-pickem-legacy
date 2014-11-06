package main

import (
	"net/http"

	"github.com/ameske/go_nfl/database"
)

func Stadings(w http.ResponseWriter, r *http.Request) {
	standings := database.GetStandingsForm(db)
	Fetch("standings.html").Execute(w, standings)
}
