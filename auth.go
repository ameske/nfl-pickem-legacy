package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/ameske/go_nfl/database"
	"github.com/gorilla/context"
)

func LoginForm(w http.ResponseWriter, r *http.Request) {
	n := context.Get(r, "next")

	if n == nil {
		log.Printf("User did not specify a next endpoint.")
		Fetch("login.html").Execute(w, []string{"", "/login"})
	} else {
		next := n.(string)
		log.Printf("User wishes to go to %s after login", next)
		next64 := base64.StdEncoding.EncodeToString([]byte(next))
		Fetch("login.html").Execute(w, []string{"", fmt.Sprintf("/login?next=%s", string(next64))})
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	next := r.FormValue("next")

	// Attempt login, taking the user back to the login page with an error message if failed
	u := r.FormValue("username")
	p := r.FormValue("password")
	if !database.CheckCredentials(db, u, p) {
		if next == "" {
			Fetch("login.html").Execute(w, []string{"Invalid username or password.", "/login"})
		} else {
			Fetch("login.html").Execute(w, []string{"Invalid username of password.", fmt.Sprintf("/login?next=%s", next)})
		}
		return
	}

	// Set session information
	session, _ := store.Get(r, "LoginState")
	session.Values["status"] = "loggedin"
	session.Values["user"] = u
	session.Save(r, w)

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

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "LoginState")
	delete(session.Values, "status")
	delete(session.Values, "username")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ChangePasswordForm(w http.ResponseWriter, r *http.Request) {
	Fetch("passwordChange.html").Execute(w, cpm{})
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
		Fetch("passwordChange.html").Execute(w, "Invalid username or password.")
		return
	}

	// Make sure the user REALLY knows their new password and it isn't empty
	if pN != pNC || pN == "" {
		m := cpm{
			Error: "Passwords do not match.",
		}
		Fetch("passwordChange.html").Execute(w, m)
	}

	bpass, err := bcrypt.GenerateFromPassword([]byte(pN), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf(err.Error())
	}

	database.UpdatePassword(db, u, bpass)

	m := cpm{
		Success: "Password succesfully update!",
	}
	Fetch("passwordChange.html").Execute(w, m)
}
