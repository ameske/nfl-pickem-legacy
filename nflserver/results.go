package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func Results(w http.ResponseWriter, r *http.Request) {
	year, week := yearWeek(r)
	u := currentUser(r)

	// This page changes throughout the day, so we can't cache it.
	name := fmt.Sprintf("/opt/ameske/gonfl/templates/%d-Week%d-Results.html", year, week)
	t := template.New("_base.html").Funcs(funcs)
	t = template.Must(t.ParseFiles("/opt/ameske/gonfl/templates/_base.html", "/opt/ameske/gonfl/templates/navbar.html", name))
	data := struct {
		User    string
		Content interface{}
	}{
		u,
		nil,
	}
	t.Execute(w, data)
}
