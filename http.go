package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
)

type ProtectedEndpoint struct {
	h http.HandlerFunc
}

// Protect wraps a HandlerFunc and only allows access to the handler func if specific criteria
// are met
func Protect(h http.HandlerFunc) *ProtectedEndpoint {
	return &ProtectedEndpoint{h: h}
}

// ServeHTTP implements the http.Handler interface for protected endpoints. Specifically,
// it acts as the gatekeeper for all incoming requests that require a user to be logged in.
func (p *ProtectedEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "LoginState")
	if err != nil {
		message := fmt.Sprintf("Application error: Please e-mail me letting me know the follwing information. Thank you! \n\n%s", err.Error())
		w.Write([]byte(message))
	}

	// Check if the user is logged in, if not redirect them to the login page
	if session.Values["status"] == "valid" {
		var userId int64
		err = db.SelectOne(&userId, "SELECT id FROM users WHERE email = $1", session.Values["username"])
		context.Set(r, "userId", userId)
		// We're logged in already, go to the index page instead
		if r.URL.String() == "/login" {
			http.Redirect(w, r, "/", 302)
		} else {
			p.h(w, r)
		}
	} else {
		context.Set(r, "next", r.URL.String())
		http.Redirect(w, r, "/login", 302)
	}
}
