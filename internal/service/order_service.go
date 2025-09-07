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

func (s *OrderService) PlaceOrder(ctx context.Context, req *models.Order) error {
	if err := s.repo.ReduceStock(ctx, req.ProductID, req.Quantity); err != nil {
		return err
	}
	return s.repo.Create(ctx, req)
}

func (s *OrderService) GetOrder(ctx context.Context, id uint64) (*models.Order, error) {
	return s.repo.GetByID(ctx, id)
}
