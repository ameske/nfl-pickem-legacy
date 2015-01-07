package main

import (
	"log"
	"net/http"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/ameske/go_nfl/database"
)

// m is contains information for the passwordChange template
type m struct {
	Error   string
	Success string
}

// ChangePasswordForm renders the change password template.
func ChangePasswordForm(w http.ResponseWriter, r *http.Request) {
	u := currentUser(r)
	Fetch("changePassword.html").Execute(w, u, m{})
}

// ChangePassword processes the password change form, informing the user of any problems or success.
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	u := currentUser(r)

	p := r.FormValue("oldPassword")
	pN := r.FormValue("newPassword")
	pNC := r.FormValue("confirmNewPassword")

	log.Printf("User: %s\nOldPassword: %s\nNewPassword: %s\nNewConfirmPassowrd: %s", u, p, pN, pNC)
	// Check that this user is actually who they claim they are
	if !database.CheckCredentials(db, u, p) {
		Fetch("changePassword.html").Execute(w, u, m{Error: "Invalid username or password"})
		return
	}

	// Make sure the user REALLY knows their new password and it isn't empty
	if pN != pNC || pN == "" {
		Fetch("changePassowrd.html").Execute(w, u, m{Error: "Passwords do not match."})
	}

	// Perform the password update
	bpass, err := bcrypt.GenerateFromPassword([]byte(pN), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf(err.Error())
	}
	database.UpdatePassword(db, u, bpass)

	Fetch("changePassword.html").Execute(w, u, m{Success: "Password updated successfully."})
}
