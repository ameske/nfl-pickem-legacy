package main

import (
	"log"

	"github.com/ameske/nfl-pickem/database"
	"github.com/codegangsta/cli"

	"golang.org/x/crypto/bcrypt"
)

func inputUser(c *cli.Context) {
	first, last, email, password := c.String("first"), c.String("last"), c.String("email"), c.String("password")

	// Check to see if required arguments were given
	if first == "" {
		log.Fatalf("First name is required. Use --first <firstname> to specify.")
	}
	if last == "" {
		log.Fatalf("Last name is required. Use --last <lastname> to specify.")
	}
	if email == "" {
		log.Fatalf("Email is required. Use --email <address> to specify.")
	}
	if password == "" {
		log.Fatalf("Desired password required. Use --password <pass> to specify.")
	}

	bpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf(err.Error())
	}

	newUser := &database.Users{
		FirstName: first,
		LastName:  last,
		Email:     email,
		Password:  string(bpass),
	}

	err = db.Insert(newUser)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
