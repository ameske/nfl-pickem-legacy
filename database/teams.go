package database

type Teams struct {
	Id       int64  `db:"id"`
	City     string `db:"city"`
	Nickname string `db:"nickname"`
	Stadium  string `db:"stadium"`
}
