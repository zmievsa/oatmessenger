// Everything related to tokens (cookies)

package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const sessionTokenName = "session_oatmessenger_token"

var hmacSecret []byte

// Token struct from the database
type Token struct {
	userID int
	data   string
}

func (t *Token) String() string {
	return fmt.Sprintf(`Token:
	user_id: %d
	data: %s`, t.userID, t.data)
}

func (t *Token) scan(dbToken *sql.Row) error {
	return dbToken.Scan(&t.userID, &t.data)
}

// YOU ABSOLUTELY HAVE TO CALL THIS FUNCTION BEFORE EXECUTION
func initSecret() {
	dat, err := ioutil.ReadFile("secret.txt")
	panicIfError(err)
	hmacSecret = dat
}

func buildNewToken(user *User) string {
	log.Println("buildNewToken() for User:")
	log.Println(user)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	claims["iat"] = time.Now().Unix()
	claims["userid"] = user.ID
	token.Claims = claims
	fmt.Println("Token is valid: ", token.Valid)

	tokenString, _ := token.SignedString(hmacSecret)

	return tokenString
}

func getTokenByData(db *sql.DB, tokenString string) (*Token, error) {
	dbToken := db.QueryRow(fmt.Sprintf("SELECT * FROM token WHERE data='%s'", tokenString))
	var token *Token = new(Token)
	err := token.scan(dbToken)
	return token, err
}

func getUserIDFromClaims(claims *jwt.MapClaims) int {
	return (int)((*claims)["userid"].(float64))
}
