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

// LoginForm renders the login form, populating the form action with the desired destination
// if the user arrived at this endpoint by requesting a protected endpoint while not authenticated.
func LoginForm(w http.ResponseWriter, r *http.Request) {
	n := context.Get(r, "next")
	if n == nil {
		log.Printf("User did not specify a next endpoint.")
		Fetch("login.html").Execute(w, "", []string{"", "/login"})
	} else {
		next := n.(string)
		log.Printf("User wishes to go to %s after login", next)
		next64 := base64.StdEncoding.EncodeToString([]byte(next))
		Fetch("login.html").Execute(w, "", []string{"", fmt.Sprintf("/login?next=%s", string(next64))})
	}
}

// Login processes the form post, determining whether or not the user succssfully logged in. If the login
// was a success the user is redirected to their desired endpoint, if such an endpoint exists. Otherwise,
// the user is taken to the login page.
func Login(w http.ResponseWriter, r *http.Request) {
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
	delete(session.Values, "username")
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// m is contains information for the passwordChange template
type m struct {
	Error   string
	Success string
}

// ChangePasswordForm renders the change password template.
func ChangePasswordForm(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "LoginState")
	u := session.Values["user"].(string)
	Fetch("passwordChange.html").Execute(w, u, m{})
}

// ChangePassword processes the password change form, informing the user of any problems or success.
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	session, _ := store.Get(r, "LoginState")
	u := session.Values["user"].(string)

	p := r.FormValue("oldPassword")
	pN := r.FormValue("newPassword")
	pNC := r.FormValue("confirmNewPassword")

	// Check that this user is actually who they claim they are
	if !database.CheckCredentials(db, u, p) {
		Fetch("passwordChange.html").Execute(w, u, m{Error: "Invalid username or password"})
		return
	}

	// Make sure the user REALLY knows their new password and it isn't empty
	if pN != pNC || pN == "" {
		Fetch("passwordChange.html").Execute(w, u, m{Error: "Passwords do not match."})
	}

	// Perform the password update
	bpass, err := bcrypt.GenerateFromPassword([]byte(pN), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf(err.Error())
	}
	database.UpdatePassword(db, u, bpass)

	Fetch("passwordChange.html").Execute(w, u, m{Success: "Password updated successfully."})
}
