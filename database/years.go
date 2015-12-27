package database

import "log"

type Years struct {
	Id   int64
	Year int
}

func yearId(year int) (id int64) {
	row := db.QueryRow("SELECT id FROM years WHERE year = ?1", year)
	err := row.Scan(&id)
	if err != nil {
		log.Fatal(err)
	}

	return id
}
