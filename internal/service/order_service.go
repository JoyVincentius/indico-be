package service

import (
	"context"

	"indico-be/internal/models"
	"indico-be/internal/repository"
)

type OrderService struct {
	repo repository.OrderRepository
}

func NewOrderService(r repository.OrderRepository) *OrderService {
	return &OrderService{repo: r}
}

// PlaceOrder tries to reduce stock atomically, then creates the order.
// Returns error "OUT_OF_STOCK" when stock insufficient.
func (s *OrderService) PlaceOrder(ctx context.Context, req *models.Order) error {
	if err := s.repo.ReduceStock(ctx, req.ProductID, req.Quantity); err != nil {
		return err
	}
	return s.repo.Create(ctx, req)
}

// GetOrder returns a stored order by its ID.
func (s *OrderService) GetOrder(ctx context.Context, id uint64) (*models.Order, error) {
	return s.repo.GetByID(ctx, id)
}
