package database

import (
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

//FormPick contains the information required to populate a picks HTML Form
type FormPick struct {
	Id        int64
	Time      time.Time
	Away      string
	AwayNick  string
	AwayId    int64
	Home      string
	HomeNick  string
	HomeId    int64
	Selection int
	Points    int
	Disabled  bool
}

// WeeklyPicks creates a []Picks representing a user's picks for the given week
func WeeklyPicks(username string) []*Picks {
	year, week := CurrentWeek()

	sql := `SELECT picks.id, picks.game_id, picks.selection, picks.points, picks.correct
		FROM picks
		JOIN games ON games.id = picks.game_id
		JOIN weeks ON weeks.id = games.week_id
		JOIN years ON years.id = weeks.year_id
		JOIN users ON users.id = picks.user_id
		WHERE years.year = ?1 AND weeks.week = ?2 AND users.email = ?3
		ORDER BY games.date ASC`

	rows, err := db.Query(sql, year, week, username)
	if err != nil {
		log.Fatal(err)
	}

	picks := make([]*Picks, 0)
	for rows.Next() {
		tmp := &Picks{}
		err := rows.Scan(&tmp.Id, &tmp.GameId, &tmp.Selection, &tmp.Points, &tmp.Correct)
		if err != nil {
			log.Fatal(err)
		}
		picks = append(picks, tmp)
	}
	rows.Close()

	return picks
}

func WeeklyPicksYearWeek(username string, year, week int) []*Picks {
	sql := `SELECT picks.id, picks.game_id, picks.selection, picks.points, picks.correct
		FROM picks
		JOIN games ON games.id = picks.game_id
		JOIN weeks ON weeks.id = games.week_id
		JOIN years ON years.id = weeks.year_id
		JOIN users ON users.id = picks.user_id
		WHERE years.year = ?1 AND weeks.week = ?2 AND users.email = ?3
		ORDER BY games.date ASC`

	rows, err := db.Query(sql, year, week, username)
	if err != nil {
		log.Fatal(err)
	}

	picks := make([]*Picks, 0)
	for rows.Next() {
		tmp := &Picks{}
		err := rows.Scan(&tmp.Id, &tmp.GameId, &tmp.Selection, &tmp.Points, &tmp.Correct)
		if err != nil {
			log.Fatal(err)
		}
		picks = append(picks, tmp)
	}
	rows.Close()

	return picks
}

func weeklySelectedPicks(username string) []*Picks {
	year, week := CurrentWeek()

	sql := `SELECT picks.id, picks.game_id, picks.selection, picks.points, picks.correct
		FROM picks
		JOIN games ON games.id = picks.game_id
		JOIN weeks ON weeks.id = games.week_id
		JOIN years ON years.id = weeks.year_id
		JOIN users ON users.id = picks.user_id
		WHERE years.year = ?1 AND weeks.week = ?2 AND users.email = ?3 AND picks.selection <> 0
		ORDER BY games.date ASC`

	rows, err := db.Query(sql, year, week, username)
	if err != nil {
		log.Fatal(err)
	}

	picks := make([]*Picks, 0)
	for rows.Next() {
		tmp := &Picks{}
		err := rows.Scan(&tmp.Id, &tmp.GameId, &tmp.Selection, &tmp.Points, &tmp.Correct)
		if err != nil {
			log.Fatal(err)
		}
		picks = append(picks, tmp)
	}
	rows.Close()

	return picks
}

// FormPicks gathers the neccessary information needed render a user's pick-em form
func FormPicks(username string, selectedOnly bool) []FormPick {
	var picks []*Picks
	if selectedOnly {
		picks = weeklySelectedPicks(username)
	} else {
		picks = WeeklyPicks(username)
	}

	formGames := make([]FormPick, 0)
	for _, p := range picks {
		// Lookup the game information
		awayId, homeId, date := gamePickInfo(p.GameId)

		// Lookup the team information
		h := teamById(homeId)
		a := teamById(awayId)

		// Convert to a go time for the form picks
		gDate := time.Unix(date, 0)

		// 8/23/15: Since we now append time zone when loading schedule, we should be ok
		//disabled := time.Now().After(gDate.Add(time.Hour * 5))
		disabled := time.Now().After(gDate)

		// Construct the FormPick
		f := FormPick{
			Id:        p.Id,
			Time:      gDate,
			Away:      a.City,
			AwayNick:  a.Nickname,
			AwayId:    a.Id,
			Home:      h.City,
			HomeNick:  h.Nickname,
			HomeId:    h.Id,
			Selection: p.Selection,
			Points:    p.Points,
			Disabled:  disabled,
		}
		formGames = append(formGames, f)
	}

	return formGames
}

func MakePick(pickID int64, selection int, points int) error {
	_, err := db.Exec("UPDATE picks SET selection = ?1, points = ?2 WHERE id = ?3", selection, points, pickID)

	return err
}

func UpdatePick(id int64, correct bool) error {
	_, err := db.Exec("UPDATE picks SET correct = ?1 WHERE id = ?2", correct, id)
	return err
}

type PickSelection struct {
	Picks
	AwayId int64
	HomeId int64
}

func CLIPickSelections(user string, year, week int) []PickSelection {
	sql := `SELECT picks.id, picks.user_id, picks.game_id, picks.selection, picks.points, picks.correct, games.away_id, games.home_id
		FROM picks
		JOIN users ON picks.user_id = users.id
		JOIN games ON picks.game_id = games.id
		JOIN weeks ON games.week_id = weeks.id
		JOIN years ON weeks.year_id = years.id
		WHERE users.email = ?1 AND years.year = ?2 AND weeks.week = ?3`

	rows, err := db.Query(sql, user, year, week)
	if err != nil {
		log.Fatal(err)
	}

	picks := make([]PickSelection, 0)
	for rows.Next() {
		tmp := PickSelection{}
		err := rows.Scan(&tmp.Id, &tmp.UserId, &tmp.GameId, &tmp.Selection, &tmp.Points, &tmp.Correct, &tmp.AwayId, &tmp.HomeId)
		if err != nil {
			log.Fatal(err)
		}
		picks = append(picks, tmp)
	}

	rows.Close()

	return picks
}

func CreateSeasonPicks(year int) {
	users := AllUsers()
	games := GamesBySeason(year)
	for _, g := range games {
		for _, u := range users {
			sql := "INSERT INTO picks(user_id, game_id) VALUES(?1, ?2)"
			_, err := db.Exec(sql, u.Id, g.Id)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
