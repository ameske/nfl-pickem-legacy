package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func currentUser(r *http.Request) (string, bool) {
	session, _ := store.Get(r, "LoginState")
	user := session.Values["user"]
	admin := session.Values["admin"]

	if user == nil || user == "" {
		return "", false
	} else {
		return user.(string), admin.(bool)
	}
}

func yearWeek(r *http.Request) (int, int) {
	v := mux.Vars(r)
	y, _ := strconv.ParseInt(v["year"], 10, 32)
	w, _ := strconv.ParseInt(v["week"], 10, 32)
	return int(y), int(w)
}
