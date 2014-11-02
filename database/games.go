package database

import "time"

type Games struct {
	Id        int64     `db:"id"`
	WeekId    int64     `db:"week_id"`
	Date      time.Time `db:"date"`
	HomeId    int64     `db:"home_id"`
	AwayId    int64     `db:"away_id"`
	HomeScore int       `db:"home_score"`
	AwayScore int       `db:"away_score"`
}
