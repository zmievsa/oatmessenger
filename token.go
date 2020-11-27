package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const tokenName = "oatmessenger_token"
const sessionTokenName = "session_" + tokenName

var hmacSecret string

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

// Claims that are used as fields in the cookie token
type Claims struct {
	userID int
	jwt.StandardClaims
}

// YOU ABSOLUTELY HAVE TO CALL THIS FUNCTION BEFORE EXECUTION
func initSecret() {
	dat, err := ioutil.ReadFile("/tmp/dat")
	panicIfError(err)
	hmacSecret = string(dat)
}

func buildNewToken(user *User) (token string, expirationTime time.Time, err error) {
	expirationTime = time.Now().Add(30 * time.Minute)

	claims := &Claims{
		userID: user.ID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	token, err = unsignedToken.SignedString(hmacSecret)

	return token, expirationTime, err
}

func getTokenByData(db *sql.DB, tokenString string) (*Token, error) {
	dbToken := db.QueryRow(fmt.Sprintf("SELECT * FROM token WHERE data='%s'", tokenString))
	var token *Token = new(Token)
	err := token.scan(dbToken)
	return token, err
}
