package repository

import (
	"log"

	"github.com/alexey-shedrin/wb-tech-demo-service/internal/config"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/model"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/repository/cache"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/repository/postgres"
)

type Repository struct {
	Postgres *postgres.Postgres
	Cache    *cache.Cache
}

func New(cfg *config.Config) *Repository {
	repo := Repository{
		Postgres: postgres.New(&cfg.Postgres),
		Cache:    cache.New(),
	}

	//Прогрев кэша

	return &repo
}

func (repo *Repository) Close() {
	repo.Postgres.Close()
}

func (repo *Repository) SaveOrder(order *model.Order) error {
	err := repo.Postgres.SaveOrder(order)
	if err != nil {
		return err
	}

	log.Printf("order created: %v", order.OrderUID)
	return nil
}

func (repo *Repository) GetOrderByUID(orderUID string) (*model.Order, error) {
	order, found := repo.Cache.GetOrderByUID(orderUID)
	if found {
		return order, nil
	}

	order, err := repo.Postgres.GetOrderByUID(orderUID)
	if err != nil {
		return nil, err
	}

	repo.Cache.PutOrder(order)
	return order, nil
}
