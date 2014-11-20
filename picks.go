package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ameske/go_nfl/database"
)

// Picks fetches this week's picks for the current logged in user and renders the
// picks template with them.
func PicksForm(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "LoginState")
	user := session.Values["user"].(string)

	year, week := yearWeek(r)

	picks := database.FormPicks(db, user, year, week)
	log.Printf("Pulled %d picks for user %s", len(picks), user)

	for _, p := range picks {
		log.Printf("%v", p)
	}

	data := struct {
		URL   string
		Picks []database.FormPick
	}{
		r.URL.String(),
		picks,
	}

	Fetch("picks.html").Execute(w, user, data)
}

// ProcessPicks validates a user's picks, and then updates the current picks in the database
func ProcessPicks(w http.ResponseWriter, r *http.Request) {
	// Gather endpoint information
	r.ParseForm()

	picks := r.Form["ids"]

	// First, validate that the user has not broken the rules for the given week
	if !validate(picks, r) {
		w.Write([]byte("Valid Picks"))
		return
	}

	for _, p := range picks {
		id, _ := strconv.ParseInt(p, 10, 64)

		// Fetch the Pick in question
		var pick database.Picks
		err := db.SelectOne(&pick, "SELECT * FROM picks WHERE id = $1", id)
		if err != nil {
			log.Fatalf(err.Error())
		}

		// Fetch the selection and the Point Value
		selection, _ := strconv.ParseInt(r.FormValue(fmt.Sprintf("%s-Selection", p)), 10, 32)
		points, _ := strconv.ParseInt(r.FormValue(fmt.Sprintf("%s-Points", p)), 10, 32)

		pick.Selection = int(selection)
		pick.Points = int(points)

		_, err = db.Update(&pick)
		if err != nil {
			log.Fatalf("ProcessPicks: %s", err.Error())
		}
	}

	w.Write([]byte("Picks submitted successfully!"))
}

func validate(picks []string, r *http.Request) bool {
	year, week := yearWeek(r)

	weekId := database.WeekId(db, year, week)
	pvs := database.GetPvs(db, weekId)

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

	log.Printf("Found:\tSeven:%d\tFive:%d\tThree:%d\tOne:%d\t", seven, five, three, one)
	log.Printf("Expected:\tSeven:%d\tFive:%d\tThree:%d\tOne:%d\t", pvs.Seven, pvs.Five, pvs.Three, pvs.One)

	if one > pvs.One || three > pvs.Three || five > pvs.Five || seven > pvs.Seven {
		return false
	}

	return true
}
