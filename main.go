package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// Hardcode login/cookie stuff for testing
const (
	username = "kyle"
	password = "password"
)

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

func main() {
	//db := database.NflDb()

	r := mux.NewRouter()
	r.HandleFunc("/", LoginGet).Methods("GET")
	r.HandleFunc("/", LoginPost).Methods("POST")
	r.HandleFunc("/state", State).Methods("GET")
	r.HandleFunc("/logout", Logout).Methods("POST")

	log.Fatal(http.ListenAndServe(":61389", r))
}

func writeJsonResponse(w http.ResponseWriter, r interface{}) {
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(j)
}
