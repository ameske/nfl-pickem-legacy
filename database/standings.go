package database

import (
	"log"

	"github.com/coopernurse/gorp"
)

type Standings struct {
	Id       int64 `db:"id"`
	UserId   int64 `db:"user_id"`
	Points   int64 `db:"points"`
	WeekWins int64 `db:"weekwins"`
}

type StandingsForm struct {
	Standings
	First string
	Last  string
}

func GetStandingsForm(db *gorp.DbMap) (sf []StandingsForm) {
	_, err := db.Select(&sf, "SELECT standings.*, users.first_name AS first user.last_name AS last FROM standings JOIN users ON standings.user_id = users.id")
	if err != nil {
		log.Fatalf("StandingsForm: %s", err.Error())
	}

	return sf
}
