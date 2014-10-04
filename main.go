package main

import (
	"log"
	"net/http"

	"github.com/ameske/go_nfl/database"
	"github.com/coopernurse/gorp"
	"github.com/gorilla/mux"
)

func main() {
	db := database.NflDb()

	r := mux.NewRouter()
	r.Handle("/", TestHandler(db))

	log.Fatal(http.ListenAndServe(":61389", r))
}

func TestHandler(db *gorp.DbMap) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
