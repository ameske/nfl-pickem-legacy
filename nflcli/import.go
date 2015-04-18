package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ameske/nfl-pickem/database"
	"github.com/codegangsta/cli"
	_ "github.com/lib/pq"
)

func Import(args []string) {
	log.Println("Reached the import subcommand")

	if len(args) == 0 {
		ImportHelp()
		return
	}

	switch args[0] {
	case "schedule":
		log.Println("Reached the import-schedule subcommand")
		//ImportSchedule(args[1:])
	case "scores":
		log.Println("Reached the import-scores subcommand")
		//ImportScores(Args[1:])
	case "teams":
		log.Println("Reached the import-teams subcommand")
		//ImportTeams(Args[1:])
	case "help":
		ImportHelp()
	default:
		ImportHelp()
	}
}

func ImportHelp() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, "nfl import - Import assets into the database")
	fmt.Fprintln(w, "\nAvailable Commands:")
	fmt.Fprintln(w, "\tschedule\t Scrape nfl.com for a year's schedule")
	fmt.Fprintln(w, "\tscores\t Scrape nfl.com for a week's scores")
	fmt.Fprintln(w, "\tteams\t Import teams into the database for setup")
	fmt.Fprintln(w, "\thelp\t Display this message\n")

	err := w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

type ResultsJson struct {
	Week      int    `json:"week"`
	Year      int    `json:"year"`
	Home      string `json:"home"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
}

func ImportScores(args []string) {
	var year, week int

	f := flag.NewFlagSet("import scores", flag.ExitOnError)
	f.IntVar(&year, "year", -1, "Year")
	f.IntVar(&week, "week", -1, "Week")

	err := f.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	if year == -1 || week == -1 {
		year, week = database.CurrentWeek(db)
	}

	cmd := exec.Command("weeklyScores", strconv.Itoa(year), strconv.Itoa(week))
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Pipe: %s", err.Error())
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Start: %s", err.Error())
	}

	var results []ResultsJson
	err = json.NewDecoder(pipe).Decode(&results)
	if err != nil {
		log.Fatalf("Decode: %s", err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Wait: %s", err.Error())
	}

	for _, result := range results {
		// Lookup the year ID
		var yearId int64
		err = db.SelectOne(&yearId, "SELECT id FROM years WHERE year = $1", result.Year)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		// Lookup the week ID
		var weekId int64
		err = db.SelectOne(&weekId, "SELECT id FROM weeks WHERE year_id = $1 AND week = $2", yearId, result.Week)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		// Lookup the home team ID
		var teamId int64
		err = db.SelectOne(&teamId, "SELECT id FROM teams WHERE nickname = $1", result.Home)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		// Lookup the game ID based on the home team and week
		var game database.Games
		err = db.SelectOne(&game, "SELECT * FROM Games WHERE week_id = $1 and home_id = $2", weekId, teamId)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		// Update the scores for that game
		game.HomeScore = result.HomeScore
		game.AwayScore = result.AwayScore

		count, err := db.Update(&game)
		if err != nil {
			log.Fatalf("%s", err.Error())
		} else if count != 1 {
			log.Fatalf("More than one game was updated by the update statement.")
		}
	}
}

type GameJson struct {
	Week       int       `json:"week"`
	Home       string    `json:"home"`
	Away       string    `json:"away"`
	DateString string    `json:"date"`
	Date       time.Time `json:"-"`
	Year       int       `json:"year"`
}

func ImportSchedule(args []string) {
	games := make([]*GameJson, 0)

	// Open the 2014 schedule, and parse the json into our golang struct
	bytes, err := ioutil.ReadFile("json/2014/2014-Schedule.json")
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = json.Unmarshal(bytes, &games)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Convert the time string into a time.Time for postgres
	for _, game := range games {
		estDateString := game.DateString + " EST"
		game.Date, err = time.Parse("2006-01-02T15:04:05 MST", estDateString)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	// Construct a mapping of team nicknames to ID's
	var teams []database.Teams
	_, err = db.Select(&teams, "SELECT * FROM teams")
	if err != nil {
		log.Fatalf(err.Error())
	}
	teamsMap := make(map[string]int64)
	for _, team := range teams {
		teamsMap[team.Nickname] = team.Id
	}

	// Now, construct a game row and add it into postgres
	for _, game := range games {
		// Note: Manually update year id!
		weekId, err := db.SelectInt("SELECT id FROM weeks WHERE week = $1 AND year_id = 1", game.Week)
		temp := database.Games{
			WeekId:    weekId,
			HomeId:    teamsMap[game.Home],
			AwayId:    teamsMap[game.Away],
			Date:      game.Date,
			HomeScore: -1,
			AwayScore: -1,
		}
		err = db.Insert(&temp)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}

type Teams struct {
	Id       int64  `db:"id"`
	City     string `db:"city"`
	Nickname string `db:"nickname"`
	Stadium  string `db:"stadium"`
}

func ImportTeams(c *cli.Context) {
	scanner := bufio.NewScanner(bytes.NewBufferString(teams))

	// Line by line, get the information and load it into the DB
	for scanner.Scan() {
		teamLine := scanner.Text()
		splitLine := strings.Split(teamLine, ",")
		newTeam := &Teams{
			City:     splitLine[0],
			Nickname: splitLine[1],
			Stadium:  splitLine[2],
		}

		err := db.Insert(newTeam)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}

var teams = `Buffalo,Bills,Ralph Wilson Stadium
Miami,Dolphins,Sun Life Stadium
New England,Patriots,Gilette Stadium
New York,Jets,MetLife Stadium
Baltimore,Ravens,M&T Bank Stadium
Cincinatti,Bengals,Paul Brown Stadium
Cleveland,Browns,First Energy Stadium
Pittsburgh,Steelers,Heinz Field
Houston,Texans,Reliant Stadium
Indianapolis,Colts,Lucas Oil Stadium
Jacksonville,Jaguars,EverBank Field
Tennessee,Titans,LP Field
Denver,Broncos,Mile High Stadium
Kansas City,Chiefs,Arrowhead Stadium
Oakland,Raiders,O.co Coliseum
San Diego,Chargers,Qualcomm Stadium
Dallas,Cowboys,AT&T Stadium
New York,Giants,MetLife Stadium
Philadelphia,Eagles,Lincoln Financial Field
Washington,Redskins,FedEx Field
Chicago,Bears,Soldier Field
Detroit,Lions,Ford Field
Green Bay,Packers,Lambeau Field
Minnesota,Vikings,TCF Bank Stadium
Atlanta,Falcons,Georiga Dome
Carolina,Panthers,Bank of America Stadium
New Orleans,Saints,Mercedes-Benz Superdome
Tampa Bay,Buccaneers,Raymond James Stadium
Arizona,Cardinals,University of Phoenix Stadium
St. Louis,Rams,Edward Jones Dome
San Francisco,49ers,Candlestick Park
Seattle,Seahawks,CenturyLink Field`
