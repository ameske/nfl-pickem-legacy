package main

import (
	"log"
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
	session, _ := store.Get(r, "LoginState")
	log.Printf("A user requested a protected endpoint: %s", r.URL.String())
	// Check if the user is logged in, if not redirect them to the login page
	if session.Values["status"] == "loggedin" {
		// We're logged in already, go to the index page instead
		if r.URL.String() == "/login" {
			http.Redirect(w, r, "/", 302)
		} else {
			p.h(w, r)
		}
	} else {
		context.Set(r, "next", r.URL.String())
		LoginForm(w, r)
	}
}
