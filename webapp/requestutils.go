package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func currentUser(r *http.Request) string {
	session, _ := store.Get(r, "LoginState")
	user := session.Values["user"]

	if user == nil || user == "" {
		return ""
	} else {
		return user.(string)
	}
}

func yearWeek(r *http.Request) (int, int) {
	v := mux.Vars(r)
	y, _ := strconv.ParseInt(v["year"], 10, 32)
	w, _ := strconv.ParseInt(v["week"], 10, 32)
	return int(y), int(w)
}
