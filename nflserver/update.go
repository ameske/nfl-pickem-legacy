package main

import (
	"encoding/json"
	"log"
	"math"
	"os/exec"
	"strconv"

	"github.com/ameske/nfl-pickem/database"
)

type ResultsJson struct {
	Week      int    `json:"week"`
	Year      int    `json:"year"`
	Home      string `json:"home"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
}

// ImportScores scrapes the NFL's website using a helper script and inserts those
// scores into the database.
func UpdateGameScores() {
	year, week := database.CurrentWeek()

	cmd := exec.Command("weeklyScores", strconv.Itoa(year), strconv.Itoa(week))
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	var results []ResultsJson
	err = json.NewDecoder(pipe).Decode(&results)
	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		err := database.UpdateScores(result.Week, result.Year, result.Home, result.HomeScore, result.AwayScore)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// Grade calculates the scores for each user in the database for the given week.
// It assumes that the scores for the graded week have already been imported, else
// results are undefined.
func Grade() {
	year, week := database.CurrentWeek()

	// Gather this week's games
	gamesSlice := database.WeeklyGames(year, week)
	gamesMap := database.GamesMap(gamesSlice)

	// Gather all of the user id's
	users := database.AllUsers()

	// For each user, score their picks for this week
	for _, u := range users {
		picks := database.WeeklyPicksYearWeek(u.Email, year, week)

		for _, p := range picks {
			// Ignore all games that haven't finished yet - clean up points though
			if gamesMap[p.GameId].HomeScore == -1 && gamesMap[p.GameId].AwayScore == -1 {
				p.Correct = false
				continue
			}

			if gamesMap[p.GameId].HomeScore == gamesMap[p.GameId].AwayScore {
				p.Correct = true
				p.Points = int(math.Floor(float64(p.Points) / 2))
			} else if gamesMap[p.GameId].HomeScore > gamesMap[p.GameId].AwayScore && p.Selection == 2 {
				p.Correct = true
			} else if gamesMap[p.GameId].HomeScore > gamesMap[p.GameId].AwayScore && p.Selection == 1 {
				p.Correct = false
			} else if gamesMap[p.GameId].AwayScore > gamesMap[p.GameId].HomeScore && p.Selection == 2 {
				p.Correct = false
			} else {
				p.Correct = true
			}

			err := database.UpdatePick(p.Id, p.Correct)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
