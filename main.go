package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	router.HandleFunc("/login", Login).Methods("POST")
	router.Handle("/logout", Protect(Logout)).Methods("POST")

	router.Handle("/changePassword", Protect(ChangePasswordForm)).Methods("GET").Name("ChangePasswordForm")
	router.Handle("/changePassowrd", Protect(ChangePassword)).Methods("POST")

	router.Handle("/picks/{year:[0-9]*}/{week:[0-9]*}", Protect(PicksForm)).Methods("GET").Name("Picks")
	router.Handle("/picks/{year:[0-9]*}/{week:[0-9]*}", Protect(ProcessPicks)).Methods("POST")

	router.HandleFunc("/results/{year:[0-9]*}/{week:[0-9]*}", Results).Methods("GET").Name("Results")
}

func main() {
	log.Printf("NFL Pick-Em Pool listening on port 61389")
	log.Fatal(http.ListenAndServe(":61389", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	Fetch("index.html").Execute(w, u, u)
}

func Results(w http.ResponseWriter, r *http.Request) {
	year, week := yearWeek(r)
	results := fmt.Sprintf("%d-Week%d-Results.html", year, week)
	u := currentUser(r)
	Fetch(results).Execute(w, u, nil)
}

func currentUser(r *http.Request) string {
	session, _ := store.Get(r, "LoginState")
	user := session.Values["user"]

	if user == nil || user == "" {
		return ""
	} else {
		return user.(string)
	}
}

func yearWeek(r *http.Request) (int, int) {
	v := mux.Vars(r)
	y, _ := strconv.ParseInt(v["year"], 10, 32)
	w, _ := strconv.ParseInt(v["week"], 10, 32)
	return int(y), int(w)
}
