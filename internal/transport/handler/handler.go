package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexey-shedrin/wb-tech-demo-service/internal/model"
	"github.com/alexey-shedrin/wb-tech-demo-service/internal/util"
	"github.com/gorilla/mux"
)

type Service interface {
	GetOrderByUID(orderUID string) (*model.Order, error)
}

type Handler struct {
	service Service
}

func New(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetOrderByUID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderUID := vars["uid"]

	order, err := h.service.GetOrderByUID(orderUID)
	if err != nil {
		if errors.Is(err, util.ErrOrderNotFound) {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}
