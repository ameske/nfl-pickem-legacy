package database

import (
	"errors"
	"fmt"
	"log"
	"time"
)

// Picks is a struct mapping to the corresponding postgres table
type Picks struct {
	Id        int64
	UserId    int64
	GameId    int64
	Selection int
	Points    int
	Correct   bool
}

// UserPicksByWeek creates a []Picks representing a user's picks for the given week
func UserPicksByWeek(username string, year int, week int) ([]Picks, error) {
	sql := `SELECT picks.id, picks.game_id, picks.selection, picks.points, picks.correct
		FROM picks
		JOIN games ON games.id = picks.game_id
		JOIN weeks ON weeks.id = games.week_id
		JOIN years ON years.id = weeks.year_id
		JOIN users ON users.id = picks.user_id
		WHERE years.year = ?1 AND weeks.week = ?2 AND users.email = ?3
		ORDER BY games.date, games.id ASC`

	rows, err := db.Query(sql, year, week, username)
	if err != nil {
		log.Println("User picks by week")
		return nil, err
	}

	picks := make([]Picks, 0)
	for rows.Next() {
		tmp := Picks{}
		err := rows.Scan(&tmp.Id, &tmp.GameId, &tmp.Selection, &tmp.Points, &tmp.Correct)
		if err != nil {
			log.Fatal(err)
		}
		picks = append(picks, tmp)
	}
	rows.Close()

	return picks, nil
}

type SelectedPicks struct {
	Id        int64
	UserId    int64
	GameId    int64
	Selection int
	Points    int
	Correct   bool
	HomeNick  string
	AwayNick  string
}

// SelectedPicks returns the currently selected picks for a user for the given week
func UserSelectedPicksByWeek(username string, year int, week int) ([]SelectedPicks, error) {
	sql := `SELECT picks.id, picks.game_id, picks.selection, picks.points, picks.correct, hometeam.nickname, awayteam.nickname
		FROM picks
		JOIN games ON games.id = picks.game_id
		JOIN weeks ON weeks.id = games.week_id
		JOIN years ON years.id = weeks.year_id
		JOIN users ON users.id = picks.user_id
		JOIN teams AS hometeam ON games.home_id = hometeam.id
		JOIN teams AS awayteam ON games.away_id = awayteam.id
		WHERE years.year = ?1 AND weeks.week = ?2 AND users.email = ?3 AND picks.selection <> -1
		ORDER BY games.date ASC, games.id ASC`

	rows, err := db.Query(sql, year, week, username)
	if err != nil {
		return nil, err
	}

	picks := make([]SelectedPicks, 0)
	for rows.Next() {
		tmp := SelectedPicks{}
		err := rows.Scan(&tmp.Id, &tmp.GameId, &tmp.Selection, &tmp.Points, &tmp.Correct, &tmp.HomeNick, &tmp.AwayNick)
		if err != nil {
			return nil, err
		}
		picks = append(picks, tmp)
	}
	rows.Close()

	return picks, nil
}

//FormPick contains the information required to populate a picks HTML Form
type FormPick struct {
	Id               int64
	Time             time.Time
	Away             string
	AwayNick         string
	AwayRecord       string
	AwayAbbreviation string
	Home             string
	HomeNick         string
	HomeRecord       string
	HomeAbbreviation string
	Selection        int
	Points           int
	Disabled         bool
}

