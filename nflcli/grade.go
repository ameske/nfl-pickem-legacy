package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/ameske/nfl-pickem/database"
)

// Grade calculates the scores for each user in the database for the given week.
// It assumes that the scores for the graded week have already been imported, else
// results are undefined.
func Grade(args []string) {
	var year, week int

	f := flag.NewFlagSet("grade", flag.ExitOnError)
	f.IntVar(&year, "year", -1, "Year")
	f.IntVar(&week, "week", -1, "Week")

	err := f.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	if year == -1 || week == -1 {
		year, week = database.CurrentWeek()
	}

	// Gather this week's games
	gamesSlice := database.WeeklyGames(year, week)
	gamesMap := database.GamesMap(gamesSlice)

	// Gather all of the user id's
	users := database.AllUsers()

	// For each user, score their picks for this week and print their total
	for _, u := range users {
		picks := database.WeeklyPicksYearWeek(u.Email, year, week)

		total := 0
		for _, p := range picks {
			// Ignore all games that haven't finished yet - clean up points though
			if gamesMap[p.GameId].HomeScore == -1 && gamesMap[p.GameId].AwayScore == -1 {
				p.Correct = false
				continue
			}

			if gamesMap[p.GameId].HomeScore == gamesMap[p.GameId].AwayScore {
				p.Correct = true
				p.Points = int(math.Floor(float64(p.Points) / 2))
				total += p.Points
			} else if gamesMap[p.GameId].HomeScore > gamesMap[p.GameId].AwayScore && p.Selection == 2 {
				p.Correct = true
				total += p.Points
			} else if gamesMap[p.GameId].HomeScore > gamesMap[p.GameId].AwayScore && p.Selection == 1 {
				p.Correct = false
			} else if gamesMap[p.GameId].AwayScore > gamesMap[p.GameId].HomeScore && p.Selection == 2 {
				p.Correct = false
			} else {
				p.Correct = true
				total += p.Points
			}

			err := database.UpdatePick(p.Id, p.Correct)
			if err != nil {
				log.Fatalf("Update: %s", err.Error())
			}
		}
		fmt.Printf("%s: %d\n", u.FirstName, total)
	}
}
