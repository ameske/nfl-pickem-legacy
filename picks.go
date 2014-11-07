package main

import (
	"log"
	"net/http"

	"github.com/ameske/go_nfl/database"
)

// Picks fetches this week's picks for the current logged in user and renders the
// picks template with them.
func PicksForm(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "LoginState")
	user := session.Values["user"].(string)

	picks := database.FormPicks(db, user, 2014, 10)
	log.Printf("Pulled %d picks for user %s", len(picks), user)
	Fetch("picks.html").Execute(w, user, picks)
}

// ProcessPicks processes a user's picks ensuring the following.
//	- The user has not entered more than one 7 point game
//	- The user has not entered more than two 5 point games
//	- The user has not entered more than the allowable 3 point games for the week
//
// This does NOT check if a user has entered something for all of the games, as a user
// might not want to enter all of their picks at once. Picks are locked after the start
// of the game, and any submissions following game time must be manually entered.
func ProcessPicks(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	for k, _ := range r.Form {
		log.Printf("%s: %s", k, r.FormValue(k))
	}

	w.Write([]byte("All done here!"))
}
