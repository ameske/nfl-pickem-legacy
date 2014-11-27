package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ameske/go_nfl/database"
	"github.com/gorilla/context"
)

// Picks fetches this week's picks for the current logged in user and renders the
// picks template with them.
func PicksForm(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)

	picks := database.FormPicks(db, user, false)

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
	picks := r.Form["ids"]

	if !validate(picks, r) {
		context.Set(r, "error", "Invalid picks")
		PicksForm(w, r)
		return
	}

	// Update the picks in the database based on the user's selection
	for _, p := range picks {
		id, _ := strconv.ParseInt(p, 10, 64)

		var pick database.Picks
		err := db.SelectOne(&pick, "SELECT * FROM picks WHERE id = $1", id)
		if err != nil {
			log.Fatalf(err.Error())
		}

		selection, _ := strconv.ParseInt(r.FormValue(fmt.Sprintf("%s-Selection", p)), 10, 32)
		points, _ := strconv.ParseInt(r.FormValue(fmt.Sprintf("%s-Points", p)), 10, 32)

		pick.Selection = int(selection)
		pick.Points = int(points)

		_, err = db.Update(&pick)
		if err != nil {
			log.Fatalf("ProcessPicks: %s", err.Error())
		}
	}

	selectedPicks := database.FormPicks(db, user, true)
	_, week := database.CurrentWeek(db)
	SendPicksEMail(user,
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
func validate(picks []string, r *http.Request) bool {
	pvs := database.WeekPvs(db)

	one := 0
	three := 0
	five := 0
	seven := 0

	for _, p := range picks {
		tmp, _ := strconv.ParseInt(r.FormValue(fmt.Sprintf("%s-Points", p)), 10, 32)
		switch tmp {
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

	if three > pvs.Three || five > pvs.Five || seven > pvs.Seven {
		return false
	}

	return true
}
