package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// ChatID is a pseudo-alias that allow future easy modification of Chat.ID
type ChatID string

func NewChatID() (ChatID, error) {
	return ChatID(uuid.New().String()), nil
}

// Validate check for ChatID integrity
func (id ChatID) Validate() bool {
	return true
}

// Equal check for ChatID equality
func (id ChatID) Equal(rhs ChatID) bool {
	return string(id) == string(rhs)
}
func (id ChatID) String() string {
	return string(id)
}

// Chat is a type that represent a Char between two Accounts
// ID could be later change by a UUID
type Chat struct {
	ID      ChatID `json:"id"`
	UserAID string `json:"user_a_id"`
	UserBID string `json:"user_b_id"`
}

// Validate check presence of minimal data required for an Chat.
func (c Chat) Validate() bool {
	return c.ID.Validate()
}

func (c Chat) String() string {
	// All info are present
	return fmt.Sprintf("%s %s %s", c.ID, c.UserAID, c.UserBID)
}
