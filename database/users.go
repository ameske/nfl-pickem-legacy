package database

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Id        int64
	FirstName string
	LastName  string
	Email     string
	Admin     bool
	LastLogin time.Time
	Password  string
}

func AddUser(u Users) error {
	bpass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf(err.Error())
	}
	u.Password = string(bpass)

	_, err = db.Exec("INSERT INTO users(first_name, last_name, email, admin, last_login, password) VALUES(?1, ?2 ?3, ?4, ?5, ?6)", u.FirstName, u.LastName, u.Email, u.Admin, u.LastLogin, u.Password)

	return err
}

func AllUsers() []Users {
	users := make([]Users, 0)
	rows, err := db.Query("SELECT id, first_name, last_name, email, admin, last_login, password FROM users ORDER BY first_name ASC")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		tmp := Users{}
		err := rows.Scan(&tmp.Id, &tmp.FirstName, &tmp.LastName, &tmp.Email, &tmp.Admin, &tmp.LastLogin, &tmp.Password)
		if err != nil {
			log.Fatal(err)
		}
	}

	rows.Close()

	return users
}

func UserId(username string) int64 {
	var userId int64

	row := db.QueryRow("SELECT id FROM users WHERE email = ?1", username)
	err := row.Scan(&userId)

	if err != nil {
		log.Fatal(err)
	}

	return userId
}

func CheckCredentials(user string, password string) bool {
	var storedPassword string
	row := db.QueryRow("SELECT password FROM users WHERE email = ?1", user)
	err := row.Scan(&storedPassword)
	if err != nil {
		log.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))

	return err == nil
}

func UpdatePassword(user string, newPassword []byte) error {
	_, err := db.Exec("UPDATE users SET password = ?1 WHERE email = ?2", string(newPassword), user)
	return err
}

func UsersMap(users []Users) map[int64]Users {
	um := make(map[int64]Users)
	for _, u := range users {
		u := u
		um[u.Id] = u
	}

	return um
}
