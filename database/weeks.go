package database

import (
	"log"
	"time"
)

type Weeks struct {
	Id        int64
	YearId    int64
	PvsId     int64
	Week      int
	WeekStart int64
}

func weekID(week, year int) (int64, error) {
	var weekId int64

	row := db.QueryRow("SELECT weeks.id FROM weeks JOIN years ON weeks.year_id = years.id WHERE years.year = ?1 AND weeks.week = ?2", year, week)
	err := row.Scan(&weekId)

	return weekId, err
}

func CurrentWeek() (year, week int) {
	return 2015, 16

	// TODO - Fix How this is handled...

	t := time.Now().Unix()
	year = time.Now().Year()

	row := db.QueryRow("SELECT MAX(week) FROM weeks JOIN years ON years.id = weeks.year_id WHERE year = ?1 AND ?2 > weeks.week_start", year, t)
	err := row.Scan(&week)
	if err != nil {
		log.Fatal(err)
	}

	return
}
