package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func getCredentials(r *http.Request) (*JSONCredentials, error) {
	creds := new(JSONCredentials)
	log.Println("DEBUG: GETCREDENTIALS()")
	err := json.NewDecoder(r.Body).Decode(creds)
	log.Println(creds.Username)
	return creds, err
}

func signin(w http.ResponseWriter, creds JSONCredentials) (err error) {
	db := connectToDB(dbName)
	defer db.Close()
	user, err := getUserByName(db, creds.Username)

	if err != nil || !passwordsEqual(user.PasswordHash, []byte(creds.Password)) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokenString := buildNewToken(user)
	addCookie(w, sessionTokenName, tokenString)
	w.WriteHeader(http.StatusOK)
	return nil
}

func parseCookie(tknStr string) (user *User, err error) {

	// Check cookie validity
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

func queryParams(r *http.Request, attrName string) (string, error) {
	for k, v := range r.URL.Query() {
		if k == attrName {
			return v[0], nil
		}
	}
	return "", fmt.Errorf("No cookie found in %s", r.URL.String())

}

func parseCookieFromURL(r *http.Request) (user *User, err error) {
	cookie, err := queryParams(r, "cookie")
	if err != nil {
		log.Println(err)
		return
	}

	tknStr := cookie
	user, err = parseCookie(tknStr)
	return
}

func parseCookieFromHeaders(r *http.Request) (user *User, cookie *http.Cookie, err error) {
	cookie, err = r.Cookie(sessionTokenName)
	if err != nil {
		return
	}

	tknStr := cookie.Value
	user, err = parseCookie(tknStr)
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
