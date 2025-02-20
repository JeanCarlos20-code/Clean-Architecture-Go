package web

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/repositories"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/use-cases/order"
	"github.com/JeanCarlos20-code/CleanArchitecture/pkg/events"
)

type WebOrderHandler struct {
	EventDispatcher   events.EventDispatcherInterface
	OrderRepository   repositories.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
}

func NewWebOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository repositories.OrderRepositoryInterface,
	OrderCreatedEvent events.EventInterface,
) *WebOrderHandler {
	return &WebOrderHandler{
		EventDispatcher:   EventDispatcher,
		OrderRepository:   OrderRepository,
		OrderCreatedEvent: OrderCreatedEvent,
	}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto repositories.OrderInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto.TypeRequisition = "rest"

	createOrder := order.NewCreateOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *WebOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}

	sort := r.URL.Query().Get("sort")
	listOrder := order.NewListOrderUseCase(h.OrderRepository)
	output, err := listOrder.Execute(pageInt, limitInt, sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
