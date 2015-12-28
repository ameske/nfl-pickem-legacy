package main

import (
	"net/http"

	"github.com/gorilla/context"
)

// ServeHTTP implements the http.Handler interface for protected endpoints. Specifically,
// it acts as the gatekeeper for all incoming requests that require a user to be logged in.
func Protect(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "LoginState")
		if session.Values["status"] == "loggedin" {
			if r.URL.String() == "/login" {
				http.Redirect(w, r, "/", 302)
			} else {
				h(w, r)
			}
		} else {
			context.Set(r, "next", r.URL.String())
			http.Redirect(w, r, "/login", 302)
		}
	}
}

func AdminOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, isAdmin := currentUser(r)
		if !isAdmin {
			http.Error(w, "404 not found", http.StatusNotFound)
			return
		}

		h(w, r)
	}
}
