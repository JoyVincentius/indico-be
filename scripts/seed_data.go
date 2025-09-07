package main

import (
	"log"
	"math/rand"
	"time"

	"indico-be/config"
	"indico-be/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("db err: %v", err)
	}
	// Ensure product row exists
	db.Exec(`INSERT INTO products (id, stock) VALUES (1, 100) ON DUPLICATE KEY UPDATE stock=100`)

	// Seed 1M transactions
	const total = 1_000_000
	merchants := 100
	start := time.Now().AddDate(0, -3, 0) // three months back

	for i := 0; i < total; i++ {
		tx := models.Transaction{
			MerchantID:  uint64(rand.Intn(merchants) + 1),
			AmountCents: int64(rand.Intn(10_000) + 100),
			FeeCents:    int64(rand.Intn(500)),
			Status:      "paid",
			PaidAt:      start.Add(time.Duration(rand.Int63n(int64(90 * 24 * time.Hour)))),
		}
		if err := db.Create(&tx).Error; err != nil {
			log.Fatalf("seed error at %d: %v", i, err)
		}
		if i%10_000 == 0 {
			log.Printf("seeded %d rows...", i)
		}
	}
	log.Println("seed completed")
}
