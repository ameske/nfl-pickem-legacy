package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const stateLoggedIn = `
<h1>State</h1>
You are logged in
<form method="post" action="/logout">
  <button type="submit">Logout</button>
</form>
`

const stateLoggedOut = `
<h1>State</h1>
You are logged out
<form method="get" action="/">
  <button type="submit">Log Me In</button>
</form>
`

func LoginGet(w http.ResponseWriter, r *http.Request) {
	// Check to see if we have a valid existing session
	session, _ := store.Get(r, "LoginState")
	if session.Values["status"] == "loggedin" {
		http.Redirect(w, r, "/state", 302)
	} else {
		t, err := template.ParseFiles("templates/_base.html", "templates/navbar.html", "templates/login.html")
		if err != nil {
			log.Fatalf(err.Error())
		}
		t.Execute(w, nil)
	}
}

func LoginPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	u := r.FormValue("username")
	p := r.FormValue("password")
	if u != username || p != password {
		http.Error(w, "Invalid username and password", http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, "LoginState")
	session.Values["status"] = "loggedin"
	session.Save(r, w)

	http.Redirect(w, r, "/state", 302)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "LoginState")
	delete(session.Values, "status")
	session.Save(r, w)
	http.Redirect(w, r, "/state", http.StatusSeeOther)
}

func State(w http.ResponseWriter, r *http.Request) {
	// Check for an existing session, if none exists prompt for login
	session, err := store.Get(r, "LoginState")
	if err != nil {
		fmt.Printf("%s", err.Error())
	}

	if session.Values["status"] == "loggedin" {
		fmt.Printf("Session Status %s\n", session.Values["status"])
		w.Write([]byte(stateLoggedIn))
	} else {
		fmt.Printf("Session Status %s\n", session.Values["status"])
		w.Write([]byte(stateLoggedOut))
	}
}
