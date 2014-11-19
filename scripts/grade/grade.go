package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/ameske/go_nfl/database"
)

var (
	year, week int
)

func init() {
	flag.IntVar(&year, "year", -1, "Year")
	flag.IntVar(&week, "week", -1, "Week")
}

func main() {
	flag.Parse()
	db := database.NflDb()

	if year == -1 || week == -1 {
		log.Fatalf("Specify both --week and --year")
	}

	weekId := database.WeekId(db, year, week)

	// Gather this week's games
	var gamesSlice []database.Games
	_, err := db.Select(&gamesSlice, "SELECT * FROM games WHERE week_id = $1", weekId)
	if err != nil {
		log.Fatalf("Games: %s", err.Error())
	}
	gamesMap := database.GamesMap(gamesSlice)

	// Gather all of the user id's
	var usersSlice []database.Users
	_, err = db.Select(&usersSlice, "SELECT * FROM users")
	if err != nil {
		log.Fatalf("Users: %s", err.Error())
	}
	usersMap := database.UsersMap(usersSlice)

	// For each user, score their picks for this week and print their total
	for u, _ := range usersMap {
		var picks []database.Picks
		_, err := db.Select(&picks, "SELECT picks.* FROM picks INNER JOIN games ON picks.game_id = games.id WHERE games.week_id = $1 AND picks.user_id = $2", weekId, u)
		if err != nil {
			log.Fatalf("Picks: %s", err.Error())
		}

		total := 0
		for _, p := range picks {
			if gamesMap[p.GameId].HomeScore == gamesMap[p.GameId].AwayScore {
				p.Correct = true
				p.Points = int(math.Floor(float64(p.Points) / 2))
				total += p.Points
			} else if gamesMap[p.GameId].HomeScore > gamesMap[p.GameId].AwayScore && p.Selection == 2 {
				p.Correct = true
				total += p.Points
			} else if gamesMap[p.GameId].HomeScore > gamesMap[p.GameId].AwayScore && p.Selection == 1 {
				p.Correct = false
			} else if p.Selection == 1 {
				p.Correct = true
				total += p.Points
			} else {
				p.Correct = false
			}
			_, err := db.Update(&p)
			if err != nil {
				log.Fatalf("Update: %s", err.Error())
			}
		}

		fmt.Printf("%s: %d points\n", usersMap[u].FirstName, total)
	}

}
