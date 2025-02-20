package order

import (
	"context"

	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/repositories"
)

type ListOrderUseCase struct {
	OrderRepository repositories.OrderRepositoryInterface
}

func NewListOrderUseCase(OrderRepository repositories.OrderRepositoryInterface) *ListOrderUseCase {
	return &ListOrderUseCase{OrderRepository: OrderRepository}
}

func (c *ListOrderUseCase) Execute(page, limit int, sort string) ([]repositories.OrderOutputDTO, error) {
	order, err := c.OrderRepository.List(context.Background(), page, limit, sort)
	if err != nil {
		return []repositories.OrderOutputDTO{}, err
	}

	var listOrder []repositories.OrderOutputDTO

	for _, o := range order {
		listOrder = append(listOrder, repositories.OrderOutputDTO{
			ID:              o.ID,
			Price:           o.Price,
			Tax:             o.Tax,
			FinalPrice:      o.FinalPrice,
			IssueDate:       o.IssueDate,
			TypeRequisition: o.TypeRequisition,
			DeleteAt:        o.DeleteAt,
		})
	}

	return listOrder, nil
}
