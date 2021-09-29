package domain

import (
	"github.com/google/uuid"
)

// AccountID is a pseudo-alias that allow future easy modification of Account.ID
type AdID string

func NewAdID() (AdID, error) {
	return AdID(uuid.New().String()), nil
}

// Validate check for AccountID integrity
func (id AdID) Validate() bool {
	return true
}

// Equal check for AccountID equality
func (id AdID) Equal(rhs AdID) bool {
	return string(id) == string(rhs)
}
func (id AdID) String() string {
	return string(id)
}

// Account is a type that represent a user account
// ID could be later change by a UUID
type Ad struct {
	ID          AdID   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	Picture     string `json:"picture"`
}

// Validate check presence of minimal data required for an Account.
func (a Ad) Validate() bool {
	// TODO: fix
	return true
}

func (a Ad) String() string {
	//Todo
	return ("")
}
