package domain

import (
	"time"

	"github.com/google/uuid"
)

func NewTransactionID() (*string, error) {
	randUUID, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}
	id := randUUID.String()
	return &id, nil
}

type Transaction struct {
	Id       string    `json:"id"`
	SellerId string    `json:"seller_id"`
	BuyerId  string    `json:"buyer_id"`
	AdId     string    `json:"ad_id"`
	Date     time.Time `json:"date"`
	Price    float64   `json:"price"`
}

type Account struct {
	Id string `json:"id"`
}

type IsAccountAdmin struct {
	IsAdmin bool `json:"is_admin"`
}

type Ad struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       uint     `json:"price"`
	Pictures    []string `json:"pictures"`
	UserId      string   `json:"owner_user_id"`
}

type Funds struct {
	ID          string    `json:"id"`
	UserId      string    `json:"user_id"`
	Balance     float64   `json:"balance"`
	LastUpdated time.Time `json:"last_updated"`
}
