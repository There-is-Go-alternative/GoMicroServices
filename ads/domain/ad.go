package domain

import (
	"github.com/google/uuid"
)

type AdID string

func NewAdID() (AdID, error) {
	return AdID(uuid.New().String()), nil
}

func (id AdID) Validate() bool {
	return true
}

func (id AdID) Equal(rhs AdID) bool {
	return string(id) == string(rhs)
}
func (id AdID) String() string {
	return string(id)
}

type Ad struct {
	ID          AdID   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
	Picture     string `json:"picture"`
}

func (a Ad) Validate() bool {
	// TODO: fix
	return true
}

func (a Ad) String() string {
	//Todo
	return ("")
}
