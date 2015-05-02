package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/ameske/nfl-pickem/database"
)

func TestData(args []string) {
	if len(args) == 0 {
		TestDataHelp()
		return
	}

	switch args[0] {
	case "new":
		NewTestDatabase()
	default:
		TestDataHelp()
	}
}

func TestDataHelp() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, "nfl-testdata generates fake data for testing")
	fmt.Fprintln(w, "\nAvailable commands:")
	fmt.Fprintln(w, "\tnew\t Creates a new testdatabase")
	fmt.Fprintln(w, "\n")

	err := w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

func NewTestDatabase() {
	// Import Teams
	ImportTeams(nil)

	// Create Test User
	err := database.NewUser("test", "user", "testuser@gmail.com", "password")
	if err != nil {
		log.Fatal(err)
	}

	// Create a Week and Year
	err = database.NewPvs(8, 5, 2, 1, "A")
	if err != nil {
		log.Fatal(err)
	}
	err = database.NewYear(2015)
	if err != nil {
		log.Fatal(err)
	}
	err = database.NewWeek(1, time.Now().Unix(), 2015, "A")
	if err != nil {
		log.Fatal(err)
	}

	// Create Fake Games
	err = database.CreateRandomGames(1, 2015)
	if err != nil {
		log.Fatal(err)
	}

	// Create Fake Picks
	database.CreateSeasonPicks(2015)

	log.Println("Fake database successfully created")
}
