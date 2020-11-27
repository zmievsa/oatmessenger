package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type credentials struct {
	password string `json:"password"`
	username string `json:"username"`
}

func signin(w http.ResponseWriter, r *http.Request) {
	var creds credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := connectToDB(dbName)
	defer db.Close()
	user, err := getUserByName(db, creds.username)

	if err != nil || !passwordsEqual(user.passwordHash, []byte(creds.password)) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString, expirationTime, err := buildNewToken(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	addCookie(w, sessionTokenName, tokenString, expirationTime)
}

func welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	claims, err := getCookie(w, r)
	if err != nil {
		return
	}
	db := connectToDB(dbName)
	defer db.Close()
	user, err := getUserByID(db, claims.userID)
	w.Write([]byte(fmt.Sprintf("Welcome %s!", user.login)))
}

func renewCookie(w http.ResponseWriter, r *http.Request) {
	claims, err := getCookie(w, r)
	if err != nil {
		return
	}
	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(15 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	addCookie(w, sessionTokenName, tokenString, expirationTime)
}

func getCookie(w http.ResponseWriter, r *http.Request) (claims *Claims, err error) {
	c, err := r.Cookie(tokenName)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims = &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	return
}

// addCookie will apply a new cookie to the response of a http request
// with the key/value specified.
func addCookie(w http.ResponseWriter, name, value string, expirationTime time.Time) {
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expirationTime,
	}
	http.SetCookie(w, &cookie)
}
