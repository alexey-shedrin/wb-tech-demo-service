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

	cacheSize := cfg.Cache.StartupSize
	if cacheSize > 0 {
		log.Printf("starting cache warming with the last %d orders", cacheSize)
		lastOrders, err := repo.Postgres.GetLastOrders(cacheSize)
		if err != nil {
			log.Printf("cache warming failed: could not get last orders: %v", err)
		} else {
			for _, order := range lastOrders {
				repo.Cache.PutOrder(order)
			}
			log.Printf("cache warming finished: loaded %d orders into the cache", len(lastOrders))
		}
	}

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
