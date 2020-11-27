package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

// User struct from the database
type User struct {
	ID           int
	login        string
	fullName     string
	passwordHash []byte
	isDisabled   bool
	IP           string
}

func (u *User) String() string {
	return fmt.Sprintf(`User:
	ID: %d
	login: %s
	fullname: %s
	passwordHash: %s
	isDisabled: %t
	IP: %s`, u.ID, u.login, u.fullName, u.passwordHash, u.isDisabled, u.IP)
}

func (u *User) scan(dbUser *sql.Row) error {
	return dbUser.Scan(&u.ID, &u.login, &u.fullName, &u.passwordHash, &u.isDisabled, &u.IP)
}

func getUserByToken(db *sql.DB, token string) (*User, error) {
	var user *User = &User{}
	split := strings.Split(token, tokenSeparator)
	// userHash := split[0]
	userID, err := strconv.Atoi(split[1])
	if err == nil {
		user, err = getUserByID(db, userID)
		// if user.passwordHash != userHash {
		// 	err = errors.New("Incorrect hash")
		// }
	}
	return user, err
}

func getUserByName(db *sql.DB, name string) (*User, error) {
	dbUser := db.QueryRow(fmt.Sprintf("SELECT * FROM user WHERE login='%s' COLLATE NOCASE", name))
	var user *User = new(User)
	err := user.scan(dbUser)
	return user, err
}

func getUserByID(db *sql.DB, id int) (*User, error) {
	dbUser := db.QueryRow(fmt.Sprintf("SELECT * FROM user WHERE ID=%d", id))
	var user *User = new(User)
	err := user.scan(dbUser)
	return user, err
}

func addUser(db *sql.DB, login string, plaintextPassword string, IP string) error {
	hashedPassword := string(hashAndSalt([]byte(plaintextPassword)))
	_, err := db.Exec(fmt.Sprintf(
		`INSERT INTO "main"."user" ("login", "password_hash", "ip", "is_disabled")
		VALUES ('%s', '%s', '%s', '%d');`, login, hashedPassword, IP, 0))
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	return err
}

func createToken(user *User) string {
	// return user.passwordHash + tokenSeparator + fmt.Sprint(user.ID)
	return "1" // TODO: stub
}
