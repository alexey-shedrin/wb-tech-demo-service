package service

import (
	"context"

	"github.com/alexey-shedrin/wb-tech-demo-service/internal/model"
)

type Repository interface {
	GetOrderByUID(ctx context.Context, orderUID string) (*model.Order, error)
	SaveOrder(ctx context.Context, order *model.Order) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) GetOrderByUID(ctx context.Context, orderUID string) (*model.Order, error) {
	order, err := s.repo.GetOrderByUID(ctx, orderUID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s Service) SaveOrder(ctx context.Context, order *model.Order) error {
	return s.repo.SaveOrder(ctx, order)
}
