package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ameske/nfl-pickem/database"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// For now, we will let all of these things be global since it's easier
var (
	store              *sessions.CookieStore
	router             = mux.NewRouter()
	emailNotifications = false
)

type Config struct {
	Server ServerConfig
	Email  EmailConfig
}

type ServerConfig struct {
	AuthKey      string `json:"authKey"`
	EncryptKey   string `json:"encryptKey"`
	DatabaseFile string `json:"databaseFile"`
}

type EmailConfig struct {
	SendAsAddress   string `json:"sendAsAddress"`
	Password        string `json:"password"`
	SMTPAddress     string `json:"smtpAddress"`
	SMTPFullAddress string `json:"smtpFullAddress"`
}

func main() {
	configFile := flag.String("config", "/opt/ameske/gonfl/conf.json", "Path to server config file")
	debug := flag.Bool("debug", false, "run the server with debug configuration instead of a config file")
	flag.Parse()

	if !*debug {
		config := loadConfig(*configFile)
		configureEmail(config)
		configureSessionStore(config)
		database.SetDefaultDb(config.Server.DatabaseFile)
	} else {
		store = sessions.NewCookieStore([]byte("something secret"), []byte("something secret"))
		err := database.SetDefaultDb("nfl.db")
		if err != nil {
			log.Fatal(err)
		}
	}

	// HTTP Server configuration
	router.HandleFunc("/", Index)
	router.HandleFunc("/login", Login)
	router.Handle("/logout", Protect(Logout))
	router.Handle("/changePassword", Protect(ChangePassword))
	router.Handle("/picks", Protect(Picks))
	router.HandleFunc("/results/{year:[0-9]*}/{week:[0-9]*}", Results)
	router.HandleFunc("/standings/{year:[0-9]*}/{week:[0-9]*}", Standings)

	log.Printf("NFL Pick-Em Pool listening on port 61389")
	log.Fatal(http.ListenAndServe("0.0.0.0:61389", router))
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
