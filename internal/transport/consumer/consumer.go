package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/alexey-shedrin/wb-tech-demo-service/internal/config"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/model"
	"github.com/segmentio/kafka-go"
)

type Service interface {
	SaveOrder(ctx context.Context, order *model.Order) error
}

type Consumer struct {
	reader  *kafka.Reader
	service Service
}

func New(cfg *config.Config, service Service) *Consumer {
	connStr := fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{connStr},
		GroupID: cfg.Kafka.Group,
		Topic:   cfg.Kafka.Topic,
	})

	return &Consumer{
		reader:  reader,
		service: service,
	}
}

func (c *Consumer) Close() {
	err := c.reader.Close()
	if err != nil {
		log.Printf("failed to close Kafka reader: %v", err)
	}
}

func (c *Consumer) ConsumeOrders(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("consumer stopped")
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					log.Println("consumer stopped")
					return
				}
				log.Printf("consumer could not read message: %v", err)
				continue
			}

			var order model.Order
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				log.Printf("failed to unmarshal message: %v. Message content: %s", err, string(msg.Value))

				if err := c.reader.CommitMessages(ctx, msg); err != nil {
					log.Printf("failed to commit message after unmarshal error: %v", err)
				}
				continue
			}

			if order.OrderUID == "" {
				log.Printf("invalid order: missing order_uid. Message content: %s", string(msg.Value))

				if err := c.reader.CommitMessages(ctx, msg); err != nil {
					log.Printf("failed to commit message after validation error: %v", err)
				}
				continue
			}

			err = c.service.SaveOrder(ctx, &order)
			if err != nil {
				log.Printf("failed to save order %s: %v", order.OrderUID, err)
			} else {
				log.Printf("order %s saved successfully", order.OrderUID)
			}

			err = c.reader.CommitMessages(ctx, msg)
			if err != nil {
				log.Printf("failed to commit message: %v", err)
			}
		}
	}
}
