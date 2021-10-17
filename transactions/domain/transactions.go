package domain

import "time"

type Transaction struct {
	Id       string    `json:"id"`
	SellerId string    `json:"seller_id"`
	BuyerId  string    `json:"buyer_id"`
	Date     time.Time `json:"date"`
	Price    float64   `json:"price"`
}
