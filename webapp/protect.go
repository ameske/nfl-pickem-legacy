package main

import (
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
	if session.Values["status"] == "loggedin" {
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