// FormPicks gathers the neccessary information needed render a user's pick-em form
func PicksFormByWeek(username string, year int, week int) ([]FormPick, error) {
	sql := `SELECT picks.id, picks.selection, picks.points, hometeam.City, hometeam.Nickname, hometeam.Abbreviation, awayteam.City, awayteam.Nickname, awayteam.Abbreviation, games.date
		FROM picks
		JOIN games ON games.id = picks.game_id
		JOIN teams AS hometeam ON games.home_id = hometeam.id
		JOIN teams AS awayteam ON games.away_id = awayteam.id
		JOIN weeks ON weeks.id = games.week_id
		JOIN years ON years.id = weeks.year_id
		JOIN users ON users.id = picks.user_id
		WHERE years.year = ?1 AND weeks.week = ?2 AND users.email = ?3
		ORDER BY games.date ASC, games.id ASC`

	rows, err := db.Query(sql, year, week, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	formPicks := make([]FormPick, 0)
	for rows.Next() {
		tmp := FormPick{}
		gametime := 0
		err := rows.Scan(&tmp.Id, &tmp.Selection, &tmp.Points, &tmp.Home, &tmp.HomeNick, &tmp.HomeAbbreviation, &tmp.Away, &tmp.AwayNick, &tmp.AwayAbbreviation, &gametime)
		if err != nil {
			return nil, err
		}

		tmp.Time = time.Unix(int64(gametime), 0)
		tmp.Disabled = time.Now().After(tmp.Time)

		hw, hl, err := teamRecord(tmp.Home, tmp.HomeNick)
		if err != nil {
			return nil, err
		}

		aw, al, err := teamRecord(tmp.Away, tmp.AwayNick)
		if err != nil {
			return nil, err
		}

		tmp.HomeRecord = fmt.Sprintf("(%d-%d)", hw, hl)
		tmp.AwayRecord = fmt.Sprintf("(%d-%d)", aw, al)

		formPicks = append(formPicks, tmp)
	}

	return formPicks, nil
}

var (
	ErrGameLocked = errors.New("game time locked out")
)

func MakePick(now time.Time, pickID int64, selection int, points int) error {
	sql := `SELECT games.date FROM games
		JOIN picks ON picks.game_id = games.id
		WHERE picks.id = ?1`

	var gameTime int64
	err := db.QueryRow(sql, pickID).Scan(&gameTime)
	if err != nil {
		return err
	}

	if gameTime < now.Unix() {
		return ErrGameLocked
	}

	_, err = db.Exec("UPDATE picks SET selection = ?1, points = ?2 WHERE id = ?3", selection, points, pickID)

	return err
}

func AdminMakePick(pickID int64, selection int, points int) error {
	_, err := db.Exec("UPDATE picks SET selection = ?1, points = ?2 WHERE id = ?3", selection, points, pickID)

	return err
}

func UpdatePick(id int64, correct bool, points int) error {
	var intBool int
	if correct {
		intBool = 1
	} else {
		intBool = 0
	}
	_, err := db.Exec("UPDATE picks SET correct = ?1, points = ?2 WHERE id = ?3", intBool, points, id)
	return err
}

type ResultPick struct {
	Id        int64
	HomeScore int
	AwayScore int
	Selection int
	Points    int
}

func UserResultPicksByWeek(user string, year int, week int) ([]ResultPick, error) {
	sql := `SELECT picks.id, games.home_score, games.away_score, picks.selection, picks.points
	       FROM picks
	       JOIN games ON games.id = picks.game_id
	       JOIN users ON users.id = picks.user_id
	       JOIN weeks ON weeks.id = games.week_id
	       JOIN years ON weeks.year_id = years.id
	       WHERE users.email = ?1 AND years.year = ?2 AND weeks.week = ?3`

	rows, err := db.Query(sql, user, year, week)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rp := make([]ResultPick, 0)
	for rows.Next() {
		tmp := ResultPick{}
		err := rows.Scan(&tmp.Id, &tmp.HomeScore, &tmp.AwayScore, &tmp.Selection, &tmp.Points)
		if err != nil {
			return nil, err
		}
		rp = append(rp, tmp)
	}

	return rp, nil
}

func CreateSeasonPicks(year int) error {
	sql := `SELECT id FROM users`

	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	users := make([]int, 0)
	for rows.Next() {
		var tmp int
		err = rows.Scan(&tmp)
		if err != nil {
			return err
		}
		users = append(users, tmp)
	}

	sql = `SELECT games.id FROM games
		JOIN weeks ON games.week_id = weeks.id
		JOIN years ON weeks.year_id = years.id
		WHERE years.year = ?1`

	rows, err = db.Query(sql, year)
	if err != nil {
		return err
	}
	defer rows.Close()

	games := make([]int, 0)
	for rows.Next() {
		var tmp int
		err = rows.Scan(&tmp)
		if err != nil {
			return err
		}

		games = append(games, tmp)
	}

	for _, g := range games {
		for _, u := range users {
			sql := "INSERT INTO picks(user_id, game_id) VALUES(?1, ?2)"
			_, err := db.Exec(sql, u, g)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
