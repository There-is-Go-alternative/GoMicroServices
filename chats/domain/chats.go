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
func (id ChatID) Validate() error {
	if _, err := uuid.Parse(id.String()); err != nil {
		return fmt.Errorf("ID (%v) is invalid: {%v}", id, err)
	}
	return nil
}

// Equal check for ChatID equality
func (id ChatID) Equal(rhs ChatID) bool {
	return string(id) == string(rhs)
}
func (id ChatID) String() string {
	return string(id)
}

// Chat is a type that represents a Chat between multiple Accounts
// ID could be later change by a UUID
type Chat struct {
	ID       ChatID   `json:"id"`
	UsersIDs []string `json:"users_ids"`
}

// Validate check presence of minimal data required for an Chat.
func (c Chat) Validate() bool {
	return c.ID.Validate() != nil
}

func (c Chat) String() string {
	// All info are present
	final_str := string(c.ID)
	for _, elem := range c.UsersIDs {
		final_str = final_str + " " + string(elem)
	}
	return final_str
}
