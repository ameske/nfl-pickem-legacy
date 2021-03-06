package database

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrOffseason = errors.New("offsesaon")
)

const (
	oneWeek      = time.Hour * 24 * 7
	seasonLength = 17
)

func PrevSeasonExists(year int) bool {
	var c int
	err := db.QueryRow("SELECT COUNT(*) FROM years WHERE year = ?1", year-1).Scan(&c)
	if err != nil {
		return false
	}

	return c != 0
}

func currentSeasonStart(t time.Time) (start time.Time, err error) {
	now := t.Unix()

	var s sql.NullInt64
	row := db.QueryRow("SELECT MAX(year_start) FROM years WHERE year_start < ?1", now)
	err = row.Scan(&s)
	if err != nil {
		return time.Unix(0, 0), err
	}

	// Special case: if now + 7 is a different value then that means we're on the cusp of a new season. So pretend we are in week 1.
	now2 := time.Date(t.Year(), t.Month(), t.Day()+7, t.Hour(), t.Minute(), t.Second(), 0, t.Location())
	var s2 sql.NullInt64
	row = db.QueryRow("SELECT MAX(year_start) FROM years WHERE year_start < ?1", now2.Unix())
	err = row.Scan(&s2)
	if err != nil {
		return time.Unix(0, 0), err
	}

	if s.Int64 != s2.Int64 && s2.Int64 != 0 {
		return time.Unix(s2.Int64-604800, 0), err
	} else if s.Valid {
		return time.Unix(s.Int64, 0), err
	} else {
		return time.Unix(0, 0), err
	}
}

func CurrentWeek(t time.Time) (year int, week int, err error) {
	/*
		The season starts on the Tuesday before the first game.
		To figure out what week we are in, calculate where we are from there.
	*/
	start, err := currentSeasonStart(t)
	if err != nil {
		return -1, -1, err
	}

	d := t.Sub(start)

	week = int(d/oneWeek) + 1

	if week > seasonLength {
		return start.Year(), -1, ErrOffseason
	}

	return start.Year(), week, nil
}

func IsOffseason(t time.Time) bool {
	_, _, err := CurrentWeek(t)

	return err == ErrOffseason
}

func AddWeek(year int, week int, numGames int) error {
	pvs := ""
	switch numGames {
	case 16:
		pvs = "A"
	case 15:
		pvs = "B"
	case 14:
		pvs = "C"
	case 13:
		pvs = "D"
	}

	_, err := db.Exec("INSERT INTO weeks(week, year_id, pvs_id) VALUES(?1, (SELECT id FROM YEARS where year = ?2), (SELECT id FROM pvs WHERE type = ?3))", week, year, pvs)

	return err
}

func AddYear(year int, yearStart int) error {
	_, err := db.Exec("INSERT INTO years(year, year_start) VALUES(?1, ?2)", year, yearStart)

	return err
}
