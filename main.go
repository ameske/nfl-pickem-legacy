package main

import (
	"flag"
	"fmt"
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
	config = Config{}
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
	configPath := flag.String("config", "", "Location of config file")
	flag.Parse()

	if *configPath == "" {
		log.Fatalf("Config path must be specified")
		flag.PrintDefaults()
	}

	if err := LoadConfig(*configPath); err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("NFL App Setting (Email): %s", config.EmailAddress)
	log.Printf("NFL App Setting (Password): %s", config.Password)
	log.Printf("NFL App Setting (SMTP Address): %s", config.SMTPAddress)
	log.Printf("NFL App Setting (SMTP Port): %s", config.SMTPPort)

	log.Printf("NFL Pick-Em Pool listening on port 61389")
	log.Fatal(http.ListenAndServe(":61389", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	s := database.Standings(db, 2014, 12)
	Fetch("index.html").Execute(w, u, s)
}

func Results(w http.ResponseWriter, r *http.Request) {
	year, week := yearWeek(r)
	u := currentUser(r)
	Fetch(fmt.Sprintf("%d-Week%d-Results.html", year, week)).Execute(w, u, nil)
}

func Standings(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	year, week := yearWeek(r)
	s := database.Standings(db, year, week)
	Fetch("standings.html").Execute(w, u, s)
}
