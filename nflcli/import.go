package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ameske/nfl-pickem/database"
)

// Import is the dispatch hub for the "import" subcommand
func Import(args []string) {

	if len(args) == 0 {
		ImportHelp()
		return
	}

	switch args[0] {
	case "schedule":
		ImportSchedule(args[1:])
	case "scores":
		ImportScores(args[1:])
	case "teams":
		ImportTeams(args[1:])
	case "help":
		ImportHelp()
	default:
		ImportHelp()
	}
}

// ImportHelp displays the help message for the "import" subcommand
func ImportHelp() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, "nfl import - Import assets into the database")
	fmt.Fprintln(w, "\nAvailable Commands:")
	fmt.Fprintln(w, "\tschedule\t Scrape nfl.com for a year's schedule")
	fmt.Fprintln(w, "\tscores\t Scrape nfl.com for a week's scores")
	fmt.Fprintln(w, "\tteams\t Import teams into the database for setup")
	fmt.Fprintln(w, "\thelp\t Display this message")
	fmt.Fprintf(w, "\n")

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

// ImportScores scrapes the NFL's website using a helper script and inserts those
// scores into the database.
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
		year, week = database.CurrentWeek()
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
		err := database.UpdateScores(result.Week, result.Year, result.Home, result.HomeScore, result.AwayScore)
		if err != nil {
			log.Fatalf("Update Scores: %v", err)
		}
	}
}

type GameJson struct {
	Week       int    `json:"week"`
	Home       string `json:"home"`
	Away       string `json:"away"`
	DateString string `json:"date"`
	Date       int64  `json:"-"`
	Year       int    `json:"year"`
}

// ImportSchedule scrapes the nfl.com website for the request year's schedule using a
// helper script and inserts the games into the database.
func ImportSchedule(args []string) {
	var year int
	var endweek int

	f := flag.NewFlagSet("import schedule", flag.ExitOnError)
	f.IntVar(&year, "year", -1, "Year")
	f.IntVar(&endweek, "endweek", 18, "End Week")

	err := f.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	if year == -1 {
		log.Fatal("year required for import schedule")
	}

	cmd := exec.Command("scrapeSchedule", strconv.Itoa(year), strconv.Itoa(endweek))
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Pipe: %s", err.Error())
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Start: %s", err.Error())
	}

	games := make([]*GameJson, 0)
	err = json.NewDecoder(pipe).Decode(&games)
	if err != nil {
		log.Fatalf("Decode: %s", err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Wait: %s", err.Error())
	}

	// Convert the time string into a time.Time for sqlite
	for _, game := range games {
		estDateString := game.DateString + " EST"
		t, err := time.Parse("2006-01-02T15:04:05 MST", estDateString)
		if err != nil {
			log.Fatalf(err.Error())
		}
		game.Date = t.Unix()
	}

	// Now, construct a game row and add it into postgres
	for _, game := range games {
		err = database.AddGame(game.Week, game.Year, game.Home, game.Away, game.Date)
		if err != nil {
			log.Fatal("AddGame: ", err)
		}
	}
}

func ImportTeams(args []string) {
	scanner := bufio.NewScanner(bytes.NewBufferString(teams))

	// Line by line, get the information and load it into the DB
	for scanner.Scan() {
		teamLine := scanner.Text()
		splitLine := strings.Split(teamLine, ",")
		err := database.AddTeam(splitLine[0], splitLine[1], splitLine[2])
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
