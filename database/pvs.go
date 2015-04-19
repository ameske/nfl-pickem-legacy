package database

import "log"

type Pvs struct {
	Id    int64  `db:"id"`
	Type  string `db:"type"`
	Seven int    `db:"seven"`
	Five  int    `db:"five"`
	Three int    `db:"three"`
	One   int    `db:"one"`
}

func WeekPvs() Pvs {
	var pvs Pvs
	sql := `SELECT pvs.*
		FROM weeks
		JOIN years ON years.id = weeks.year_id
		JOIN pvs ON pvs.id = weeks.pvs_id
		WHERE years.year = $1 AND weeks.week = $2`

	year, week := CurrentWeek()
	err := db.SelectOne(&pvs, sql, year, week)
	if err != nil {
		log.Fatalf("WeekPvs: %s", err.Error())
	}

	return pvs
}
