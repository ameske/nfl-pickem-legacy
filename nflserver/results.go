package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ameske/nfl-pickem/database"
)

func Results(w http.ResponseWriter, r *http.Request) {
	year, week := yearWeek(r)
	u, a := currentUser(r)

	data := GenerateResultsData(year, week)

	err := Fetch("results.html").Execute(w, u, a, data)
	if err != nil {
		log.Println(err)
	}
}

type ResultsTemplateData struct {
	Users  []database.Users
	Rows   []ResultsTableRow
	Totals []int
	Year   int
	Week   int
}

type ResultsTableRow struct {
	Matchup string
	Picks   []UserPicks
}

type UserPicks struct {
	Pick   string
	Points int
	Status PickStatus
}

type PickStatus int

const (
	Correct PickStatus = iota
	Incorrect
	Pending
)

// GenerateResultsHTML creates an HTML file based on a template that displays the results for a given week.
func GenerateResultsData(year int, week int) ResultsTemplateData {
	teams := database.TeamAbbreviationMap()
	users := database.AllUsers()
	games := database.WeeklyGames(year, week)

	picks := make([][]*database.Picks, len(users))
	for i, u := range users {
		picks[i] = database.WeeklyPicksYearWeek(u.Email, year, week)
		reorderPicks(games, picks[i])
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

			if games[i].HomeScore == -1 && games[i].AwayScore == -1 {
				tr.Picks[j].Status = Pending
			} else if p[i].Correct {
				tr.Picks[j].Status = Correct
			} else {
				tr.Picks[j].Status = Incorrect
			}
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
	data.Year = year
	data.Week = week

	return data
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
