package main

import (
	"flag"
	"log"

	"github.com/ameske/nfl-pickem/database"
)

func main() {
	year := flag.Int("year", -1, "year")
	db := flag.String("db", "", "database file")
	flag.Parse()

	if *year == -1 {
		log.Fatal("must provide year")
	}

	if *db == "" {
		log.Fatal("must provide db")
	}

	err := database.SetDefaultDb(*db)
	if err != nil {
		log.Fatal(err)
	}

	err = database.CreateSeasonPicks(*year)
	if err != nil {
		log.Fatal(err)
	}
}
