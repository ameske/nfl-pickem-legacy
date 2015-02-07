package main

import (
	"log"

	"github.com/ameske/go_nfl/database"
	"github.com/codegangsta/cli"
)

type GamesWeeksJoin struct {
	database.Games
	database.Weeks
}

func seasonPicks(c *cli.Context) {
	year := c.Int("year")
	if year == -1 {
		log.Fatalf("Year is a required argument")
	}

	var users []database.Users
	_, err := db.Select(&users, "SELECT * FROM users")
	if err != nil {
		log.Fatalf("Users error: %s", err.Error())
	}

	var yearId int64
	err = db.SelectOne(&yearId, "SELECT id FROM years WHERE year = $1", year)
	if err != nil {
		log.Fatalf("Years error: %s", err.Error())
	}

	var games []int64
	_, err = db.Select(&games, "SELECT games.id FROM weeks INNER JOIN games ON games.week_id = weeks.id WHERE weeks.year_id = $1", yearId)
	if err != nil {
		log.Fatalf("Games error: %s", err.Error())
	}

	for _, g := range games {
		for _, u := range users {
			tmp := &database.Picks{
				UserId: u.Id,
				GameId: g,
			}
			err = db.Insert(tmp)
			if err != nil {
				log.Fatalf("Insert error: %s", err.Error())
			}

		}
	}

}
