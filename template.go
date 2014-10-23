package main

import (
	"html/template"
	"log"
	"path/filepath"
	"sync"
)

/*
This file is built up from https://github.com/zeebo/gostbook/blob/master/template.go
*/

var (
	cache     = map[string]*template.Template{}
	cacheLock sync.Mutex
)

var funcs = template.FuncMap{
	"reverse": reverse,
}

// Fetch returns the specified template, creating it and adding it to the
// map of cached templates if it has not yet been created
func Fetch(name string) *template.Template {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	if t, ok := cache[name]; ok {
		return t
	}

	t := template.New("_base.html").Funcs(funcs)
	t = template.Must(t.ParseFiles("templates/_base.html", "templates/navbar.html", filepath.Join("templates", name)))
	cache[name] = t

	return t
}

// reverse takes a string representing a named route and returns a url for said route
func reverse(name string, params ...string) string {
	url, err := router.GetRoute(name).URL(params...)
	if err != nil {
		log.Fatalf("reverse: %s", err.Error())
	}

	return url.Path
}
