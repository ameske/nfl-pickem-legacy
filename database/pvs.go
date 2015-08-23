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

func NewPvs(one, three, five, seven int, typeID string) error {
	pvs := Pvs{
		One:   one,
		Three: three,
		Five:  five,
		Seven: seven,
		Type:  typeID,
	}

	return db.Insert(&pvs)
}

func WeekPvs() Pvs {
	var pvs Pvs
	sql := `SELECT pvs.*
		FROM weeks
		JOIN years ON years.id = weeks.year_id
		JOIN pvs ON pvs.id = weeks.pvs_id
		WHERE years.year = ?1 AND weeks.week = ?2`

	year, week := CurrentWeek()
	err := db.SelectOne(&pvs, sql, year, week)
	if err != nil {
		log.Fatalf("WeekPvs: %s", err.Error())
	}

	return pvs
}

func pvsID(typeID string) (int64, error) {
	return db.SelectInt("SELECT id FROM pvs WHERE type = ?1", typeID)
}
