package database

import (
	"log"

	"github.com/coopernurse/gorp"
)

type Years struct {
	Id   int64 `db:"id"`
	Year int   `db:"year"`
}

type Weeks struct {
	Id     int64 `db:"id"`
	YearId int64 `db:"year_id"`
	PvsId  int64 `db:"pvs_id"`
	Week   int   `db:"week"`
}

func WeekId(db *gorp.DbMap, year, week int) int64 {
	var weekId int64
	err := db.SelectOne(&weekId, "SELECT weeks.id FROM weeks JOIN years ON weeks.year_id = years.id WHERE years.year = $1 AND weeks.week = $2", year, week)
	if err != nil {
		log.Fatalf("WeekId: %s", err.Error())
	}

	return weekId
}
