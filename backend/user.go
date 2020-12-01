package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// User struct from the database
type User struct {
	ID           int
	Login        string
	FullName     string
	PasswordHash []byte
	Dialogues    string
}

func (u *User) String() string {
	return fmt.Sprintf(`User:
	ID: %d
	login: %s
	fullname: %s
	passwordHash: %s
	dialogues: %s`, u.ID, u.Login, u.FullName, u.PasswordHash, u.Dialogues)
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func (u *User) scan(dbUser scanner) error {
	return dbUser.Scan(&u.ID, &u.Login, &u.FullName, &u.PasswordHash, &u.Dialogues)
}

func getUserByName(db *sql.DB, name string) (*User, error) {
	dbUser := db.QueryRow(fmt.Sprintf("SELECT * FROM user WHERE login='%s' COLLATE NOCASE", name))
	var user *User = new(User)
	err := user.scan(dbUser)
	return user, err
}

func searchUsersByName(db *sql.DB, name string, searcherUserid int) (users []*User) {
	users = []*User{}
	q := fmt.Sprintf("SELECT * FROM user WHERE login LIKE '%%%s%%' AND NOT (ID=%d)", name, searcherUserid)
	log.Println(q)
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
		`INSERT INTO "main"."user" ("login", "password_hash")
		VALUES ('%s', '%s');`, login, hashedPassword))
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	return err
}

func setUserFullName(db *sql.DB, userID int, newName string) error {
	q := fmt.Sprintf("UPDATE user SET full_name = '%s' WHERE (ID=%d)", newName, userID)
	_, err := db.Exec(q)
	return err
}

func getUserDialogues(user *User) (userIDs []int, err error) {
	userIdsAsStrings := strings.Split(user.Dialogues, ";")
	// Convert str arr to int arr
	for _, idAsStr := range userIdsAsStrings {
		var idAsInt int
		idAsInt, err = strconv.Atoi(idAsStr)
		if err != nil {
			return
		}
		userIDs = append(userIDs, idAsInt)
	}
	return
}

func usersHaveADialogue(db *sql.DB, user1 *User, user2 *User) bool {
	user1Dialogues, _ := getUserDialogues(user1)
	for _, id := range user1Dialogues {
		if id == user2.ID {
			return true
		}
	}
	return false
}

func addUserIDToDialogue(user1 *User, user2ID int) {
	user2IDAsString := fmt.Sprint(user2ID)
	if user1.Dialogues == "" {
		user1.Dialogues = user2IDAsString
	} else {
		user1.Dialogues = user1.Dialogues + ";" + user2IDAsString
	}
}

func createUserDialogue(db *sql.DB, user1ID int, user2ID int) error {
	user1, _ := getUserByID(db, user1ID)
	user2, _ := getUserByID(db, user2ID)
	if !usersHaveADialogue(db, user1, user2) {
		addUserIDToDialogue(user1, user2ID)
		addUserIDToDialogue(user2, user1ID)

		// user1 = "Pidor"
		q := fmt.Sprintf(`	UPDATE user SET 'dialogues' = CASE
								WHEN (id = %d) THEN '%s'
								WHEN (id = %d) THEN '%s'
   							END`, user1ID, user1.Dialogues, user2ID, user2.Dialogues)
		_, err := db.Exec(q)
		log.Println("SQL Query err: ", err)
		return err
	}
	return nil
}
