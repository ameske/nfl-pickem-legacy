package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ameske/go_nfl/database"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// For now, we will let all of these things be global since it's easier
var (
	store  = sessions.NewCookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
	db     = database.NflDb()
	router = mux.NewRouter()
)

func init() {
	router.HandleFunc("/", Index).Methods("GET").Name("Index")
	router.HandleFunc("/login", LoginForm).Methods("GET").Name("LoginForm")
	router.HandleFunc("/login", Login).Methods("POST").Name("Login")
	router.HandleFunc("/logout", Logout).Methods("POST").Name("Logout")
	router.Handle("/picks", Protect(Picks)).Methods("GET").Name("Picks")
}

func main() {
	log.Fatal(http.ListenAndServe(":61389", router))
}

func writeJsonResponse(w http.ResponseWriter, r interface{}) {
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(j)
}

func Index(w http.ResponseWriter, r *http.Request) {
	t := Fetch("index.html")
	t.Execute(w, nil)
}
