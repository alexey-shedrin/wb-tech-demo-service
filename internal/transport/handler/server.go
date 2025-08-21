package handler

import (
	"fmt"
	"net/http"

	"github.com/alexey-shedrin/wb-tech-demo-service/internal/config"
	"github.com/gorilla/mux"
)

func NewServer(cfg *config.Config, h *Handler) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/", h.ServeIndex).Methods(http.MethodGet)
	router.HandleFunc("/order/{uid}", h.GetOrderByUID).Methods(http.MethodGet)

	addrStr := fmt.Sprintf("%s:%d", cfg.HTTPServer.Host, cfg.HTTPServer.Port)
	return &http.Server{
		Addr:         addrStr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
	}
}
