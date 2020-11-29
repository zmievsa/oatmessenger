package main

import (
	"database/sql"
	"fmt"
)

// User struct from the database
type User struct {
	ID           int
	Login        string
	FullName     string
	PasswordHash []byte
	IsDisabled   bool
}

func (u *User) String() string {
	return fmt.Sprintf(`User:
	ID: %d
	login: %s
	fullname: %s
	passwordHash: %s
	isDisabled: %t`, u.ID, u.Login, u.FullName, u.PasswordHash, u.IsDisabled)
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func (u *User) scan(dbUser scanner) error {
	return dbUser.Scan(&u.ID, &u.Login, &u.FullName, &u.PasswordHash, &u.IsDisabled)
}

func getUserByName(db *sql.DB, name string) (*User, error) {
	dbUser := db.QueryRow(fmt.Sprintf("SELECT * FROM user WHERE login='%s' COLLATE NOCASE", name))
	var user *User = new(User)
	err := user.scan(dbUser)
	return user, err
}

func searchUsersByName(db *sql.DB, name string, searcherUserid int) (users []*User) {
	users = []*User{}
	q := fmt.Sprintf("SELECT * FROM user WHERE (login LIKE '%%%s%%' AND NOT (ID=%d)) COLLATE NOCASE", name, searcherUserid)
	rows, err := db.Query(q)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var user *User = new(User)

		err = user.scan(rows)
		if err != nil {
			return
		}

		users = append(users, user)

	}
	return

}

func getUserByID(db *sql.DB, id int) (*User, error) {
	dbUser := db.QueryRow(fmt.Sprintf("SELECT * FROM user WHERE ID=%d", id))
	var user *User = new(User)
	err := user.scan(dbUser)
	return user, err
}

func getUserByToken(db *sql.DB, tokenString string) (*User, error) {
	token, err := getTokenByData(db, tokenString)
	if err != nil {
		return nil, err
	}
	user, err := getUserByID(db, token.userID)
	return user, err
}

func addUser(db *sql.DB, login string, plaintextPassword string) error {
	hashedPassword := string(hashAndSalt([]byte(plaintextPassword)))
	_, err := db.Exec(fmt.Sprintf(
		`INSERT INTO "main"."user" ("login", "password_hash", "is_disabled")
		VALUES ('%s', '%s', '%d');`, login, hashedPassword, 0))
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	return err
}
