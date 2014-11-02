package main

import (
	"log"
	"net/http"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/ameske/go_nfl/database"
	"github.com/gorilla/context"
)

func LoginForm(w http.ResponseWriter, r *http.Request) {
	Fetch("login.html").Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Attempt login, taking the user back to the login page with an error message if failed
	u := r.FormValue("username")
	p := r.FormValue("password")
	if !database.CheckCredentials(db, u, p) {
		t := Fetch("login.html")
		t.Execute(w, "Invalid username or password.")
		return
	}

	// Set session information
	session, _ := store.Get(r, "LoginState")
	session.Values["status"] = "loggedin"
	session.Values["username"] = u
	session.Save(r, w)

	// Redirect to where they intended to go, or the home page if they explicitly were logging in
	n := context.Get(r, "next")
	if n == nil {
		http.Redirect(w, r, "/", 302)
		return
	}

	next := n.(string)
	if next == "/login" {
		http.Redirect(w, r, "/", 302)
		return
	} else {
		http.Redirect(w, r, next, 302)
		return
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "LoginState")
	delete(session.Values, "status")
	delete(session.Values, "username")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ChangePasswordForm(w http.ResponseWriter, r *http.Request) {
	t := Fetch("passwordChange.html").Execute(w, cpm{})
}

type cpm struct {
	Success string
	Error   string
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	u := r.FormValue("username")
	p := r.FormValue("oldPassword")
	pN := r.FormValue("newPassword")
	pNC := r.FormValue("confirmNewPassword")

	// Check that this user is actually who they claim they are
	if !database.CheckCredentials(db, u, p) {
		t := Fetch("passwordChange.html").Execute(w, "Invalid username or password.")
		return
	}

	// Make sure the user REALLY knows their new password and it isn't empty
	if pN != pNC || pN == "" {
		m := cpm{
			Error: "Passwords do not match.",
		}
		t := Fetch("passwordChange.html").Execute(w, m)
	}

	bpass, err := bcrypt.GenerateFromPassword([]byte(pN), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf(err.Error())
	}

	database.UpdatePassword(db, u, bpass)

	m := cpm{
		Success: "Password succesfully update!",
	}
	t := Fetch("passwordChange.html").Execute(w, m)
}
