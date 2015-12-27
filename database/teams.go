package database

import (
	"fmt"
	"log"
)

type Teams struct {
	Id           int64
	City         string
	Nickname     string
	Stadium      string
	Abbreviation string
}

func teamById(id int64) Teams {
	var t Teams
	row := db.QueryRow("SELECT id, city, nickname, stadium, abbreviation FROM teams WHERE id = ?1", id)
	err := row.Scan(&t.Id, &t.City, &t.Nickname, &t.Stadium, &t.Abbreviation)
	if err != nil {
		log.Fatal(err)
	}

	return t
}

func TeamAbbreviationMap() map[int64]string {
	rows, err := db.Query("SELECT id, abbreviation FROM teams")
	if err != nil {
		log.Fatal(err)
	}

	teams := make([]Teams, 0)
	for rows.Next() {
		tmp := Teams{}
		err := rows.Scan(&tmp.Id, &tmp.Abbreviation)
		if err != nil {
			log.Fatal(err)
		}
		teams = append(teams, tmp)
	}
	rows.Close()

	teamMap := make(map[int64]string)
	for _, t := range teams {
		teamMap[t.Id] = t.Abbreviation
	}

	return teamMap
}

func TeamMap() map[int64]string {
	rows, err := db.Query("SELECT id, city, nickname FROM teams")
	if err != nil {
		log.Fatal(err)
	}

	teams := make([]Teams, 0)
	for rows.Next() {
		tmp := Teams{}
		err := rows.Scan(&tmp.Id, &tmp.City, &tmp.Nickname)
		if err != nil {
			log.Fatal(err)
		}
		teams = append(teams, tmp)
	}
	rows.Close()

	teamMap := make(map[int64]string)
	for _, t := range teams {
		teamMap[t.Id] = fmt.Sprintf("%s %s", t.City, t.Nickname)
	}

	return teamMap
}
