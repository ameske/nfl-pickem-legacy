package main

import (
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/ameske/nfl-pickem/database"
)

type Test struct {
	Name string
	Data int
}

// m contains information for the passwordChange template
type m struct {
	Error   string
	Success string
}

// ChangePassword processes the password change form, informing the user of any problems or success.
func ChangePassword(templateDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, a := currentUser(r)
		if u == "" {
			http.Error(w, "no user information", http.StatusUnauthorized)
			return
		}

		if r.Method == "GET" {
			Fetch(templateDir, "changePassword.html").Execute(w, u, a, m{})
			return
		}

		r.ParseForm()

		p := r.FormValue("oldPassword")
		pN := r.FormValue("newPassword")
		pNC := r.FormValue("confirmNewPassword")

		// Check that this user is actually who they claim they are
		ok, err := database.CheckCredentials(u, p)
		if err != nil {
			log.Fatal(err)
		} else if !ok {
			Fetch(templateDir, "changePassword.html").Execute(w, u, a, m{Error: "Invalid username or password"})
			return
		}

		// Make sure the user REALLY knows their new password and it isn't empty
		if pN != pNC || pN == "" {
			Fetch(templateDir, "changePassword.html").Execute(w, u, a, m{Error: "Passwords do not match."})
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

		Fetch(templateDir, "changePassword.html").Execute(w, u, a, m{Success: "Password updated successfully."})
	}
}
