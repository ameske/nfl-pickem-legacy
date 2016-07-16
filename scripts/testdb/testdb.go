package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/ameske/nfl-pickem/database"
)

var (
	teams = map[int]string{
		1:  "Ravens",
		2:  "Bengals",
		3:  "Browns",
		4:  "Steelers",
		5:  "Bears",
		6:  "Lions",
		7:  "Packers",
		8:  "Vikings",
		9:  "Texans",
		10: "Colts",
		11: "Jaguars",
		12: "Titans",
		13: "Falcons",
		14: "Panthers",
		15: "Saints",
		16: "Buccaneers",
		17: "Bills",
		18: "Dolphins",
		19: "Patriots",
		20: "Jets",
		21: "Cowboys",
		22: "Giants",
		23: "Eagles",
		24: "Redskins",
		25: "Broncos",
		26: "Chiefs",
		27: "Raiders",
		28: "Chargers",
		29: "Cardinals",
		30: "Rams",
		31: "49ers",
		32: "Seahawks",
	}
)

func main() {
	os.Remove("nfl-test.db")

	fd, err := os.Open("ddl.sql")
	if err != nil {
		log.Fatal(err)
	}

	var sout bytes.Buffer
	var serr bytes.Buffer

	init := exec.Command("sqlite3", "nfl-test.db")
	init.Stdin = fd
	init.Stdout = &sout
	init.Stderr = &serr

	err = init.Run()
	if err != nil {
		fmt.Println(sout.String())
		fmt.Println(serr.String())
		log.Fatal(err)
	}

	err = database.SetDefaultDb("nfl-test.db")
	if err != nil {
		log.Fatal(err)
	}

	err = addTestUsers()
	if err != nil {
		log.Fatal(err)
	}

	start, err := addWeek(time.Now(), 16)
	if err != nil {
		log.Fatal(err)
	}

	err = addGames(start, 16)
	if err != nil {
		log.Fatal(err)
	}

	err = database.CreateSeasonPicks(time.Now().Year())
	if err != nil {
		log.Fatal(err)
	}
}

func addWeek(t time.Time, numGames int) (time.Time, error) {
	next := testWeekDate(t)

	err := database.AddYear(t.Year(), int(next.Unix()))
	if err != nil {
		return next, err
	}

	return next, database.AddWeek(t.Year(), 1, 16)
}

func testWeekDate(t time.Time) time.Time {
	var next time.Time
	switch t.Weekday() {
	case time.Sunday, time.Monday:
		next = nextDay(t, time.Tuesday)
	case time.Tuesday, time.Wednesday:
		next = nextDay(time.Date(t.Year(), t.Month(), t.Day()-7, t.Hour(), t.Minute(), t.Second(), 0, t.Location()), time.Tuesday)
	default:
		next = nextDay(t, time.Tuesday)
	}

	return next
}

func addGames(start time.Time, numGames int) error {
	curTeam := 1

	// One game on Thursday
	thur := nextDay(start, time.Thursday)
	thur = time.Date(thur.Year(), thur.Month(), thur.Day(), 20, 30, 0, 0, thur.Location())
	err := database.AddGame(thur, teams[curTeam], teams[curTeam+1])
	if err != nil {
		return err
	}

	curTeam += 2

	// Nine games at 1:00 Sunday
	sunday := nextDay(start, time.Sunday)
	sunday = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 13, 0, 0, 0, sunday.Location())
	for i := 0; i < 9; i++ {
		err = database.AddGame(sunday, teams[curTeam], teams[curTeam+1])
		if err != nil {
			return err
		}

		curTeam += 2
	}

	// Three games at 4:00 Sunday
	sunday = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 16, 0, 0, 0, sunday.Location())
	for i := 0; i < 3; i++ {
		err = database.AddGame(sunday, teams[curTeam], teams[curTeam+1])
		if err != nil {
			return err
		}

		curTeam += 2
	}

	// One game at 4:25 Sunday
	sunday = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 16, 25, 0, 0, sunday.Location())
	err = database.AddGame(sunday, teams[curTeam], teams[curTeam+1])
	if err != nil {
		return err
	}

	curTeam += 2

	// One game on Sunday Night
	sunday = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 20, 30, 0, 0, sunday.Location())
	err = database.AddGame(sunday, teams[curTeam], teams[curTeam+1])
	if err != nil {
		return err
	}

	curTeam += 2

	// One game on Monday Night
	monday := nextDay(start, time.Monday)
	monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 20, 30, 0, 0, monday.Location())
	err = database.AddGame(monday, teams[curTeam], teams[curTeam+1])
	if err != nil {
		return err
	}

	return nil
}

func addTestUsers() error {
	alice := database.Users{
		FirstName: "Alice",
		LastName:  "Tester",
		Email:     "alice@gmail.com",
		Admin:     true,
		Password:  "password",
	}

	err := database.AddUser(alice)
	if err != nil {
		return err
	}

	bob := database.Users{
		FirstName: "Bob",
		LastName:  "Tester",
		Email:     "bob@gmail.com",
		Admin:     false,
		Password:  "password",
	}

	err = database.AddUser(bob)
	if err != nil {
		return err
	}

	return nil
}

func nextDay(now time.Time, day time.Weekday) time.Time {
	// We only want to go forwards, so use modular arith to force going ahead
	diff := int(day-now.Weekday()+7) % 7

	next := now.AddDate(0, 0, diff)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())

	return next
}
