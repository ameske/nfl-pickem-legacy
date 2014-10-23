package main

import (
	"html/template"
	"log"
	"net/http"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/ameske/go_nfl/database"
	"github.com/gorilla/context"
)

func checkCredentials(user string, password string) bool {
	var u database.Users
	_ = db.SelectOne(&u, "SELECT * FROM users WHERE email = $1", user)

	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}

func LoginForm(w http.ResponseWriter, r *http.Request) {
	// Check to see if we have a valid existing session
	session, _ := store.Get(r, "LoginState")
	if session.Values["status"] == "valid" {
		http.Redirect(w, r, "/state", 302)
	} else {
		t, err := template.ParseFiles("templates/_base.html", "templates/navbar.html", "templates/login.html")
		if err != nil {
			log.Fatalf(err.Error())
		}
		t.Execute(w, nil)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Attempt login, taking the user back to the login page with an error message if failed
	u := r.FormValue("username")
	p := r.FormValue("password")
	if !checkCredentials(u, p) {
		t, err := template.ParseFiles("templates/_base.html", "templates/navbar.html", "templates/login.html")
		if err != nil {
			log.Fatalf(err.Error())
		}
		t.Execute(w, "Invalid username or password.")
	}

	// Set the session and redirect to where they intended to go
	session, _ := store.Get(r, "LoginState")
	session.Values["status"] = "loggedin"
	session.Values["username"] = u
	session.Save(r, w)

	next := context.Get(r, "next").(string)
	http.Redirect(w, r, next, 302)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "LoginState")
	delete(session.Values, "status")
	delete(session.Values, "username")
	session.Save(r, w)
	http.Redirect(w, r, "/state", http.StatusSeeOther)
}
