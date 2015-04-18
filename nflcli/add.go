package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/ameske/nfl-pickem/database"
	"golang.org/x/crypto/bcrypt"
)

// Add is the subcommand hub for the "add" actions
func Add(args []string) {
	log.Println("Reached the add subcommand")

	if len(args) == 0 {
		AddHelp()
		return
	}

	switch args[0] {
	case "user":
		log.Println("Reached the add-user subcommand")
		//AddUser(args[1:])
	case "picks":
		log.Println("Reached the add-picks subcommand")
		//AddPicks(args[1:])
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

	// Check to see if required arguments were given
	if first == "" {
		log.Fatalf("First name is required. Use --first <firstname> to specify.")
	}
	if last == "" {
		log.Fatalf("Last name is required. Use --last <lastname> to specify.")
	}
	if email == "" {
		log.Fatalf("Email is required. Use --email <address> to specify.")
	}
	if password == "" {
		log.Fatalf("Desired password required. Use --password <pass> to specify.")
	}

	bpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf(err.Error())
	}

	newUser := &database.Users{
		FirstName: first,
		LastName:  last,
		Email:     email,
		Password:  string(bpass),
	}

	err = db.Insert(newUser)
	if err != nil {
		log.Fatalf(err.Error())
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

	teams := database.TeamMap(db)

	var userID int64
	err = db.SelectOne(&userID, "SELECT id from users WHERE email = $1", user)
	if err != nil {
		log.Fatalf("UserId: %s", err.Error())
	}
	fmt.Printf("%d\n", userID)

	type PickSelection struct {
		database.Picks
		AwayId int64
		HomeId int64
	}

	var picks []PickSelection
	weekID := database.WeekId(db, year, week)
	_, err = db.Select(&picks, `SELECT picks.*, games.away_id AS AwayId, games.home_id AS HomeId 
				    FROM picks INNER JOIN games ON picks.game_id = games.id
				    WHERE picks.user_id = $1 AND games.week_id = $2"`, userID, weekID)
	if err != nil {
		log.Fatalf("PickSelection: %s", err.Error())
	}

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

		p.Picks.Selection = selection
		p.Picks.Points = points
		_, err := db.Update(&p.Picks)
		if err != nil {
			log.Fatalf("Updating Pick: %s", err.Error())
		}
	}

	// Have the user verify the picks that were just made. We have already updated the database, but
	// since this is designed to be an admin only app we can get away with it.
	var enteredPicks []PickSelection
	_, err = db.Select(&enteredPicks, `SELECT picks.*, games.away_id AS AwayId, games.home_id AS HomeId
					   FROM picks INNER JOIN games ON picks.game_id = games.id
					   WHERE picks.user_id = $1 AND games.week_id = $2"`, userID, weekID)
	if err != nil {
		log.Fatalf("PickSelection Verify: %s", err.Error())
	}

	fmt.Printf("Picks successfully entered. Please verify.\n")
	for _, p := range enteredPicks {
		fmt.Printf("Game: (1) %s at (2) %s\n", teams[p.AwayId], teams[p.HomeId])
		fmt.Printf("Selection: %d \t Points: %d\n\n", p.Picks.Selection, p.Picks.Points)
	}
}
