package models

import "time"

type Settlement struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	MerchantID  uint64    `json:"merchant_id"`
	Date        time.Time `json:"date"`
	GrossCents  int64     `json:"gross_cents"`
	FeeCents    int64     `json:"fee_cents"`
	NetCents    int64     `json:"net_cents"`
	TxnCount    int64     `json:"txn_count"`
	GeneratedAt time.Time `json:"generated_at"`
	RunID       string    `json:"run_id"`
}
