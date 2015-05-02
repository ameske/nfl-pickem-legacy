package database

type Years struct {
	Id   int64 `db:"id"`
	Year int   `db:"year"`
}

func NewYear(year int) error {
	y := Years{
		Year: year,
	}

	return db.Insert(&y)
}

func yearID(year int) (int64, error) {
	return db.SelectInt("SELECT id FROM years WHERE year = $1", year)
}
