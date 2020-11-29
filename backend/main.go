package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
)

const cookieName = "oatmessenger_auth"

var allowedOrigins = []string{"http://127.0.0.1:8080", "http://127.0.0.1:8090"}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
func handleLogin(w http.ResponseWriter, r *http.Request) {
	// setupResponse(&w, r)
	creds, err := getCredentials(r)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = signin(w, *creds)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func handleRegistration(w http.ResponseWriter, r *http.Request) {
	// setupResponse(&w, r)
	log.Println("Received a new registration request.")
	log.Printf("Method: %s\n", r.Method)
	creds, err := getCredentials(r)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		log.Println("Bad request (incorrect Credentials).")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Username: %s, Password: %s.\n", creds.Username, creds.Password)
	db := connectToDB(dbName)
	defer db.Close()
	_, err = getUserByName(db, creds.Username)
	if err == nil {
		fmt.Fprintf(w, "User %s already exists.\n", creds.Username)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err = addUser(db, creds.Username, creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = signin(w, *creds)
	if err != nil {
		fmt.Fprintln(w, "We have successfully registered you but weren't able to sign you in.")
	}
	fmt.Fprintf(w, "You successfully registered, %s!\n", creds.Username)
	log.Printf("Successfully registered %s.\n", creds.Username)
}
func handleRoot(w http.ResponseWriter, r *http.Request) {
	// setupResponse(&w, r)
	log.Println("The root was pinged.")
	_, err := welcome(w, r)
	var fileNameToServe string
	if err == nil {
		fileNameToServe = "html/index.html"
	} else {
		fileNameToServe = "html/auth.html"
	}
	http.ServeFile(w, r, fileNameToServe)
}

func handleCheckCookieExistence(w http.ResponseWriter, r *http.Request) {
	log.Println("checkCookieExistence()")
	_, err := getCookie(w, r)
	if err != nil {
		log.Println("No cookie found.")
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
	} else {
		log.Println("Found a cookie!")
		w.WriteHeader(http.StatusOK)
	}
	return

}

func main() {
	initSecret()
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/register/", handleRegistration)
	mux.HandleFunc("/login/", handleLogin)
	mux.HandleFunc("/checkCookieExistence/", handleCheckCookieExistence)
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug:          true,
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	})
	handler := c.Handler(mux)
	fmt.Println("You can now access the web app at http://127.0.0.1:8090")
	log.Fatal(http.ListenAndServe(":8090", handler))
}
