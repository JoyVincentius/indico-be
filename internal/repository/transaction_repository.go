package repository

import (
	"context"
	"fmt"
	"time"

	"indico-be/internal/models"

	"gorm.io/gorm"
)

type Transaction struct {
	models.Transaction
	// embed for GORM table name handling
}

type TransactionRepository interface {
	// Batch fetches rows ordered by ID (or any offset) – useful for pagination.
	FetchBatch(ctx context.Context, offset, limit int) ([]models.Transaction, error)
	CountAll(ctx context.Context) (int64, error)
	CountByPeriod(ctx context.Context, from, to time.Time) (int64, error)
	GetBatch(ctx context.Context, from, to time.Time, offset, limit int) ([]Transaction, error)
}

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) TransactionRepository {
	return &transactionRepo{db: db}
}

// Fetch a slice of transactions – deterministic ordering by ID.
func (r *transactionRepo) FetchBatch(ctx context.Context, offset, limit int) ([]models.Transaction, error) {
	var txs []models.Transaction
	err := r.db.WithContext(ctx).
		Order("id").
		Offset(offset).
		Limit(limit).
		Find(&txs).Error
	return txs, err
}

func (r *transactionRepo) CountAll(ctx context.Context) (int64, error) {
	var cnt int64
	err := r.db.WithContext(ctx).Model(&models.Transaction{}).Count(&cnt).Error
	return cnt, err
}

func (r *transactionRepo) CountByPeriod(ctx context.Context, from, to time.Time) (int64, error) {
	var count int64

	// Gunakan Query builder untuk mendapatkan jumlah transaksi dalam rentang
	err := r.db.WithContext(ctx).
		Model(&Transaction{}).
		Where("paid_at >= ? AND paid_at <= ?", from, to).
		Count(&count).
		Error

	if err != nil {
		return 0, fmt.Errorf("gagal menghitung transaksi: %w", err)
	}

	return count, nil
}

func (r *transactionRepo) GetBatch(ctx context.Context, from, to time.Time, offset, limit int) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.WithContext(ctx).
		Model(&Transaction{}).
		Where("paid_at >= ? AND paid_at <= ?", from, to).
		Order("paid_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).
		Error

	if err != nil {
		return nil, fmt.Errorf("gagal mengambil batch transaksi: %w", err)
	}

	return transactions, nil
}
