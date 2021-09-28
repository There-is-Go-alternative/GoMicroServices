package domain

import (
	"fmt"
	"github.com/google/uuid"
)

// AccountID is a pseudo-alias that allow future easy modification of Account.ID
type AccountID string

func NewAccountID() (AccountID, error) {
	return AccountID(uuid.New().String()), nil
}

// Validate check for AccountID integrity
func (id AccountID) Validate() bool {
	return true
}

// Equal check for AccountID equality
func (id AccountID) Equal(rhs AccountID) bool {
	return string(id) == string(rhs)
}
func (id AccountID) String() string {
	return string(id)
}

// Account is a type that represent a user account
// ID could be later change by a UUID
type Account struct {
	ID        AccountID `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Admin     bool      `json:"admin,omitempty"`
	Address   Address   `json:"address,omitempty"`
}

// Validate check presence of minimal data required for an Account.
func (a Account) Validate() bool {
	// TODO: fix
	//return a.ID.Validate() && a.Email != "" && (a.Admin || a.Address.Validate())
	return a.ID.Validate() && a.Email != ""
}

func (a Account) String() string {
	// All info are present
	if a.Firstname != "" && a.Lastname != "" {
		return fmt.Sprintf("%s %s (%s)", a.Firstname, a.Lastname, a.Email)
	}
	// User did not fill its lastname
	if a.Firstname != "" {
		return fmt.Sprintf("%s (%s)", a.Firstname, a.Email)
	}
	// User did not fill its lastname and firstname
	return a.Email
}

// IsAdmin check if a user is website administrator.
func (a Account) IsAdmin() bool {
	return a.Admin
}

// Address holds information about location.
type Address struct {
	Country       string `json:"country,omitempty"`
	State         string `json:"state,omitempty"`
	City          string `json:"city,omitempty"`
	Street        string `json:"street,omitempty"`
	StreetNumber  int    `json:"street_number,omitempty"`
	Complementary string `json:"complementary,omitempty"`
}

// Validate check presence of minimal data required for an Address.
func (a Address) Validate() bool {
	return a.Country != "" && a.City != "" && a.Street != "" && a.StreetNumber < 0
}
