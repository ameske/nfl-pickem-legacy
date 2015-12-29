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

type Record struct {
	Wins   int
	Losses int
}

func TeamRecordMap() map[int64]Record {
	sql := `SELECT tmp.id, home_wins+away_wins AS wins, home_losses+away_losses AS losses
	  FROM (SELECT home_id AS id, COUNT(*) AS home_wins FROM games WHERE (home_score > away_score) GROUP BY home_id) tmp
	  JOIN (SELECT away_id AS id, COUNT(*) AS away_wins FROM games WHERE (away_score > home_score) GROUP BY away_id) tmp2 ON tmp.id = tmp2.id
	  JOIN (SELECT home_id AS id, COUNT(*) AS home_losses FROM games WHERE (home_score < away_score) GROUP BY home_id) tmp3 ON tmp2.id = tmp3.id
	  JOIN (SELECT away_id AS id, COUNT(*) AS away_losses FROM games WHERE (away_score < home_score) GROUP BY away_id) tmp4 ON tmp3.id = tmp4.id;`

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}

	records := make(map[int64]Record)

	for rows.Next() {
		var id int64
		var wins, losses int
		err := rows.Scan(&id, &wins, &losses)
		if err != nil {
			log.Fatal(err)
		}

		records[id] = Record{Wins: wins, Losses: losses}
	}

	rows.Close()

	return records
}
