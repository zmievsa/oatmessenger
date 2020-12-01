package main

// Credentials json struct
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JsonName struct {
	Name string `json:"name"`
}

type JsonFullName struct {
	Name string `json:"fullname"`
}

type JsonGetMessages struct {
	UserWithID int `json:"user_with_id"`
	Offset     int `json:"offset"`
}

type JsonUserID struct {
	UserID int `json:"userID"`
}

type MessageJson struct {
	UserIDFor   int    `json:"userID_for"`
	Text        string `json:"text"`
	Attachments string `json:"attachments"`
}

type CookieJson struct {
	Cookie string `json:"cookie"`
}
