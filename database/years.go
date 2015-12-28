package database

import (
	"log"
	"time"
)

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

func YearStart(year int) time.Time {
	var start int64
	row := db.QueryRow("SELECT year_start FROM years WHERE year = ?1", year)
	err := row.Scan(&start)
	if err != nil {
		log.Fatal(err)
	}

	return time.Unix(start, 0)
}
