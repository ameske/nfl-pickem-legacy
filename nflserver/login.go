package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/ameske/nfl-pickem/database"
	"github.com/gorilla/context"
)

// Login processes the form post, determining whether or not the user succssfully logged in. If the login
// was a success the user is redirected to their desired endpoint, if such an endpoint exists. Otherwise,
// the user is taken to the login page.
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		n := context.Get(r, "next")
		if n == nil {
			Fetch("login.html").Execute(w, "", []string{"", "/login"})
		} else {
			next := n.(string)
			next64 := base64.StdEncoding.EncodeToString([]byte(next))
			Fetch("login.html").Execute(w, "", []string{"", fmt.Sprintf("/login?next=%s", string(next64))})
		}
		return
	}

	r.ParseForm()
	next := r.FormValue("next")

	// Attempt login, taking the user back to the login page with an error message if failed
	u := r.FormValue("username")
	p := r.FormValue("password")
	if !database.CheckCredentials(db, u, p) {
		if next == "" {
			Fetch("login.html").Execute(w, "", []string{"Invalid username or password.", "/login"})
		} else {
			Fetch("login.html").Execute(w, "", []string{"Invalid username of password.", fmt.Sprintf("/login?next=%s", next)})
		}
		return
	}

	// Set session information
	session, _ := store.Get(r, "LoginState")
	session.Values["status"] = "loggedin"
	session.Values["user"] = u
	session.Save(r, w)

	// Redirect appropriately
	if next == "" {
		http.Redirect(w, r, "/", 302)
	} else {
		n, err := base64.StdEncoding.DecodeString(next)
		if err != nil {
			log.Fatalf("Decoding next paramter: %s", err.Error())
		}
		http.Redirect(w, r, string(n), 302)
	}
}

// Logout clears the session information, which effectively logs the user out.
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "LoginState")
	delete(session.Values, "status")
	delete(session.Values, "user")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
