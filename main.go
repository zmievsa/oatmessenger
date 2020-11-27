package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

const cookieName = "oatmessenger_auth"

func handleLogin(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	username, password := getCredentials(r)
	user, err := getUserByName(db, username)
	if err != nil {
		fmt.Fprintln(w, "User not found")
		return
	}
	if comparePasswords(user.passwordHash, []byte(password)) {
		fmt.Fprintln(w, "You have successfully logged in")
	} else {
		fmt.Fprintln(w, "Incorrect password.")
	}

}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

// addCookie will apply a new cookie to the response of a http request
// with the key/value specified.
func addCookie(w http.ResponseWriter, name, value string, ttl time.Duration) {
	expire := time.Now().Add(ttl)
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expire,
	}
	http.SetCookie(w, &cookie)
}

func getCredentials(r *http.Request) (username string, password string) {
	r.ParseForm()
	return r.Form["username"][0], r.Form["password"][0]
}

func handleRegistration(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	username, password := getCredentials(r)
	_, err := getUserByName(db, username)
	if err == nil {
		fmt.Fprintf(w, "User %s already exists.\n", username)
		return
	}
	ip := GetIP(r)
	err = addUser(db, username, password, ip)
	user, err := getUserByName(db, username)
	fmt.Println(user)
	token := createToken(user)
	addCookie(w, cookieName, token, 48*time.Hour)
	fmt.Fprintln(w, "You successfully registered!")
}

func getAuthToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, err
}

func hasAuthToken(r *http.Request) bool {
	_, err := getAuthToken(r)
	return err == nil
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	// t, _ := template.ParseFiles("html/index.html")
	var fileNameToServe string
	if hasAuthToken(r) {
		fileNameToServe = "html/index.html"
	} else {
		fileNameToServe = "html/auth.html"
	}
	http.ServeFile(w, r, fileNameToServe)
}

// func main() {
// 	http.HandleFunc("/", handleRoot)
// 	http.HandleFunc("/register/", handleRegistration)
// 	http.HandleFunc("/login/", handleLogin)
// 	db := connectToDB("oatmessenger.db")
// 	defer db.Close()
// 	fmt.Println("You can now access the web app at http://127.0.0.1:8080")
// 	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
// }
