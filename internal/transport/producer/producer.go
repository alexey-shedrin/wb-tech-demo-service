package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/alexey-shedrin/wb-tech-demo-service/internal/config"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/util"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func New(cfg *config.Config) *Producer {
	connStr := fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{connStr},
		Topic:        cfg.Kafka.Topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: int(kafka.RequireAll),
	})

	return &Producer{
		writer: writer,
	}
}

func (p *Producer) Close() {
	err := p.writer.Close()
	if err != nil {
		log.Printf("failed to close Kafka writer: %v", err)
	}
}

func (p *Producer) ProduceOrders(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			order := util.GenerateOrder()

			orderJSON, err := json.Marshal(order)
			if err != nil {
				log.Printf("failed to marshal order to JSON: %v", err)
				continue
			}

			msg := kafka.Message{
				Key:   []byte(order.OrderUID),
				Value: orderJSON,
			}

			err = p.writer.WriteMessages(ctx, msg)
			if err != nil {
				log.Printf("failed to write message to Kafka: %v", err)
			} else {
				log.Printf("successfully produced order with UID: %s", order.OrderUID)
			}

		case <-ctx.Done():
			log.Println("producer stopped")
			return
		}
	}
}
