package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"

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
	AuthKey            string `json:"authKey"`
	EncryptKey         string `json:"encryptKey"`
	DatabaseFile       string `json:"databaseFile"`
	TemplatesDirectory string `json:"templatesDirectory"`
	LogosDirectory     string `json:"logosDirectory"`
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
		templatesDir = config.Server.TemplatesDirectory
		logosDir = config.Server.LogosDirectory
	} else {
		store = sessions.NewCookieStore([]byte("something secret"), []byte("something secret"))
		err := database.SetDefaultDb("nfl.db")
		if err != nil {
			log.Fatal(err)
		}
		templatesDir = "/Users/ameske/Documents/go/src/github.com/ameske/nfl-pickem/templates/"
		logosDir = "/Users/ameske/Documents/go/src/github.com/ameske/nfl-pickem/logos"
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	scheduleUpdates()

	// TODO - Get this served by something static?
	router.HandleFunc("/logo/{team:.*}", Logo)

	router.HandleFunc("/", Index)
	router.HandleFunc("/login", Login)
	router.HandleFunc("/logout", Protect(Logout))
	router.HandleFunc("/changePassword", Protect(ChangePassword))
	router.HandleFunc("/picks", Protect(Picks))
	router.HandleFunc("/admin/{year:[0-9]*}/{week:[0-9]*}", Protect(AdminOnly(AdminPickForm)))
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
		log.Fatal(err)
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

// scheduleUpdates sets up goroutines that will import the results of games and update the
// picks after every wave of games completes.
func scheduleUpdates() {
	// Friday at 8:00
	go func() {
		nextFriday := adjustIfPast(nextDay(time.Friday).Add(time.Hour * 8))
		time.Sleep(nextFriday.Sub(time.Now()))
		update(false)
	}()

	// Sunday at 18:00
	go func() {
		nextSunday := adjustIfPast(nextDay(time.Sunday).Add(time.Hour * 18))
		time.Sleep(nextSunday.Sub(time.Now()))
		update(false)
	}()

	// Sunday at 21:00
	go func() {
		nextSunday := adjustIfPast(nextDay(time.Sunday).Add(time.Hour * 21))
		time.Sleep(nextSunday.Sub(time.Now()))
		update(false)
	}()

	// Monday at 8:00
	go func() {
		nextMonday := adjustIfPast(nextDay(time.Monday).Add(time.Hour * 8))
		time.Sleep(nextMonday.Sub(time.Now()))
		update(false)
	}()

	// Tuesday at 8:00. Here we need to update the current week - 1
	go func() {
		nextTuesday := adjustIfPast(nextDay(time.Tuesday).Add(time.Hour * 8))
		time.Sleep(nextTuesday.Sub(time.Now()))
		update(true)
	}()
}

func update(updatePreviousWeek bool) {
	for {
		year, week := database.CurrentWeek()
		if updatePreviousWeek {
			week -= 1
		}
		UpdateGameScores(year, week)
		Grade(year, week)
		time.Sleep(time.Hour * 24 * 7)
	}
}

func adjustIfPast(next time.Time) time.Time {
	now := time.Now()
	ty, tm, td := now.Date()
	ny, nm, nd := next.Date()

	// if the next day is today, but the hour we want has past then we must advance next a week
	if (ty == ny && tm == nm && td == nd) && now.Hour() > next.Hour() {
		next = next.AddDate(0, 0, 7)
	}

	return next
}

func nextDay(day time.Weekday) time.Time {
	now := time.Now()

	// We only want to go forwards, so use modular arith to force going ahead
	diff := int(day-now.Weekday()+7) % 7

	next := now.AddDate(0, 0, diff)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())

	return next
}
