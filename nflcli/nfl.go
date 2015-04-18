package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/coopernurse/gorp"
)

var (
	db *gorp.DbMap
)

func init() {
	//	db = database.NflDb("host=localhost port=5432")
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		Help()
		return
	}

	// Select the appropriate subcommand
	switch args[0] {
	case "add":
		Add(args[1:])
	case "import":
		Import(args[1:])
	case "grade":
		Grade(args[1:])
	case "generate":
		Generate(args[1:])
	case "help":
		Help()
	default:
		Help()
	}
}

func Help() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, "nfl is a tool for managing an nfl-pickem pool server")
	fmt.Fprintln(w, "\nAvailable commands:")
	fmt.Fprintln(w, "\tadd\t Manually add information to the database")
	fmt.Fprintln(w, "\timport\t Import information into the database")
	fmt.Fprintln(w, "\tgrade\t Grade a given week's picks")
	fmt.Fprintln(w, "\tgenerate\t Generate HTML based on information in the database")
	fmt.Fprintln(w, "\thelp\t Display this message")
	fmt.Fprintf(w, "\n")

	err := w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
