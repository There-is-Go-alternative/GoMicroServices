package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type FundsID string

func NewFundsID() (*FundsID, error) {
	randUUID, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}
	id := FundsID(randUUID.String())
	return &id, nil
}

func (id FundsID) Validate() error {
	if _, err := uuid.Parse(string(id)); err != nil {
		return fmt.Errorf("ID (%v) is invalid: {%v}", id, err)
	}
	return nil
}

func (id FundsID) String() string {
	return string(id)
}

type Funds struct {
	ID          FundsID   `json:"id"`
	UserId      string    `json:"user_id"`
	Balance     int       `json:"balance"`
	LastUpdated time.Time `json:"last_updated"`
}

func (f Funds) Validate() error {
	if err := f.ID.Validate(); err != nil {
		return err
	}
	return nil
}
