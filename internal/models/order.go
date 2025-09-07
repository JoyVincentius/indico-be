package models

import "time"

type Order struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	ProductID uint64    `json:"product_id"`
	BuyerID   string    `json:"buyer_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}
