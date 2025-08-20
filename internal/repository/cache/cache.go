package cache

import (
	"log"
	"sync"

	"github.com/alexey-shedrin/wb-tech-demo-service/internal/model"
)

type Cache struct {
	data map[string]model.Order //[order_uid]model.Order
	mu   sync.Mutex
}

func New() *Cache {
	return &Cache{
		data: make(map[string]model.Order),
		mu:   sync.Mutex{},
	}
}

func (c *Cache) GetOrderByUID(orderUID string) (*model.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	order, exists := c.data[orderUID]
	if !exists {
		return nil, false
	}

	return &order, true
}

func (c *Cache) PutOrder(order *model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[order.OrderUID] = *order
	log.Printf("order cached: %v. cache size: %v", order.OrderUID, len(c.data))
}
