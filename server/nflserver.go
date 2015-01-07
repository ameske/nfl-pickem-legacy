package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ameske/go_nfl/database"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gopkg.in/yaml.v2"
)

// For now, we will let all of these things be global since it's easier
var (
	store  *sessions.CookieStore
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
	if err := LoadEmailConfig("/opt/ameske/gonfl/go_nfl.yaml"); err != nil {
		log.Fatalf(err.Error())
	}

	if err := LoadWebappConfig("/opt/ameske/gonfl/nfl.yaml"); err != nil {
		log.Fatalf(err.Error())
	}

	log.Printf("NFL Pick-Em Pool listening on port 61389")
	log.Fatal(http.ListenAndServe(":61389", router))
}

func LoadWebappConfig(path string) error {
	configBytes, err := ioutil.ReadFile(path)

	keys := struct {
		AuthKey    string `yaml:"AuthKey"`
		EncryptKey string `yaml:"EncryptKey"`
	}{}

	err = yaml.Unmarshal(configBytes, &keys)
	if err != nil {
		return err
	}

	decodedAuth, err := base64.StdEncoding.DecodeString(keys.AuthKey)
	if err != nil {
		return err
	}

	decodedEncrypt, err := base64.StdEncoding.DecodeString(keys.EncryptKey)
	if err != nil {
		return err
	}

	store = sessions.NewCookieStore(decodedAuth, decodedEncrypt)

	return nil
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
	name := fmt.Sprintf("/opt/ameske/gonfl/templates/%d-Week%d-Results.html", year, week)
	t := template.New("_base.html").Funcs(funcs)
	t = template.Must(t.ParseFiles("/opt/ameske/gonfl/templates/_base.html", "/opt/ameske/gonfl/templates/navbar.html", name))
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
