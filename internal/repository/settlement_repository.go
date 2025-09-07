package repository

import (
	"context"

	"indico-be/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Settlement struct {
	models.Settlement
	// embed for GORM table name handling
}

type SettlementRepository interface {
	Upsert(ctx context.Context, s *models.Settlement) error
}

type settlementRepo struct {
	db *gorm.DB
}

func NewSettlementRepo(db *gorm.DB) SettlementRepository {
	return &settlementRepo{db: db}
}

// Upsert based on (merchant_id,date) unique key.
func (r *settlementRepo) Upsert(ctx context.Context, s *models.Settlement) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "merchant_id"}, {Name: "date"}},
		DoUpdates: clause.AssignmentColumns([]string{"gross_cents", "fee_cents", "net_cents", "txn_count"}),
	}).Create(s).Error
}
