package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Pick struct {
	ID        int64
	UserID    int64
	GameID    int64
	Selection int64
	Points    int64
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	db, err := sql.Open("sqlite3", "nfl.db")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("CREATE TABLE newPicks(id integer PRIMARY KEY, user_id integer REFERENCES user(id), game_id integer REFERENCES games(id), selection_id integer REFERENCES teams(id), points integer DEFAULT 0)")
	if err != nil {
		log.Fatal(err)
	}

	// for each pick, recreate it in newpicks
	rows, err := db.Query("SELECT picks.id, picks.user_id, picks.game_id, picks.selection, picks.points FROM picks")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	picks := make([]Pick, 0)

	for rows.Next() {
		var p Pick

		err := rows.Scan(&p.ID, &p.UserID, &p.GameID, &p.Selection, &p.Points)
		if err != nil {
			log.Fatal(err)
		}

		picks = append(picks, p)
	}

	for _, p := range picks {
		var updateSQL string

		if p.Selection == 1 {
			updateSQL = "INSERT INTO newPicks(user_id, game_id, selection_id, points) VALUES (?1, ?2, (SELECT games.away_id FROM picks JOIN games ON picks.game_id = games.id WHERE picks.id = ?3), ?4)"
			_, err = db.Exec(updateSQL, p.UserID, p.GameID, p.ID, p.Points)
		} else if p.Selection == 2 {
			updateSQL = "INSERT INTO newPicks(user_id, game_id, selection_id, points) VALUES (?1, ?2, (SELECT games.home_id FROM picks JOIN games ON picks.game_id = games.id WHERE picks.id = ?3), ?4)"
			_, err = db.Exec(updateSQL, p.UserID, p.GameID, p.ID, p.Points)
		} else {
			updateSQL = "INSERT INTO newPicks(user_id, game_id, selection_id, points) VALUES (?1, ?2, NULL, ?3)"
			_, err = db.Exec(updateSQL, p.UserID, p.GameID, p.Points)
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	_, err = db.Exec("DROP TABLE picks")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("ALTER TABLE newPicks RENAME TO picks")
	if err != nil {
		log.Fatal(err)
	}
}
