package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

func main() {
	//db := database.NflDb()

	r := mux.NewRouter()
	r.HandleFunc("/", Login).Methods("POST", "GET")
	r.HandleFunc("/state", State).Methods("GET")
	r.HandleFunc("/logout", Logout).Methods("POST")

	log.Fatal(http.ListenAndServe(":61389", r))
}

// Hardcode login/cookie stuff for testing
const (
	username = "kyle"
	password = "password"
)

var s = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

const loginPage = `
<h1>Login</h1>
<form method="post" action="/">
    <label for="name">User name</label>
    <input type="text" id="username" name="username">
    <label for="password">Password</label>
    <input type="password" id="password" name="password">
    <button type="submit">Login</button>
</form>
`

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Write([]byte(loginPage))
		return
	}

	r.ParseForm()
	u := r.FormValue("username")
	p := r.FormValue("password")

	if u != username || p != password {
		http.Error(w, "Invalid username and password", http.StatusUnauthorized)
		return
	}

	// For now just store the username, later we'll chain this to a session with a
	// secret value
	value := map[string]string{
		"username": u,
	}

	if encoded, err := s.Encode("LoginState", value); err == nil {
		cookie := &http.Cookie{
			Name:  "LoginState",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}

	w.Write([]byte("You have successfully logged in!"))
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "LoginState",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}

func State(w http.ResponseWriter, r *http.Request) {
	// For now, just check if the cookie exists. Later we will check a session to see
	// if the secret value is set.
	if cookie, err := r.Cookie("LoginState"); err == nil {
		value := make(map[string]string)
		if err = s.Decode("LoginState", cookie.Value, &value); err == nil {
			fmt.Printf("Yay! COOKIE! It contains: %s", value["username"])
			w.Write([]byte("Yes, you are logged in"))
			return
		}
	}

	w.Write([]byte("You don't appear to be logged in"))
	return
}

func writeJsonResponse(w http.ResponseWriter, r interface{}) {
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Write(j)
}
