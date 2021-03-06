package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/ameske/nfl-pickem/database"
	"github.com/ameske/nfl-pickem/results"
)

func logNextScheduleUpdate(t time.Time) {
	slog.Info("Scheduling update for " + t.Format(time.RFC1123))
}

// scheduleUpdates sets up goroutines that will import the results of games and update the
// picks after every wave of games completes.
func scheduleUpdates() {
	// Friday at 8:00
	go func() {
		nextFriday := adjustIfPast(nextDay(time.Friday).Add(time.Hour * 8))
		logNextScheduleUpdate(nextFriday)
		time.Sleep(nextFriday.Sub(time.Now()))
		for {
			go update(false)
			logNextScheduleUpdate(time.Now().AddDate(0, 0, 7))
			time.Sleep(time.Hour * 24 * 7)
		}
	}()

	// Sunday at 18:00
	go func() {
		nextSunday := adjustIfPast(nextDay(time.Sunday).Add(time.Hour * 18))
		logNextScheduleUpdate(nextSunday)
		time.Sleep(nextSunday.Sub(time.Now()))
		for {
			go update(false)
			logNextScheduleUpdate(time.Now().AddDate(0, 0, 7))
			time.Sleep(time.Hour * 24 * 7)
		}
	}()

	// Sunday at 21:00
	go func() {
		nextSunday := adjustIfPast(nextDay(time.Sunday).Add(time.Hour * 21))
		logNextScheduleUpdate(nextSunday)
		time.Sleep(nextSunday.Sub(time.Now()))
		for {
			go update(false)
			logNextScheduleUpdate(time.Now().AddDate(0, 0, 7))
			time.Sleep(time.Hour * 24 * 7)
		}
	}()

	// Monday at 8:00
	go func() {
		nextMonday := adjustIfPast(nextDay(time.Monday).Add(time.Hour * 8))
		logNextScheduleUpdate(nextMonday)
		time.Sleep(nextMonday.Sub(time.Now()))
		for {
			go update(false)
			logNextScheduleUpdate(time.Now().AddDate(0, 0, 7))
			time.Sleep(time.Hour * 24 * 7)
		}
	}()

	// Tuesday at 8:00. Here we need to update the current week - 1
	go func() {
		nextTuesday := adjustIfPast(nextDay(time.Tuesday).Add(time.Hour * 8))
		logNextScheduleUpdate(nextTuesday)
		time.Sleep(nextTuesday.Sub(time.Now()))
		for {
			go update(true)
			logNextScheduleUpdate(time.Now().AddDate(0, 0, 7))
			time.Sleep(time.Hour * 24 * 7)
		}
	}()
}

func update(updatePreviousWeek bool) {
	year, week, err := database.CurrentWeek(time.Now())
	if err == database.ErrOffseason {
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	if updatePreviousWeek {
		week -= 1
	}

	results, err := getGameResults(year, week)
	if err != nil {
		log.Println(err)
		return
	}

	err = updateGameScores(year, week, results)
	if err != nil {
		log.Println(err)
		return
	}

	err = grade(year, week)
	if err != nil {
		log.Println(err)
		return
	}
}

type ResultsJson struct {
	Week      int    `json:"week"`
	Year      int    `json:"year"`
	Home      string `json:"home"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
}

func getGameResults(year, week int) ([]results.Result, error) {
	url := fmt.Sprintf("http://www.nfl.com/schedules/%d/REG%d", year, week)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	p := results.NewParser(resp.Body)

	return p.Parse()
}

func updateGameScores(year int, week int, results []results.Result) error {
	for _, result := range results {
		err := database.UpdateScore(week, year, result.Home, result.HomeScore, result.AwayScore)
		if err != nil {
			return err
		}
	}

	return nil
}

// Grade calculates the scores for each user in the database for the given week.
// It assumes that the scores for the graded week have already been imported, else
// results are undefined.
func grade(year, week int) error {
	users, err := database.Usernames()
	if err != nil {
		return err
	}

	// For each user, score their picks for this week
	for _, u := range users {
		picks, err := database.UserResultPicksByWeek(u, year, week)
		if err != nil {
			return err
		}

		for _, p := range picks {
			var correct bool
			var points int

			// Ignore all games that haven't finished yet - clean up points though
			if p.HomeScore == -1 && p.AwayScore == -1 {
				err := database.UpdatePick(p.Id, false, p.Points)
				if err != nil {
					return err
				}
				continue
			}

			if p.HomeScore == p.AwayScore {
				correct = true
				points = int(math.Floor(float64(p.Points) / 2))
			} else if p.HomeScore > p.AwayScore && p.Selection == 2 {
				correct = true
				points = p.Points
			} else if p.HomeScore > p.AwayScore && p.Selection == 1 {
				correct = false
				points = p.Points
			} else if p.AwayScore > p.HomeScore && p.Selection == 2 {
				correct = false
				points = p.Points
			} else {
				correct = true
				points = p.Points
			}

			err := database.UpdatePick(p.Id, correct, points)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func adjustIfPast(next time.Time) time.Time {
	now := time.Now()
	ty, tm, td := now.Date()
	ny, nm, nd := next.Date()

	// if the next day is today, but the hour we want has past then we must advance next a week
	if (ty == ny && tm == nm && td == nd) && now.Hour() > next.Hour() {
		next = next.AddDate(0, 0, 7)
	}

	return next
}

func nextDay(day time.Weekday) time.Time {
	now := time.Now()

	// We only want to go forwards, so use modular arith to force going ahead
	diff := int(day-now.Weekday()+7) % 7

	next := now.AddDate(0, 0, diff)
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())

	return next
}
