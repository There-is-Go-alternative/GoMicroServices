package domain

import (
	account "github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	"github.com/google/uuid"
)

// ChatID is a pseudo-alias that allow future easy modification of Chat.ID
type ChatID string

func NewChatID() (ChatID, error) {
	return ChatID(uuid.New().String()), nil
}

// Validate check for ChatID integrity
func (id ChatID) Validate() bool {
	// Need to find a way to check if the UUID is valid
	return true
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
	ID       ChatID              `json:"id"`
	UsersIDs []account.AccountID `json:"users_ids"`
}

// Validate check presence of minimal data required for an Chat.
func (c Chat) Validate() bool {
	return c.ID.Validate()
}

func (c Chat) String() string {
	// All info are present
	final_str := string(c.ID)
	for _, elem := range c.UsersIDs {
		final_str = final_str + " " + string(elem)
	}
	return final_str
}
