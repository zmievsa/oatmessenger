package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

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

	if err != nil || !passwordsEqual(user.PasswordHash, []byte(creds.Password)) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString := buildNewToken(user)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	addCookie(w, sessionTokenName, tokenString)
	w.WriteHeader(http.StatusOK)
	return nil
}

// func welcome(w http.ResponseWriter, r *http.Request) (user *User, err error) {
// 	// We can obtain the session token from the requests cookies, which come with every request
// 	user, err := parseCookie(w, r)
// 	if err != nil {
// 		if err == http.ErrNoCookie {
// 			// w.WriteHeader(http.StatusUnauthorized)
// 			return
// 		}
// 		// w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
// 	db := connectToDB(dbName)
// 	defer db.Close()
// 	user, err = getUserByID(db, getUserIDFromClaims(claims))
// 	// w.Write([]byte(fmt.Sprintf("Welcome %s!", user.login)))
// 	return user, err
// }

// func renewCookie(w http.ResponseWriter, r *http.Request) (err error) {
// 	claims, err := parseCookie(w, r)
// 	if err != nil {
// 		return
// 	}

// 	// Now, create a new token for the current use, with a renewed expiration time
// 	// expirationTime := time.Now().Add(15 * time.Minute)
// 	// (*claims)["nbf"] = NEW NBF
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString(hmacSecret)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	addCookie(w, sessionTokenName, tokenString)
// 	return nil
// }

func parseCookie(w http.ResponseWriter, r *http.Request) (user *User, err error) {
	c, err := r.Cookie(sessionTokenName)
	if err != nil {
		return
	}

	// Check cookie validity
	tknStr := c.Value
	claims := &jwt.MapClaims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})

	// Invalid token
	if err != nil || !tkn.Valid {

		signedString, _ := tkn.SignedString(hmacSecret)
		return nil, fmt.Errorf("err: %s, Invalid cookie token: %s", err, signedString)
	}
	if err != nil {
		// Invalid token Signature
		if err == jwt.ErrSignatureInvalid {
			return
		}
		// Unknown token problem
		return
	}
	db := connectToDB(dbName)
	defer db.Close()
	fmt.Println((*claims)["userid"])
	user, err = getUserByID(db, getUserIDFromClaims(claims))
	if err != nil {
		log.Printf("User with ID %d not found.\n", getUserIDFromClaims(claims))
		return
	}
	log.Printf("Found user with ID from the cookie! (Name: %s, ID: %d)\n", user.Login, user.ID)
	return
}

// addCookie will apply a new cookie to the response of a http request
// with the key/value specified and will apply Set-Cookie header
func addCookie(w http.ResponseWriter, name, value string) {
	log.Println("Setting a cookie")
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Domain:   "127.0.0.1",
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	log.Println(w.Header())
}
