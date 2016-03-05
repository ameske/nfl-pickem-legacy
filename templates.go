package main

import (
	"html/template"
	"io"
	"log"
	"time"

	"github.com/ameske/nfl-pickem/database"
)

/*
This file is built up from https://github.com/zeebo/gostbook/blob/master/template.go
*/

var funcs = template.FuncMap{
	"reverse":  reverse,
	"gametime": gametime,
}

type GoNflTemplate struct {
	t *template.Template
}

func (t *GoNflTemplate) Execute(w io.Writer, user string, admin bool, content interface{}) error {
	inSeason := true

	year, week, err := database.CurrentWeek(time.Now())
	if err == database.ErrOffseason {
		inSeason = false
		week = 17
	} else if err != nil {
		log.Fatal(err)
	}

	weeks := make([]int, 0, week)
	for i := 1; i <= week; i++ {
		weeks = append(weeks, i)
	}

	data := struct {
		Navbar struct {
			Name     string
			Admin    bool
			Weeks    []int
			Year     int
			InSeason bool
		}

		Content interface{}
	}{
		Navbar: struct {
			Name     string
			Admin    bool
			Weeks    []int
			Year     int
			InSeason bool
		}{user, admin, weeks, year, inSeason},
		Content: content,
	}

	return t.t.Execute(w, data)
}

// Fetch constructs the requested template, applying the common base and navbar templates
func Fetch(templatesDir, name string) *GoNflTemplate {
	t := template.New("_base.html").Funcs(funcs)
	t = template.Must(t.ParseFiles(templatesDir+"_base.html", templatesDir+"navbar.html", templatesDir+name))

	return &GoNflTemplate{t: t}
}

// reverse takes a string representing a named route and returns a url for said route
func reverse(name string, params ...string) string {
	url, err := router.GetRoute(name).URL(params...)
	if err != nil {
		log.Fatalf("reverse: %s", err.Error())
	}

	return url.Path
}

func gametime(time time.Time) string {
	return time.Format("Mon Jan 2 2006 15:04")
}
