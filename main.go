package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func main() {
	//db := database.NflDb()

	r := mux.NewRouter()
	r.HandleFunc("/", LoginGet).Methods("GET")
	r.HandleFunc("/", LoginPost).Methods("POST")
	r.HandleFunc("/state", State).Methods("GET")
	r.HandleFunc("/logout", Logout).Methods("POST")

	log.Fatal(http.ListenAndServe(":61389", r))
}

// Hardcode login/cookie stuff for testing
const (
	username = "kyle"
	password = "password"
)

var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

const loginPage = `
<h1>Login</h1>
<form method="post" action="/">
    <label for="name">User name</label>
    <input type="text" id="username" name="username">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

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
	session, _ := store.Get(r, "LoginState")

	// Check for our logged in status
	if session.Values["status"] == "loggedin" {
		http.Redirect(w, r, "/state", 302)
	} else {
		w.Write([]byte(loginPage))
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
	session, _ := store.Get(r, "LoginState")

	if session.Values["status"] == "loggedin" {
		w.Write([]byte(stateLoggedIn))
	} else {
		w.Write([]byte(stateLoggedOut))
	}
}

func writeJsonResponse(w http.ResponseWriter, r interface{}) {
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(j)
}
