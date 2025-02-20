package repositories

import (
	"context"
	"time"

	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/entities"
)

type OrderInputDTO struct {
	Price           float64   `json:"price"`
	Tax             float64   `json:"tax"`
	IssueDate       time.Time `json:"issueDate"`
	TypeRequisition string    `json:"typeRequisition"`
}

type OrderOutputDTO struct {
	ID              string     `json:"id"`
	Price           float64    `json:"price"`
	Tax             float64    `json:"tax"`
	FinalPrice      float64    `json:"finalPrice"`
	IssueDate       time.Time  `json:"issueDate"`
	TypeRequisition string     `json:"typeRequisition"`
	DeleteAt        *time.Time `json:"deleteAt"`
}
type OrderRepositoryInterface interface {
	Save(ctx context.Context, order *entities.Order) error
	List(ctx context.Context, page, limit int, sort string) ([]entities.Order, error)
}
