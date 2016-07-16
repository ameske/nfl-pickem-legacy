package main

import (
	"html/template"
	"io"
	"log"
	"time"

	"github.com/ameske/nfl-pickem/database"
)

var funcs = template.FuncMap{
	"gametime": gametime,
}

func gametime(time time.Time) string {
	return time.Format("Mon Jan 2 2006 15:04")
}

type GoNflTemplate struct {
	t *template.Template
}

// Execute gathers data needed to construct the navbar, chains is to the current template, and
// writes the output to the provided io.Writer.
func (t *GoNflTemplate) Execute(w io.Writer, user string, admin bool, content interface{}) error {
	inSeason := true

	year, week, err := database.CurrentWeek(time.Now())
	if err == database.ErrOffseason {
		inSeason = false
		if database.PrevSeasonExists(time.Now().Year()) {
			week = 17
		} else {
			week = 0
		}
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
