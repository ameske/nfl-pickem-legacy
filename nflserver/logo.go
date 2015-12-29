package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	logosDir string
)

func Logo(w http.ResponseWriter, r *http.Request) {
	team := mux.Vars(r)["team"]

	file := fmt.Sprintf("%s/%s.gif", logosDir, team)

	http.ServeFile(w, r, file)
}
