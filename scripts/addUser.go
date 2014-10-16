package main

import (
	"flag"
	"log"

	"github.com/ameske/go_nfl/database"

	"code.google.com/p/go.crypto/bcrypt"
)

var (
	first    = flag.String("first", "", "User's first name")
	last     = flag.String("last", "", "User's last name")
	email    = flag.String("email", "", "User's email address")
	admin    = flag.Bool("admin", false, "Is user admin")
	password = flag.String("password", "", "User's desired password")
)

func main() {
	flag.Parse()
	db := database.NflDb()

	// Check to see if required arguments were given
	if *first == "" {
		log.Fatalf("First name is required. Use --first <firstname> to specify.")
	}
	if *last == "" {
		log.Fatalf("Last name is required. Use --last <lastname> to specify.")
	}
	if *email == "" {
		log.Fatalf("Email is required. Use --email <address> to specify.")
	}
	if *password == "" {
		log.Fatalf("Desired password required. Use --password <pass> to specify.")
	}

	bpass, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf(err.Error())
	}

	newUser := &database.Users{
		FirstName: *first,
		LastName:  *last,
		Email:     *email,
		Admin:     *admin,
		Password:  string(bpass),
	}

	err = db.Insert(newUser)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
