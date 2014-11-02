package database

type Pvs struct {
	Id    int64  `db:"id"`
	Type  string `db:"type"`
	Seven int    `db:"seven"`
	Five  int    `db:"five"`
	Three int    `db:"three"`
	One   int    `db:"one"`
}
