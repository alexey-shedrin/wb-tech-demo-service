package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexey-shedrin/wb-tech-demo-service/internal/config"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/repository"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/service"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/transport/handler"
	"github.com/gorilla/mux"
)

func Run() {
	cfg := config.New()
	log.Println("config initialized")

	repo := repository.New(cfg)
	defer repo.Close()
	log.Println("repository initialized")

	svc := service.New(repo)
	log.Println("service initialized")

	hdr := handler.New(svc)
	log.Println("handler initialized")

	router := mux.NewRouter()
	router.HandleFunc("/order/{uid}", hdr.GetOrderByUID).Methods(http.MethodGet)
	log.Println("router initialized")

	addrStr := fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	srv := &http.Server{
		Addr:         addrStr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
	}
	log.Println("server initialized")

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	log.Println(addrStr)
	log.Print("server is running...")
}
