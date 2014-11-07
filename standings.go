package main

import (
	"net/http"

	"github.com/ameske/go_nfl/database"
)

func Stadings(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "LoginState")
	user := session.Values["user"]

	standings := database.GetStandingsForm(db)

	if user == nil || user == "" {
		Fetch("standings.html").Execute(w, "", standings)
	} else {
		Fetch("standings.html").Execute(w, user.(string), standings)
	}
}
