package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ameske/nfl-pickem/database"
)

func Results(templatesDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		year, week := yearWeek(r)
		if year < 2015 || week <= 0 || week >= 17 {
			http.Error(w, "invalid year/week", http.StatusBadRequest)
			return
		}

		u, a := currentUser(r)
		if u == "" {
			http.Error(w, "no user information", http.StatusUnauthorized)
			return
		}

		data, err := GenerateResultsData(year, week)
		if err != nil {
			log.Fatal(err)
		}

		err = Fetch(templatesDir, "results.html").Execute(w, u, a, data)
		if err != nil {
			log.Println(err)
		}
	}
}

type ResultsTemplateData struct {
	Users  []string
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
func GenerateResultsData(year int, week int) (*ResultsTemplateData, error) {
	users, err := database.Usernames()
	if err != nil {
		return nil, err
	}

	games, err := database.WeeklyGames(year, week)
	if err != nil {
		return nil, err
	}

	picks := make([][]database.Picks, len(users))
	for i, u := range users {
		picks[i], err = database.UserPicksByWeek(u, year, week)
		if err != nil {
			return nil, err
		}
	}

	// Build each row of the table, where each row represents one game and all of the user's picks for that game
	rows := make([]ResultsTableRow, len(games))
	for i, g := range games {

		if time.Now().Before(g.Date) {
			continue
		}

		tr := ResultsTableRow{
			Matchup: fmt.Sprintf("%s/%s", g.AwayAbbreviation, g.HomeAbbreviation),
			Picks:   make([]UserPicks, len(users)),
		}

		// Each element of picks is a slice containing all the game picks for a user.
		// games and these slices are in the same order. So we just deal with the i'th pick each time, but put it in the j'th slot of the UserPick.
		for j, p := range picks {
			switch p[i].Selection {
			case 1:
				tr.Picks[j].Pick = g.AwayAbbreviation
			case 2:
				tr.Picks[j].Pick = g.HomeAbbreviation
			}
			tr.Picks[j].Points = p[i].Points

			if g.HomeScore == -1 && g.AwayScore == -1 {
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

	userNames, err := database.UserFirstNames()
	if err != nil {
		return nil, err
	}

	data := &ResultsTemplateData{
		Users:  userNames,
		Rows:   rows,
		Totals: totals,
		Year:   year,
		Week:   week,
	}

	return data, nil
}
