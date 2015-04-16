package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ameske/nfl-pickem/database"
	"github.com/codegangsta/cli"
)

type PickSelection struct {
	database.Picks
	AwayId int64
	HomeId int64
}

func inputPicks(c *cli.Context) {
	user, year, week := c.String("user"), c.Int("year"), c.Int("week")
	// Make sure that the user put something for all of the flags
	if user == "" || year == -1 || week == -1 {
		fmt.Fprintf(os.Stderr, "Please explicitly supply all flags.\nUsage of %s\n", os.Args[0])
		flag.PrintDefaults()
		return
	}

	// Go ahead and grab a mapping of team ID's to names for the CLI
	teams := database.TeamMap(db)

	// Look up the userId for the user specified
	fmt.Printf("Looking up user ID...")
	var userId int64
	err := db.SelectOne(&userId, "SELECT id from users WHERE email = $1", user)
	if err != nil {
		log.Fatalf("UserId: %s", err.Error())
	}
	fmt.Printf("%d\n", userId)

	// Gather the set of picks for the request year/week
	fmt.Printf("Gathering picks to be made...")
	var picks []PickSelection
	weekId := database.WeekId(db, year, week)
	_, err = db.Select(&picks, "SELECT picks.*, games.away_id AS AwayId, games.home_id AS HomeId FROM picks INNER JOIN games ON picks.game_id = games.id WHERE picks.user_id = $1 AND games.week_id = $2", userId, weekId)
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
					fmt.Printf("Points must be the value 1, 3, 5, or 7.\n")
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

	var enteredPicks []PickSelection
	_, err = db.Select(&enteredPicks, "SELECT picks.*, games.away_id AS AwayId, games.home_id AS HomeId FROM picks INNER JOIN games ON picks.game_id = games.id WHERE picks.user_id = $1 AND games.week_id = $2", userId, weekId)
	if err != nil {
		log.Fatalf("PickSelection Verify: %s", err.Error())
	}

	fmt.Printf("Picks successfully entered. Please verify.\n")
	for _, p := range enteredPicks {
		fmt.Printf("Game: (1) %s at (2) %s\n", teams[p.AwayId], teams[p.HomeId])
		fmt.Printf("Selection: %d \t Points: %d\n\n", p.Picks.Selection, p.Picks.Points)
	}
}
