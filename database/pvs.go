package database

import "log"

type Pvs struct {
	Id    int64
	Type  string
	Seven int
	Five  int
	Three int
	One   int
}

func WeekPvs(year int, week int) Pvs {
	var pvs Pvs
	sql := `SELECT pvs.Type, pvs.Seven, pvs.Five, pvs.Three, pvs.One
		FROM weeks
		JOIN years ON years.id = weeks.year_id
		JOIN pvs ON pvs.id = weeks.pvs_id
		WHERE years.year = ?1 AND weeks.week = ?2`

	row := db.QueryRow(sql, year, week)
	err := row.Scan(&pvs.Type, &pvs.Seven, &pvs.Five, &pvs.Three, &pvs.One)

	if err != nil {
		log.Fatal(err)
	}
	return pvs
}

func pvsID(typeID string) (int64, error) {
	var id int64

	row := db.QueryRow("SELECT id FROM pvs WHERE type = ?1", typeID)
	err := row.Scan(&id)

	return id, err
}
