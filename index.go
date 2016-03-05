package main

import (
	"log"
	"net/http"
)

func Index(templatesDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, a := currentUser(r)

		welcome := struct {
			User string
		}{u}

		err := Fetch(templatesDir, "index.html").Execute(w, u, a, welcome)
		if err != nil {
			log.Println(err)
		}
	}
}
