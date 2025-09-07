package models

import "time"

type Transaction struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	MerchantID  uint64    `json:"merchant_id"`
	AmountCents int64     `json:"amount_cents"`
	FeeCents    int64     `json:"fee_cents"`
	Status      string    `json:"status"`
	PaidAt      time.Time `json:"paid_at"`
}
