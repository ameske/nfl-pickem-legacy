package database

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
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

func NewUser(first, last, email, password string) error {
	u := Users{
		FirstName: first,
		LastName:  last,
		Email:     email,
		Admin:     false,
		LastLogin: time.Now(),
	}

	bpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf(err.Error())
	}
	u.Password = string(bpass)

	return db.Insert(&u)
}

func AddUser(u Users) error {
	bpass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf(err.Error())
	}
	u.Password = string(bpass)

	return db.Insert(&u)
}

func AllUsers() []Users {
	var users []Users
	_, err := db.Select(&users, "SELECT * from users ORDER BY first_name ASC")
	if err != nil {
		log.Fatalf("AllUsers: %s", err.Error())
	}

	return users
}

func UserId(username string) int64 {
	var userId int64
	err := db.SelectOne(&userId, "SELECT id FROM users WHERE email = ?1", username)
	if err != nil {
		log.Fatalf("UserId: %s", err.Error())
	}

	return userId
}

func CheckCredentials(user string, password string) bool {
	var u Users

	_ = db.SelectOne(&u, "SELECT * FROM users WHERE email = ?1", user)
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	return err == nil
}

func UpdatePassword(user string, newPassword []byte) {
	var u Users
	err := db.SelectOne(&u, "SELECT * FROM users WHERE email = ?1", user)
	if err != nil {
		log.Fatalf("UpdatePassword: %s", err.Error())
	}

	u.Password = string(newPassword)

	_, err = db.Update(&u)
	if err != nil {
		log.Fatalf("UpdatePassword: %s", err.Error())
	}
}

func UsersMap(users []Users) map[int64]Users {
	um := make(map[int64]Users)
	for _, u := range users {
		u := u
		um[u.Id] = u
	}

	return um
}
