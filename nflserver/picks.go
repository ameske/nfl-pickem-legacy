package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ameske/nfl-pickem/database"
	"github.com/gorilla/context"
)

// Picks fetches this week's picks for the current logged in user and renders the
// picks template with them.
func PicksForm(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)

	picks := database.FormPicks(user, false)

	e, s := "", ""
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
	}{
		e,
		s,
		r.URL.String(),
		picks,
	}

	Fetch("picks.html").Execute(w, user, data)
}

// ProcessPicks validates a user's picks, and then updates the current picks in the database
func ProcessPicks(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)

	r.ParseForm()
	pickedGames := r.Form["ids"]

	if three, five, seven := validate(user, pickedGames, r); !(three && five && seven) {
		var message bytes.Buffer
		message.WriteString("Invalid Picks: ")
		if !three {
			message.WriteString("Too many three point games. ")
		}
		if !five {
			message.WriteString("Too many five point games. ")
		}
		if !seven {
			message.WriteString("Too many seven point games. ")
		}

		context.Set(r, "error", message.String())
		PicksForm(w, r)
		return
	}

	// Update the picks in the database based on the user's selection, ignoring unselected picks
	for _, p := range pickedGames {
		id, _ := strconv.ParseInt(p, 10, 64)

		selection, _ := strconv.ParseInt(r.FormValue(fmt.Sprintf("%s-Selection", p)), 10, 32)
		if selection == 0 {
			continue
		}
		points, _ := strconv.ParseInt(r.FormValue(fmt.Sprintf("%s-Points", p)), 10, 32)

		err := database.MakePick(int(id), int(selection), int(points))
		if err != nil {
			log.Fatalf("ProcessPicks: %s", err.Error())
		}
	}

	selectedPicks := database.FormPicks(user, true)
	_, week := database.CurrentWeek()
	SendPicksEmail(user,
		fmt.Sprintf("Current Week %d Picks", week),
		week,
		selectedPicks,
	)

	context.Set(r, "success", "Picks submitted successfully!")
	PicksForm(w, r)
}

// validate handles server side validation of the point distribution of a submitted
// point set. It allows users to "under-point" their picks to allow them to submit
// games at their leisure.
func validate(user string, updatedPicks []string, r *http.Request) (threes, fives, sevens bool) {
	pvs := database.WeekPvs()

	one := 0
	three := 0
	five := 0
	seven := 0

	currentPicks := make([]database.FormPick, 0, 16)

	// Add all of the locked old picks to the current "view" of the user's picks
	oldPicks := database.FormPicks(user, true)
	for _, op := range oldPicks {
		locked := true
		for _, up := range updatedPicks {
			pickId, _ := strconv.ParseInt(up, 10, 32)
			if op.Id == pickId {
				locked = false
			}
		}

		if locked {
			currentPicks = append(currentPicks, op)
		}
	}

	// Add all of the updated and selected picks to the current "view" of the user's picks
	for _, up := range updatedPicks {
		selection, _ := strconv.ParseInt(r.FormValue(fmt.Sprintf("%s-Selection", up)), 10, 32)
		if selection == 0 {
			continue
		}

		points, _ := strconv.ParseInt(r.FormValue(fmt.Sprintf("%s-Points", up)), 10, 32)
		currentPicks = append(currentPicks, database.FormPick{
			Points: int(points),
		})
	}

	// Verify that we still have a valid point distribution
	for _, p := range currentPicks {
		switch p.Points {
		case 1:
			one += 1
		case 3:
			three += 1
		case 5:
			five += 1
		case 7:
			seven += 1
		}
	}

	return three > pvs.Three, five > pvs.Five, seven > pvs.Seven
}
