package main

import (
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

	//	router.Handle("/results/{year:[0-9]{4}}/{week:[0-9]{2}}", Protect(Results)).Methods("GET").Name("Results")
}

func main() {
	log.Printf("NFL Pick-Em Pool listening on port 61389")
	log.Fatal(http.ListenAndServe(":61389", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	t := Fetch("index.html")

	session, _ := store.Get(r, "LoginState")
	user := session.Values["user"]

	if user == nil || user == "" {
		t.Execute(w, "", "Random Person")
	} else {
		t.Execute(w, user.(string), fmt.Sprintf("%s", user))
	}
}
