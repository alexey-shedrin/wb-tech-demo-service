package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexey-shedrin/wb-tech-demo-service/internal/config"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/repository"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/service"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/transport/consumer"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/transport/handler"
)

func Run() {
	cfg := config.New()
	log.Println("config initialized")

	repo := repository.New(cfg)
	defer repo.Close()
	log.Println("repository initialized")

	srvc := service.New(repo)
	log.Println("service initialized")

	cnsmr := consumer.New(cfg, srvc)
	defer cnsmr.Close()
	log.Println("consumer initialized")

	hndlr := handler.New(srvc)
	log.Println("handler initialized")

	server := handler.NewServer(cfg, hndlr)
	log.Println("server initialized")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		cnsmr.ConsumeOrders(ctx)
	}()
	log.Println("consumer is running...")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to start server: %v", err)
			cancel()
		}
	}()
	log.Println("server is running...")
	log.Printf("server is listening on %s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case <-ctx.Done():
	}
	log.Println("shutting down gracefully...")

	sdCtx, sdCancel := context.WithTimeout(context.Background(), 2*cfg.HTTPServer.Timeout)
	defer sdCancel()

	if err := server.Shutdown(sdCtx); err != nil {
		log.Printf("failed to shutdown server: %v", err)
	} else {
		log.Println("server stopped")
	}

	cancel()
	log.Println("application stopped")
}
