package domain

import (
	"fmt"
	"github.com/There-is-Go-alternative/GoMicroServices/account/internal/xerrors"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/mail"
	"time"
)

// AccountID is a pseudo-alias that allow future easy modification of Account.ID
type AccountID string

func NewAccountID() (*AccountID, error) {
	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	id := AccountID(randomUUID.String())
	return &id, nil
}

// Validate check for AccountID integrity
func (id AccountID) Validate() error {
	if _, err := uuid.Parse(string(id)); err != nil {
		return fmt.Errorf("ID (%v) is invalid: {%v}", id, err)
	}
	return nil
}

// Equal check for AccountID equality
func (id AccountID) Equal(rhs AccountID) bool {
	return string(id) == string(rhs)
}
func (id AccountID) String() string {
	return string(id)
}

func validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

// Account is a type that represent a user account
// ID could be later change by a UUID
type Account struct {
	ID        AccountID `json:"id,omitempty"`
	Email     string    `json:"email"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Admin     bool      `json:"admin,omitempty"`
	Address   Address   `json:"address,omitempty"`
	Balance   int       `json:"balance,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Validate check presence of minimal data required for an Account.
func (a Account) Validate() error {
	var err error
	var errs []error

	if err = a.ID.Validate(); err != nil {
		errs = append(errs, err)
	}
	if err = validateEmail(a.Email); err != nil {
		errs = append(errs, err)
	}

	return xerrors.Concat(errs...)
}

func (a Account) String() string {
	// All info are present
	if a.Firstname != "" && a.Lastname != "" && a.Email != "" {
		return fmt.Sprintf("%s %s (%s)", a.Firstname, a.Lastname, a.Email)
	}
	// User did not fill its lastname
	if a.Firstname != "" && a.Email != "" {
		return fmt.Sprintf("%s (%s)", a.Firstname, a.Email)
	}
	// All info are present
	if a.Firstname != "" && a.Lastname != "" {
		return fmt.Sprintf("%s %s", a.Firstname, a.Lastname)
	}
	// User did not fill its lastname
	if a.Firstname != "" {
		return a.Firstname
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
func (a Address) Validate() error {
	// TODO: Errlist
	if a.Country != "" && a.City != "" && a.Street != "" && a.StreetNumber > 0 {
		return nil
	}
	return errors.New("address not good formed")
}
