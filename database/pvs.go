package database

import (
	"log"

	"github.com/coopernurse/gorp"
)

type Pvs struct {
	Id    int64  `db:"id"`
	Type  string `db:"type"`
	Seven int    `db:"seven"`
	Five  int    `db:"five"`
	Three int    `db:"three"`
	One   int    `db:"one"`
}

func WeekPvs(db *gorp.DbMap, year, week int) Pvs {
	var pvs Pvs
	err := db.SelectOne(&pvs, "SELECT pvs.* FROM weeks JOIN years ON years.id = week.year_id JOIN pvs ON weeks.pvs_id = pvs.id WHERE years.year = $1 AN weeks.week = $1", year, week)
	if err != nil {
		log.Fatalf("WeekPvs: %s", err.Error())
	}

	return pvs
}
