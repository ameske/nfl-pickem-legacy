package main

import (
	"html/template"
	"io"
	"log"
	"sync"
	"time"

	"github.com/ameske/nfl-pickem/database"
)

/*
This file is built up from https://github.com/zeebo/gostbook/blob/master/template.go
*/

var (
	cache        = map[string]*template.Template{}
	cacheLock    sync.Mutex
	templatesDir string
)

var funcs = template.FuncMap{
	"reverse":  reverse,
	"gametime": gametime,
}

type GoNflTemplate struct {
	t *template.Template
}

func (t *GoNflTemplate) Execute(w io.Writer, user string, admin bool, content interface{}) error {
	_, week := database.CurrentWeek()

	weeks := make([]int, 0)
	for i := 0; i < week; i++ {
		weeks = append(weeks, i+1)
	}

	data := struct {
		Navbar struct {
			Name  string
			Admin bool
			Weeks []int
		}

		Content interface{}
	}{
		Navbar: struct {
			Name  string
			Admin bool
			Weeks []int
		}{user, admin, weeks},
		Content: content,
	}

	return t.t.Execute(w, data)
}

// Fetch returns the specified template, creating it and adding it to the
// map of cached templates if it has not yet been created
func Fetch(name string) *GoNflTemplate {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	if t, ok := cache[name]; ok {
		return &GoNflTemplate{t: t}
	}

	t := template.New("_base.html").Funcs(funcs)
	t = template.Must(t.ParseFiles(templatesDir+"_base.html", templatesDir+"navbar.html", templatesDir+name))
	cache[name] = t

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
