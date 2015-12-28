package database

import "time"

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

var oneWeek = time.Hour * 24 * 7

func CurrentWeek() (year, week int) {
	year = time.Now().Year()

	/*
		The season starts on the Tuesday before the first game.
		To figure out what week we are in, calculate where we are from there.
	*/
	seasonStart := YearStart(year)
	today := time.Now()
	d := today.Sub(seasonStart)

	week = int(d/oneWeek) + 1

	return
}
