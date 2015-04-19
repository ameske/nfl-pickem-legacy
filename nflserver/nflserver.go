package main

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ameske/nfl-pickem/database"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// For now, we will let all of these things be global since it's easier
var (
	store  *sessions.CookieStore
	router = mux.NewRouter()
)

type Config struct {
	Server ServerConfig
	Email  EmailConfig
}

type ServerConfig struct {
	AuthKey            string `json:"authKey"`
	EncryptKey         string `json:"encryptKey"`
	PostgresConnString string `json:"postgresConnString"`
}

type EmailConfig struct {
	SendAsAddress   string `json:"sendAsAddress"`
	Password        string `json:"password"`
	SMTPAddress     string `json:"smtpAddress"`
	SMTPFullAddress string `json:"smtpFullAddress"`
}

func init() {
	// App configuration
	config := loadConfig("/opt/ameske/gonfl/conf.json")

	database.SetDefaultDb(config.Server.PostgresConnString)

	configureSessionStore(config)
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

func loadConfig(path string) Config {
	configBytes, err := ioutil.ReadFile(path)

	config := Config{}
	err = json.Unmarshal(configBytes, &config)

	if err != nil {
		log.Fatalf(err.Error())
	}

	return config
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
