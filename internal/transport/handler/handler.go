package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/alexey-shedrin/wb-tech-demo-service/internal/model"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/util"
	"github.com/gorilla/mux"
)

type Service interface {
	GetOrderByUID(ctx context.Context, orderUID string) (*model.Order, error)
}

type Handler struct {
	service Service
}

func New(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./web/index.html")
}

func (h *Handler) GetOrderByUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["uid"]

	ctx := r.Context()
	order, err := h.service.GetOrderByUID(ctx, orderUID)
	if err != nil {
		if errors.Is(err, util.ErrOrderNotFound) {
			http.Error(w, "Order not found", http.StatusNotFound)
			log.Printf("order with UID %s not found", orderUID)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("failed to retrieve order with UID %s: %v", orderUID, err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
	log.Printf("order with UID %s retrieved successfully", orderUID)
}
