package database

import (
	"log"
	"time"

	"github.com/coopernurse/gorp"

	"code.google.com/p/go.crypto/bcrypt"
)

type Users struct {
	Id        int64     `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Admin     bool      `db:"admin"`
	LastLogin time.Time `db:"last_login"`
	Password  string    `db:"password"`
}

func CheckCredentials(db *gorp.DbMap, user string, password string) bool {
	var u Users

	_ = db.SelectOne(&u, "SELECT * FROM users WHERE email = $1", user)
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}

func UpdatePassword(db *gorp.DbMap, user string, newPassword []byte) {
	var u Users
	err := db.SelectOne(&u, "SELECT * FROM users WHERE email = $1", user)
	if err != nil {
		log.Fatalf("UpdatePassword: %s", err.Error())
	}

	u.Password = string(newPassword)

	_, err = db.Update(&u)
	if err != nil {
		log.Fatalf("UpdatePassword: %s", err.Error())
	}
}
