package database

type PVS struct {
	Seven int
	Five  int
	Three int
	One   int
}

func WeekPVS(year int, week int) (PVS, error) {
	var pvs PVS
	sql := `SELECT pvs.Seven, pvs.Five, pvs.Three, pvs.One
		FROM weeks
		JOIN years ON years.id = weeks.year_id
		JOIN pvs ON pvs.id = weeks.pvs_id
		WHERE years.year = ?1 AND weeks.week = ?2`

	row := db.QueryRow(sql, year, week)
	err := row.Scan(&pvs.Seven, &pvs.Five, &pvs.Three, &pvs.One)

	return pvs, err
}
