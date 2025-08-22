package sup

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexey-shedrin/wb-tech-demo-service/internal/config"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/transport/producer"
)

func Run() {
	cfg := config.New()
	log.Println("config initialized")

	prdcr := producer.New(cfg)
	defer prdcr.Close()
	log.Println("producer initialized")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		prdcr.ProduceOrders(ctx)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case <-ctx.Done():
	}
	log.Println("shutting down gracefully...")
}
