package main

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"time"
)

const messageDatetimeLayout = "2006-01-02 15:04:05"

// Db message
type Message struct {
	ID          int
	UserIDFrom  int
	UserIDFor   int
	Text        string
	Attachments string
	Datetime    time.Time
}

func (m *Message) String() string {
	return fmt.Sprintf(`User:
	ID: %d
	UserIDFrom: %d
	UserIDFor: %d
	Text: %s
	attachments: %s
	datetime: %s`, m.ID, m.UserIDFrom, m.UserIDFor, m.Text, m.Attachments, m.Datetime)
}
func (m *Message) scan(dbMsg scanner) error {
	var rawDatetime string
	err := dbMsg.Scan(&m.ID, &m.UserIDFrom, &m.UserIDFor, &m.Text, &m.Attachments, &rawDatetime)
	m.Datetime, _ = time.Parse(messageDatetimeLayout, rawDatetime)
	return err
}

func getMessagesByUserIds(db *sql.DB, userID1 int, userID2 int) (messages []*Message, err error) {
	log.Println("getMessagesByUserIds()")
	messages = []*Message{}
	q := fmt.Sprintf("SELECT * FROM message WHERE (user_id_from=%d AND user_id_for=%d) OR (user_id_from=%d AND user_id_for=%d)", userID1, userID2, userID2, userID1)
	rows, err := db.Query(q)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var message *Message = new(Message)

		err = message.scan(rows)
		if err != nil {
			return
		}
		messages = append(messages, message)
	}
	sort.SliceStable(messages, func(i, j int) bool {
		return messages[i].Datetime.Before(messages[j].Datetime)
	})
	return

}

func addMessage(db *sql.DB, userIDFrom int, userIDFor int, text string, attachments string, datetime time.Time) (message *Message, err error) {
	createUserDialogue(db, userIDFrom, userIDFor)
	time := datetime.Format(messageDatetimeLayout)
	result, err := db.Exec(fmt.Sprintf(
		`INSERT INTO "main"."message" ("user_id_from", "user_id_for", "text", "attachments", "datetime") VALUES (%d, %d, '%s', '%s', '%s');`, userIDFrom, userIDFor, text, attachments, time))
	if err != nil {
		fmt.Printf("Failed to insert message %s into database: %s\n", text, err.Error())
		return
	} else {
		var lastInsertID int64
		lastInsertID, err = result.LastInsertId()
		message = &Message{
			int(lastInsertID),
			userIDFrom,
			userIDFor,
			text,
			attachments,
			datetime,
		}
		return
	}
}
