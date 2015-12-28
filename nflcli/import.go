package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"text/tabwriter"

	"github.com/ameske/nfl-pickem/database"
)

// Import is the dispatch hub for the "import" subcommand
func Import(args []string) {

	if len(args) == 0 {
		ImportHelp()
		return
	}

	switch args[0] {
	case "scores":
		ImportScores(args[1:])
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
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	var results []ResultsJson
	err = json.NewDecoder(pipe).Decode(&results)
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		err := database.UpdateScores(result.Week, result.Year, result.Home, result.HomeScore, result.AwayScore)
		if err != nil {
			log.Fatal(err)
		}
	}
}
