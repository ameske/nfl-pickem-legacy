package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ameske/nfl-pickem/database"
	"github.com/gorilla/context"
)

type Pick struct {
	Id        int64
	Selection int
	Points    int
}

func Picks(templatesDir string, notifier Notifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			processPicks(templatesDir, notifier, w, r)
			return
		}

		picksForm(templatesDir, w, r)
	}
}

// Picks fetches this week's picks for the current logged in user and renders the
// picks template with them.
func picksForm(templatesDir string, w http.ResponseWriter, r *http.Request) {
	user, isAdmin := currentUser(r)
	if user == "" {
		http.Error(w, "no user information found", http.StatusUnauthorized)
		return
	}

	year, week, err := database.CurrentWeek(time.Now())
	if err == database.ErrOffseason {
		err := Fetch(templatesDir, "picks_offseason.html").Execute(w, user, isAdmin, nil)
		if err != nil {
			log.Println(err)
		}
		return
	} else if err != nil {
		log.Fatal(err)
	}

	picks, err := database.PicksFormByWeek(user, year, week)
	if err != nil {
		log.Fatal(err)
	}

	pvs, err := database.WeekPVS(year, week)
	if err != nil {
		log.Fatal(err)
	}

	var e, s string
	if context.Get(r, "error") != nil {
		e = context.Get(r, "error").(string)
	}
	if context.Get(r, "success") != nil {
		s = context.Get(r, "success").(string)
	}

	data := struct {
		Error   string
		Success string
		URL     string
		Picks   []database.FormPick
		database.PVS
	}{
		e,
		s,
		r.URL.String(),
		picks,
		pvs,
	}

	err = Fetch(templatesDir, "picks.html").Execute(w, user, isAdmin, data)
	if err != nil {
		log.Println(err)
	}
}

// ProcessPicks validates a user's picks, and then updates the current picks in the database
func processPicks(templatesDir string, notifier Notifier, w http.ResponseWriter, r *http.Request) {
	user, _ := currentUser(r)
	if user == "" {
		http.Error(w, "no user information found", http.StatusUnauthorized)
		return
	}

	year, week, err := database.CurrentWeek(time.Now())
	if err != nil {
		log.Fatal(err)
	}

	newPicks, err := parsePicksForm(r)
	if err != nil {
		log.Fatal(err)
	}

	existingPicks, err := database.UserSelectedPicksByWeek(user, year, week)
	if err != nil {
		log.Fatal(err)
	}

	allPicks := merge(existingPicks, newPicks)

	pvs, err := database.WeekPVS(year, week)
	if err != nil {
		log.Fatal(err)
	}

	// validate
	sevens, fives, threes := validate(allPicks, pvs)
	if !(sevens && fives && threes) {
		var b bytes.Buffer
		b.WriteString("Too Many:")

		if !sevens {
			b.WriteString(" Sevens")
		}

		if !fives {
			b.WriteString(" Fives")
		}

		if !threes {
			b.WriteString(" Threes")
		}

		context.Set(r, "error", b.String())
		picksForm(templatesDir, w, r)
		return
	}

	// otherwise update the database
	for _, p := range newPicks {
		err := database.MakePick(time.Now(), p.Id, p.Selection, p.Points)
		if err == database.ErrGameLocked {
			log.Println("Blocked locked game from being update")
		} else if err != nil {
			log.Fatal(err)
		}
	}

	selectedPicks, err := database.UserSelectedPicksByWeek(user, year, week)
	if err != nil {
		log.Fatal(err)
	}

	notifier.Notify(user, week, selectedPicks)

	context.Set(r, "success", "Picks submitted successfully!")
	picksForm(templatesDir, w, r)
}

func validate(picks []Pick, allowed database.PVS) (sevens, fives, threes bool) {
	var seven, five, three int

	for _, p := range picks {
		switch p.Points {
		case 7:
			seven++
		case 5:
			five++
		case 3:
			three++
		}
	}

	return seven <= allowed.Seven, five <= allowed.Five, three <= allowed.Three
}

func parsePicksForm(r *http.Request) ([]Pick, error) {
	r.ParseForm()

	pickedGames := r.Form["ids"]

	picks := make([]Pick, 0, len(pickedGames))

	for _, id := range pickedGames {
		selectionStr := r.FormValue(fmt.Sprintf("%s-Selection", id))
		if selectionStr == "" {
			return nil, fmt.Errorf("incomplete picks form: %s-Selection", id)
		}

		selection, err := strconv.ParseInt(selectionStr, 10, 32)
		if err != nil {
			return nil, err
		}

		if selection == -1 {
			continue
		}

		pointsStr := r.FormValue(fmt.Sprintf("%s-Points", id))
		if selectionStr == "" {
			return nil, fmt.Errorf("incomplete picks form: %s-Points", id)
		}

		points, err := strconv.ParseInt(pointsStr, 10, 32)
		if err != nil {
			return nil, err
		}

		id_int, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			return nil, err
		}

		p := Pick{
			Id:        id_int,
			Selection: int(selection),
			Points:    int(points),
		}

		picks = append(picks, p)
	}

	return picks, nil
}

func merge(existing []database.SelectedPicks, new []Pick) []Pick {
	all := make([]Pick, 0)

	for _, n := range new {
		all = append(all, n)
	}

	for _, e := range existing {
		found := false

		for _, n := range new {
			if n.Id == e.Id {
				found = true
				break
			}
		}

		if !found {
			all = append(all, Pick{e.Id, e.Selection, e.Points})
		}
	}

	return all
}
