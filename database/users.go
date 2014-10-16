package database

import "time"

type Users struct {
	Id        int64     `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Admin     bool      `db:"admin"`
	LastLogin time.Time `db:"last_login"`
	Password  string    `db:"password"`
}

func Login() {

}
