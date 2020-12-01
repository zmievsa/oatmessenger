package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err = addUser(db, creds.Username, creds.Email, creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = signin(w, *creds)
	if err != nil {
		fmt.Fprintf(w, "We have successfully registered you but weren't able to sign you in. (err: %s)\n", err.Error())
	}
	fmt.Fprintf(w, "You successfully registered, %s!\n", creds.Username)
	log.Printf("Successfully registered %s.\n", creds.Username)
}

func handleCheckCookieExistence(w http.ResponseWriter, r *http.Request) {
	log.Println("checkCookieExistence()")
	_, _, err := parseCookieFromHeaders(r)
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
func handleFindUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("findUsers()")
	user, _, err := parseCookieFromHeaders(r)
	if err != nil {
		log.Println("No cookie found.")
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	jsonData := new(JsonName)

	log.Println("Decoding json...")
	err = json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println("Error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := connectToDB(dbName)
	defer db.Close()

	users := searchUsersByName(db, jsonData.Name, user.ID)
	encodedJSON, err := json.Marshal(users)
	log.Printf("Users: %s, Users json: %s, \nError: %s \n", users, encodedJSON, err)
	json.NewEncoder(w).Encode(users)
}

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	t, _ := msToTime("0")
	cookie := http.Cookie{
		Name:     sessionTokenName,
		Value:    "",
		Secure:   true,
		Expires:  t,
		SameSite: http.SameSiteNoneMode,
		Domain:   "127.0.0.1",
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
}

func handleGetAllDialogues(w http.ResponseWriter, r *http.Request) {
	log.Println("getAllDialogues()")
	user, _, err := parseCookieFromHeaders(r)
	if err != nil {
		log.Println("No cookie found.")
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	db := connectToDB(dbName)
	defer db.Close()
	users := []*User{}
	if user.Dialogues != "" {
		userIDs, _ := getUserDialogues(user)
		for _, id := range userIDs {
			dialogueUser, _ := getUserByID(db, id)
			users = append(users, dialogueUser)
		}
	}
	json.NewEncoder(w).Encode(users)
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	log.Println("getUser()")
	user, cookie, err := parseCookieFromHeaders(r)
	if err != nil {
		log.Println("No cookie found.")
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	m := map[string]interface{}{"user": user, "cookie": cookie.Value}
	json.NewEncoder(w).Encode(m)
}

func handleSetFullname(w http.ResponseWriter, r *http.Request) {
	log.Println("setFullname()")
	jsonData := new(JsonFullName)

	log.Println("Decoding json...")
	err := json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println("No json arguments supplied.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, _, err := parseCookieFromHeaders(r)
	if err != nil {
		log.Println("No cookie found.")
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	db := connectToDB(dbName)
	defer db.Close()
	err = setUserFullName(db, user.ID, jsonData.Name)
}

func handleGetMessages(w http.ResponseWriter, r *http.Request) {
	log.Println("getMessages()")
	user, _, err := parseCookieFromHeaders(r)
	if err != nil {
		log.Println("No cookie found.")
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	jsonData := new(JsonGetMessages)

	log.Println("Decoding json...")
	err = json.NewDecoder(r.Body).Decode(&jsonData)
	if err != nil {
		log.Println("No json arguments supplied.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := connectToDB(dbName)
	defer db.Close()
	messages, err := getMessagesByUserIds(db, user.ID, jsonData.UserWithID)
	log.Println(err)
	err = json.NewEncoder(w).Encode(messages)
	log.Println(err)
}

func handleGetAnotherUser(w http.ResponseWriter, r *http.Request) {
	log.Println("getUser()")
	_, _, err := parseCookieFromHeaders(r)
	if err != nil {
		log.Println("No cookie found.")
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	jsonData := new(JsonUserID)
	err = json.NewDecoder(r.Body).Decode(jsonData)
	if err != nil {
		log.Println("No json arguments supplied.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	db := connectToDB(dbName)
	defer db.Close()
	anotherUser, err := getUserByID(db, jsonData.UserID)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No user with id=%d found.", jsonData.UserID)
		return
	}

	json.NewEncoder(w).Encode(anotherUser)
}

func main() {
	initSecret()

	hub := newHub()
	go hub.run()

	mux := http.NewServeMux()
	mux.HandleFunc("/register/", handleRegistration)
	mux.HandleFunc("/login/", handleLogin)
	mux.HandleFunc("/checkCookieExistence/", handleCheckCookieExistence)
	mux.HandleFunc("/findUsers/", handleFindUsers)
	mux.HandleFunc("/logout/", handleLogout)
	mux.HandleFunc("/getAllDialogues/", handleGetAllDialogues)
	mux.HandleFunc("/getUser/", handleGetUser)
	mux.HandleFunc("/getAnotherUser/", handleGetAnotherUser)
	mux.HandleFunc("/setFullName/", handleSetFullname)
	mux.HandleFunc("/getMessages/", handleGetMessages)
	mux.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
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
