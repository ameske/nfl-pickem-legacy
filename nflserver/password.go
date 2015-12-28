package main

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/ameske/nfl-pickem/database"
)

// m contains information for the passwordChange template
type m struct {
	Error   string
	Success string
}

// ChangePassword processes the password change form, informing the user of any problems or success.
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		u, a := currentUser(r)
		Fetch("changePassword.html").Execute(w, u, a, m{})
		return
	}

	r.ParseForm()
	u, a := currentUser(r)

	p := r.FormValue("oldPassword")
	pN := r.FormValue("newPassword")
	pNC := r.FormValue("confirmNewPassword")

	// Check that this user is actually who they claim they are
	if !database.CheckCredentials(u, p) {
		Fetch("changePassword.html").Execute(w, u, a, m{Error: "Invalid username or password"})
		return
	}

	// Make sure the user REALLY knows their new password and it isn't empty
	if pN != pNC || pN == "" {
		Fetch("changePassword.html").Execute(w, u, a, m{Error: "Passwords do not match."})
		return
	}

	// Perform the password update
	bpass, err := bcrypt.GenerateFromPassword([]byte(pN), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	err = database.UpdatePassword(u, bpass)
	if err != nil {
		log.Fatal(err)
	}

	Fetch("changePassword.html").Execute(w, u, a, m{Success: "Password updated successfully."})
}
