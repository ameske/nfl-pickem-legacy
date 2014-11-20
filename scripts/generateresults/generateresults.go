package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/ameske/go_nfl/database"
)

type ResultsTemplateData struct {
	Users   []database.Users
	Rows    []ResultsTableRow
	Totals  []int
	Title   string
	Content string
	End     string
}

type ResultsTableRow struct {
	Matchup string
	Picks   []UserPicks
}

type UserPicks struct {
	Pick    string
	Points  int
	Correct bool
}

var (
	week, year int
	debug      bool
)

/*
* Creates an HTML file based on a template that displays the results for a given week.
 */
func main() {
	flag.IntVar(&week, "week", -1, "Week to generate results html for")
	flag.IntVar(&year, "year", -1, "Week to generate results html for")
	flag.BoolVar(&debug, "debug", false, "Whether or not to print to console")
	flag.Parse()
	if week == -1 || year == -1 {
		fmt.Printf("Week and year must be provided explicitly.\n\n")
		flag.PrintDefaults()
	}

	db := database.NflDb()

	teams := database.TeamAbbreviationMap(db)

	users := database.AllUsers(db)
	games := database.WeeklyGames(db, year, week)
	picks := make([][]*database.Picks, len(users))
	for i, u := range users {
		picks[i] = database.WeeklyPicks(db, u.Id, year, week)
		reorderPicks(games, picks[i])
	}

	if debug {
		fmt.Printf("[ ")
		for _, g := range games {
			fmt.Printf("%d ", g.Id)
		}
		fmt.Printf("]\n")

		for _, u := range picks {
			fmt.Printf("[ ")
			for _, p := range u {
				fmt.Printf("%d ", p.GameId)
			}
			fmt.Printf("]\n")
		}
		return
	}

	// Build each row of the table, where each row represents one game and all of the user's picks for that game
	rows := make([]ResultsTableRow, len(games))
	for i, g := range games {
		tr := ResultsTableRow{}

		tr.Matchup = fmt.Sprintf("%s/%s", teams[g.AwayId], teams[g.HomeId])

		tr.Picks = make([]UserPicks, len(users))
		for j, p := range picks {
			switch p[i].Selection {
			case 1:
				tr.Picks[j].Pick = fmt.Sprintf("%s", teams[games[i].AwayId])
			case 2:
				tr.Picks[j].Pick = fmt.Sprintf("%s", teams[games[i].HomeId])
			}
			tr.Picks[j].Points = p[i].Points
			tr.Picks[j].Correct = p[i].Correct
		}
		rows[i] = tr
	}

	// Get the user's total for the week
	totals := make([]int, len(users))
	for i, user := range picks { //each user
		total := 0
		for _, p := range user { //each pick for the user
			if p.Correct {
				total += p.Points
			}
		}
		totals[i] = total
	}

	data := ResultsTemplateData{}
	data.Users = users
	data.Rows = rows
	data.Totals = totals
	data.Title = fmt.Sprintf("{{define \"title\"}}%d - Week %d Results{{end}}", year, week)
	data.Content = "{{define \"content\"}}"
	data.End = "{{end}}"

	t := template.New("results.html")
	t = template.Must(t.ParseFiles("results.html"))

	weekResults, err := os.Create(fmt.Sprintf("%d-Week%d-Results.html", year, week))
	if err != nil {
		log.Fatalf("CreatingFile: %s", err.Error())
	}
	defer weekResults.Close()

	err = t.Execute(weekResults, data)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

}

// Uses a stupid naive and slow algorithm to make sure that the picks line up with the games down the side.
// Fuck you postgres.
func reorderPicks(gamesOrder []database.Games, picks []*database.Picks) {
	for i, g := range gamesOrder {
		if picks[i].GameId != g.Id {
			// if it doesn't match go find the one that should be here and swap them
			for j := i; j < len(picks); j++ {
				if picks[j].GameId == g.Id {
					picks[i], picks[j] = picks[j], picks[i]
					break
				}
			}

		}
	}
}
