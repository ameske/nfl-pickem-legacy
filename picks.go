package main

import (
	"net/http"

	"github.com/gorilla/context"
)

// Picks fetches this week's picks for the current logged in user and renders the
// picks template with them.
func Picks(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(string)

	// STUB
	w.Write([]byte(user))
}
