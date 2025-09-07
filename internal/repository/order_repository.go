package repository

import (
	"context"
	"errors"

	"indico-be/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Order struct {
	models.Order
	// embed for GORM table name handling
}

type OrderRepository interface {
	Create(ctx context.Context, o *models.Order) error
	GetByID(ctx context.Context, id uint64) (*models.Order, error)
	ReduceStock(ctx context.Context, productID uint64, qty int) error
}

type orderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepository {
	return &orderRepo{db: db}
}

// Create a new order row.
func (r *orderRepo) Create(ctx context.Context, o *models.Order) error {
	return r.db.WithContext(ctx).Create(o).Error
}

// Get order by primary key.
func (r *orderRepo) GetByID(ctx context.Context, id uint64) (*models.Order, error) {
	var o models.Order
	if err := r.db.WithContext(ctx).First(&o, id).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

// ReduceStock uses SELECT … FOR UPDATE to guarantee atomicity.
func (r *orderRepo) ReduceStock(ctx context.Context, productID uint64, qty int) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var stock struct {
		Qty int `gorm:"column:stock"`
	}

	// SELECT ... FOR UPDATE
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Table("products").
		Where("id = ?", productID).
		First(&stock).Error; err != nil {
		tx.Rollback()
		return err
	}

	if stock.Qty < qty {
		tx.Rollback()
		return errors.New("OUT_OF_STOCK")
	}

	// UPDATE stock = stock - qty
	if err := tx.Table("products").
		Where("id = ?", productID).
		UpdateColumn("stock", gorm.Expr("stock - ?", qty)).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}