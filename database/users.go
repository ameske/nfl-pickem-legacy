package database

import (
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

func IsAdmin(username string) (admin bool, err error) {
	row := db.QueryRow("SELECT admin FROM users WHERE email = ?1", username)
	err = row.Scan(&admin)
	if err != nil {
		return false, err
	}

	return
}

func AddUser(u Users) error {
	bpass, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bpass)

	_, err = db.Exec("INSERT INTO users(first_name, last_name, email, admin, last_login, password) VALUES(?1, ?2, ?3, ?4, ?5, ?6)", u.FirstName, u.LastName, u.Email, u.Admin, u.LastLogin, u.Password)

	return err
}

func UserFirstNames() ([]string, error) {
	rows, err := db.Query("SELECT first_name FROM users ORDER BY id ASC")
	if err != nil {
		return nil, err
	}

	users := make([]string, 0)
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			return nil, err
		}
		users = append(users, tmp)
	}
	rows.Close()

	return users, nil
}

func Usernames() ([]string, error) {
	rows, err := db.Query("SELECT email FROM users ORDER BY id ASC")
	if err != nil {
		return nil, err
	}

	users := make([]string, 0)
	for rows.Next() {
		var tmp string
		err := rows.Scan(&tmp)
		if err != nil {
			return nil, err
		}
		users = append(users, tmp)
	}
	rows.Close()

	return users, nil
}

func CheckCredentials(user string, password string) (bool, error) {
	var storedPassword string
	row := db.QueryRow("SELECT password FROM users WHERE email = ?1", user)
	err := row.Scan(&storedPassword)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))

	return err == nil, nil
}

func UpdatePassword(user string, newPassword []byte) error {
	_, err := db.Exec("UPDATE users SET password = ?1 WHERE email = ?2", string(newPassword), user)
	return err
}
