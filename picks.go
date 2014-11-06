package main

import (
	"log"
	"net/http"

	"github.com/ameske/go_nfl/database"
)

// Picks fetches this week's picks for the current logged in user and renders the
// picks template with them.
func Picks(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "LoginState")
	user := session.Values["user"].(string)

	picks := database.FormPicks(db, user, 2014, 10)
	log.Printf("Pulled %d picks for user %s", len(picks), user)
	Fetch("picks.html").Execute(w, picks)
}
