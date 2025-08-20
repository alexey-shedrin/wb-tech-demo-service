package service

import "github.com/alexey-shedrin/wb-tech-demo-service/internal/model"

type Repository interface {
	GetOrderByUID(orderUID string) (*model.Order, error)
	SaveOrder(order *model.Order) error
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s Service) GetOrderByUID(orderUID string) (*model.Order, error) {
	order, err := s.repo.GetOrderByUID(orderUID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s Service) SaveOrder(order *model.Order) error {
	return s.repo.SaveOrder(order)
}
