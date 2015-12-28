package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/ameske/nfl-pickem/database"
)

func Add(args []string) {
	if len(args) == 0 {
		AddHelp()
		return
	}

	switch args[0] {
	case "user":
		AddUser(args[1:])
	case "picks":
		AddPicks(args[1:])
	case "help":
		AddHelp()

	default:
		AddHelp()
	}

}

// AddHelp displays the help message for the "add" subcommand
func AddHelp() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	fmt.Fprintln(w, "nfl-add manually adds assets to the database")
	fmt.Fprintln(w, "\nAvailable commands:")
	fmt.Fprintln(w, "\tuser\t Add a new user to the database")
	fmt.Fprintln(w, "\tpicks\t Input a user's picks to the database")
	fmt.Fprintln(w, "\thelp\t Display this message")
	fmt.Fprintf(w, "\n")

	err := w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

// AddUser adds a new user to the database. The password supplied is encrypted
// using bcrypt before being added.
func AddUser(args []string) {
	var first, last, email, password string

	f := flag.NewFlagSet("adduser", flag.ExitOnError)
	f.StringVar(&first, "first", "", "First Name")
	f.StringVar(&last, "last", "", "Last Name")
	f.StringVar(&email, "email", "", "Email Address")
	f.StringVar(&password, "password", "", "Password")

	err := f.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	if first == "" {
		log.Fatal("First name is required. Use --first <firstname> to specify.")
	}
	if last == "" {
		log.Fatal("Last name is required. Use --last <lastname> to specify.")
	}
	if email == "" {
		log.Fatal("Email is required. Use --email <address> to specify.")
	}
	if password == "" {
		log.Fatal("Desired password required. Use --password <pass> to specify.")
	}

	newUser := database.Users{
		FirstName: first,
		LastName:  last,
		Email:     email,
		Password:  password,
	}

	err = database.AddUser(newUser)
	if err != nil {
		log.Fatal(err)
	}
}

// AddPicks walks the admin through manual addition of picks for the request user
// and week.
func AddPicks(args []string) {
	var user string
	var year, week int

	f := flag.NewFlagSet("picks", flag.ExitOnError)
	f.StringVar(&user, "user", "", "Username")
	f.IntVar(&year, "year", -1, "Year")
	f.IntVar(&week, "week", -1, "Week")

	err := f.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	// Check to see if required arguments were given
	if user == "" || year == -1 || week == -1 {
		fmt.Fprintf(os.Stderr, "Please explicitly supply all flags.\nUsage of %s\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	teams := database.TeamMap()

	picks := database.CLIPickSelections(user, year, week)
	// For each pick, display the pick and prompt for selection and point value
	for _, p := range picks {
		var selection, points int
		valid := false
		for !valid {
			fmt.Printf("Game: (1) %s at (2) %s\n", teams[p.AwayId], teams[p.HomeId])
			fmt.Printf("Selection:")
			fmt.Fscanf(os.Stdin, "%d", &selection)
			fmt.Printf("Points (1,3,5,7):")
			fmt.Fscanf(os.Stdin, "%d", &points)
			if selection != 1 && selection != 2 {
				fmt.Printf("Selection must be (1) or (2).\n")
				if points != 1 && points != 3 && points != 5 && points != 7 {
					fmt.Printf("Point values must be one of 1, 3, 5, or 7.\n")
				}
				fmt.Printf("Please re-enter the game.\n\n")
				continue
			}
			valid = true
		}

		err := database.MakePick(p.Id, selection, points)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Have the user verify the picks that were just made. We have already updated the database, but
	// since this is designed to be an admin only app we can get away with it.
	fmt.Printf("Picks successfully entered. Please verify.\n")
	enteredPicks := database.CLIPickSelections(user, year, week)
	for _, p := range enteredPicks {
		fmt.Printf("Game: (1) %s at (2) %s\n", teams[p.AwayId], teams[p.HomeId])
		fmt.Printf("Selection: %d \t Points: %d\n\n", p.Picks.Selection, p.Picks.Points)
	}
}
