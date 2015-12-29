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

func TeamRecordMap() map[int64]*Record {
	homeWinsSQL := "SELECT home_id, COUNT(*) FROM games WHERE (home_score > away_score) GROUP BY home_id"
	awayWinsSQL := "SELECT away_id, COUNT(*) FROM games WHERE (away_score > home_score) GROUP BY away_id"
	homeLossesSQL := "SELECT home_id, COUNT(*) FROM games WHERE (home_score < away_score) GROUP BY home_id"
	awayLossesSQL := "SELECT away_id, COUNT(*) FROM games WHERE (away_score < home_score) GROUP BY away_id"

	records := make(map[int64]*Record)

	rows, err := db.Query(homeWinsSQL)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id int64
		var wins int
		err := rows.Scan(&id, &wins)
		if err != nil {
			log.Fatal(err)
		}
		records[id] = &Record{
			Wins: wins,
		}
	}
	rows.Close()

	rows, err = db.Query(awayWinsSQL)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id int64
		var wins int
		err := rows.Scan(&id, &wins)
		if err != nil {
			log.Fatal(err)
		}
		records[id].Wins += wins
	}
	rows.Close()

	rows, err = db.Query(homeLossesSQL)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id int64
		var losses int
		err := rows.Scan(&id, &losses)
		if err != nil {
			log.Fatal(err)
		}
		records[id].Losses += losses
	}
	rows.Close()

	rows, err = db.Query(awayLossesSQL)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id int64
		var losses int
		err := rows.Scan(&id, &losses)
		if err != nil {
			log.Fatal(err)
		}
		records[id].Losses += losses
	}
	rows.Close()

	return records
}
