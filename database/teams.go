package database

import (
	"fmt"
	"log"

	"github.com/coopernurse/gorp"
)

type Teams struct {
	Id           int64  `db:"id"`
	City         string `db:"city"`
	Nickname     string `db:"nickname"`
	Stadium      string `db:"stadium"`
	Abbreviation string `db:"abbreviation"`
}

func TeamAbbreviationMap(db *gorp.DbMap) map[int64]string {
	var teams []Teams

	teamMap := make(map[int64]string)

	_, err := db.Select(&teams, "SELECT * FROM teams")
	if err != nil {
		log.Fatalf("TeamMap: %s", err.Error())
	}

	for _, t := range teams {
		teamMap[t.Id] = t.Abbreviation
	}

	return teamMap

}

func TeamMap(db *gorp.DbMap) map[int64]string {
	var teams []Teams

	teamMap := make(map[int64]string)

	_, err := db.Select(&teams, "SELECT * FROM teams")
	if err != nil {
		log.Fatalf("TeamMap: %s", err.Error())
	}

	for _, t := range teams {
		teamMap[t.Id] = fmt.Sprintf("%s %s", t.City, t.Nickname)
	}

	return teamMap
}
