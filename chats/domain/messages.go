package domain

import (
	"fmt"
	"time"

	chats "github.com/There-is-Go-alternative/GoMicroServices/chats/domain"

	"github.com/google/uuid"
)

// MessageID is a pseudo-alias that allow future easy modification of Message.ID
type MessageID string

func NewMessageID() (MessageID, error) {
	return MessageID(uuid.New().String()), nil
}

// Validate check for MessageID integrity
func (id MessageID) Validate() bool {
	// Need to find a way to check if the UUID is valid
	return true
}

// Equal check for MessageID equality
func (id MessageID) Equal(rhs MessageID) bool {
	return string(id) == string(rhs)
}
func (id MessageID) String() string {
	return string(id)
}

// Message is a type that represent a message sent by a user to a specific conversation
// ID could be later change by a UUID
type Message struct {
	ID           MessageID    `json:"id"`
	ChatID       chats.ChatID `json:"chat_id"`
	Content      string       `json:"content"`
	SenderID     string       `json:"sender_id"`
	CreatedAt    time.Time    `json:"created_at"`
	Attachements [][]byte     `json:"attachements"`
}

// Validate check presence of minimal data required for an Message.
func (m Message) Validate() bool {
	return m.ID.Validate() && (m.ChatID != "") && (m.Content != "") && (m.SenderID != "")
}

func (m Message) String() string {
	// All info are present
	return fmt.Sprintf("%s %s %s %s %s",
		m.ID, m.ChatID, m.SenderID, m.Content, m.CreatedAt.Format(time.UnixDate))
}
