package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Credentials json struct
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func getCredentials(r *http.Request) (*Credentials, error) {
	creds := new(Credentials)
	log.Println("DEBUG: GETCREDENTIALS()")
	err := json.NewDecoder(r.Body).Decode(creds)
	log.Println(creds.Username)
	return creds, err
}

func signin(w http.ResponseWriter, creds Credentials) (err error) {
	db := connectToDB(dbName)
	defer db.Close()
	user, err := getUserByName(db, creds.Username)

	if err != nil || !passwordsEqual(user.passwordHash, []byte(creds.Password)) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString, expirationTime, err := buildNewToken(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	addCookie(w, sessionTokenName, tokenString, expirationTime)
	return nil
}

func welcome(w http.ResponseWriter, r *http.Request) (user *User, err error) {
	// We can obtain the session token from the requests cookies, which come with every request
	claims, err := getCookie(w, r)
	if err != nil {
		if err == http.ErrNoCookie {
			// w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// w.WriteHeader(http.StatusBadRequest)
		return
	}
	db := connectToDB(dbName)
	defer db.Close()
	user, err = getUserByID(db, claims.userID)
	// w.Write([]byte(fmt.Sprintf("Welcome %s!", user.login)))
	return user, err
}

func renewCookie(w http.ResponseWriter, r *http.Request) (err error) {
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
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	addCookie(w, sessionTokenName, tokenString, expirationTime)
	return nil
}

func getCookie(w http.ResponseWriter, r *http.Request) (claims *Claims, err error) {
	c, err := r.Cookie(tokenName)
	if err != nil {
		return
	}

	// Check cookie validity
	tknStr := c.Value
	claims = &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	// Invalid token
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		// Invalid token Signature
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// Unknown token problem
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
