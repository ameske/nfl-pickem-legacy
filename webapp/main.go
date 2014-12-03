package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

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
	router.Handle("/logout", Protect(Logout)).Methods("GET")
	router.Handle("/changePassword", Protect(ChangePasswordForm)).Methods("GET").Name("ChangePasswordForm")
	router.Handle("/changePassword", Protect(ChangePassword)).Methods("POST")
	router.Handle("/picks", Protect(PicksForm)).Methods("GET").Name("Picks")
	router.Handle("/picks", Protect(ProcessPicks)).Methods("POST")
	router.HandleFunc("/results/{year:[0-9]*}/{week:[0-9]*}", Results).Methods("GET").Name("Results")
	router.HandleFunc("/standings/{year:[0-9]*}/{week:[0-9]*}", Standings).Methods("GET").Name("Standings")
}

func main() {
	if err := LoadEmailConfig("/opt/ameske/etc/go_nfl/go_nfl.yaml"); err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("NFL Pick-Em Pool listening on port 61389")
	log.Fatal(http.ListenAndServe(":61389", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	year, week := database.CurrentWeek(db)
	s := database.Standings(db, year, week)
	Fetch("index.html").Execute(w, u, s)
}

func Results(w http.ResponseWriter, r *http.Request) {
	year, week := yearWeek(r)
	u := currentUser(r)

	// To support this page updating throughout the day, this either needs to be
	// generated dynamically each time, or we will have to be rude and parse the
	// template each time. The previous model of serving static content just doesn't work.
	name := fmt.Sprintf("%d-Week%d-Results.html", year, week)
	t := template.New("_base.html").Funcs(funcs)
	t = template.Must(t.ParseFiles("templates/_base.html", "templates/navbar.html", filepath.Join("templates", name)))
	data := struct {
		User    string
		Content interface{}
	}{
		u,
		nil,
	}
	t.Execute(w, data)
}

func Standings(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	year, week := yearWeek(r)
	s := database.Standings(db, year, week)
	Fetch("standings.html").Execute(w, u, s)
}
