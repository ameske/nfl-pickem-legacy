package main

import (
	"encoding/base64"
	"log"
	"net/http"

	"github.com/ameske/go_nfl/database"
	"github.com/coopernurse/gorp"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// For now, we will let all of these things be global since it's easier
var (
	store  *sessions.CookieStore
	db     *gorp.DbMap
	router = mux.NewRouter()
)

func init() {
	// App configuration
	config := loadConfig("/opt/ameske/gonfl/conf.json")
	configureSessionStore(config)
	configureDb(config)
	configureEmail(config)

	// HTTP Server configuration
	router.HandleFunc("/", Index)
	router.HandleFunc("/login", Login)
	router.Handle("/logout", Protect(Logout))
	router.Handle("/changePassword", Protect(ChangePassword))
	router.Handle("/picks", Protect(PicksForm))
	router.Handle("/picks", Protect(ProcessPicks))
	router.HandleFunc("/results/{year:[0-9]*}/{week:[0-9]*}", Results)
	router.HandleFunc("/standings/{year:[0-9]*}/{week:[0-9]*}", Standings)
}

func main() {
	log.Printf("NFL Pick-Em Pool listening on port 61389")
	log.Fatal(http.ListenAndServe(":61389", router))
}

func configureDb(config Config) {
	db = database.NflDb(config.PostgresPort)
}

func configureSessionStore(config Config) {
	decodedAuth, err := base64.StdEncoding.DecodeString(config.Server.AuthKey)
	if err != nil {
		log.Fatalf(err.Error())
	}

	decodedEncrypt, err := base64.StdEncoding.DecodeString(config.Server.EncryptKey)
	if err != nil {
		log.Fatalf(err.Error())
	}

	store = sessions.NewCookieStore(decodedAuth, decodedEncrypt)
}
