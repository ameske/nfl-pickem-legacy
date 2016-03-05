package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ameske/nfl-pickem/database"
	"github.com/gorilla/context"
)

func AdminPickForm(templatesDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			processAdminPicksForm(templatesDir, w, r)
		} else {
			renderAdminPicksForm(templatesDir, w, r)
		}
	}
}

func processAdminPicksForm(templatesDir string, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	ids := r.Form["ids"]

	for _, id := range ids {
		intID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		selection, err := strconv.ParseInt(r.FormValue(fmt.Sprintf("%s-Selection", id)), 10, 32)
		if err != nil {
			log.Fatal(err)
		}

		points, err := strconv.ParseInt(r.FormValue(fmt.Sprintf("%s-Points", id)), 10, 32)
		if err != nil {
			log.Fatal(err)
		}

		err = database.AdminMakePick(intID, int(selection), int(points))
		if err != nil {
			log.Fatal(err)
		}
	}

	context.Set(r, "success", "Picks submitted succesfully!")

	renderAdminPicksForm(templatesDir, w, r)
}

func renderAdminPicksForm(templatesDir string, w http.ResponseWriter, r *http.Request) {
	user, isAdmin := currentUser(r)
	if user == "" {
		http.Error(w, "no user information", http.StatusUnauthorized)
		return
	} else if !isAdmin {
		http.Error(w, "admin only", http.StatusForbidden)
		return

	}

	year, week := yearWeek(r)

	users, rows := database.AdminForm(year, week)

	var e, s string
	if context.Get(r, "error") != nil {
		e = context.Get(r, "error").(string)
	}
	if context.Get(r, "success") != nil {
		s = context.Get(r, "success").(string)
	}

	data := struct {
		Error           string
		Success         string
		URL             string
		Users           []string
		UserPicksByGame []database.AdminPickRow
	}{
		e,
		s,
		r.URL.String(),
		users,
		rows,
	}

	err := Fetch(templatesDir, "admin.html").Execute(w, user, isAdmin, data)
	if err != nil {
		log.Println(err)
	}
}
