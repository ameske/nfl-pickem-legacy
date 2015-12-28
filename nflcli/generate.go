package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/ameske/nfl-pickem/database"
)

func Generate(args []string) {
	if len(args) == 0 {
		GenerateHelp()
		return
	}

	switch args[0] {
	case "picks":
		GenerateSeasonPicks(args[1:])
	case "help":
		GenerateHelp()
	default:
		GenerateHelp()
	}

}

func GenerateHelp() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, "nfl-generate creates static assets or inits database tables")
	fmt.Fprintln(w, "\nAvailable commands:")
	fmt.Fprintln(w, "\tpicks\t Generate empty pick rows for a sesason")
	fmt.Fprintf(w, "\n")

	err := w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

// SeasonPicks generates empty pick rows for a year's games. The games must be
// loaded into the database before this can be called.
func GenerateSeasonPicks(args []string) {
	var year int

	f := flag.NewFlagSet("generateSeasonPicks", flag.ExitOnError)
	f.IntVar(&year, "year", -1, "Year")

	err := f.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	if year == -1 {
		log.Fatalf("Year is a required argument")
	}

	database.CreateSeasonPicks(year)
}
