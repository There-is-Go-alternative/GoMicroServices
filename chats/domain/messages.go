package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// MessageID is a pseudo-alias that allow future easy modification of Message.ID
type MessageID string

func NewMessageID() (MessageID, error) {
	return MessageID(uuid.New().String()), nil
}

// Validate check for MessageID integrity
func (id MessageID) Validate() bool {
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
	ID        MessageID `json:"id"`
	ChatID    string    `json:"chat_id"`
	Content   string    `json:"content"`
	SenderID  string    `json:"sender_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate check presence of minimal data required for an Message.
func (m Message) Validate() bool {
	return m.ID.Validate()
}

func (m Message) String() string {
	// All info are present
	return fmt.Sprintf("%s %s %s %d-%02d-%02dT%02d:%02d:%02d-00:00\n",
		m.ID, m.SenderID, m.Content, m.CreatedAt.Year(), m.CreatedAt.Month(),
		m.CreatedAt.Day(), m.CreatedAt.Hour(), m.CreatedAt.Minute(),
		m.CreatedAt.Second())
}
