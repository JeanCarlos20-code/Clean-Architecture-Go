package order

import (
	"context"

	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/entities"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/repositories"
	"github.com/JeanCarlos20-code/CleanArchitecture/pkg/events"
)

type CreateOrderUseCase struct {
	OrderRepository repositories.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewCreateOrderUseCase(
	OrderRepository repositories.OrderRepositoryInterface,
	OrderCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: OrderRepository,
		OrderCreated:    OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (c *CreateOrderUseCase) Execute(input repositories.OrderInputDTO) (repositories.OrderOutputDTO, error) {
	order := entities.Order{
		Price:           input.Price,
		Tax:             input.Tax,
		IssueDate:       input.IssueDate,
		TypeRequisition: input.TypeRequisition,
	}
	order.CalculateFinalPrice()

	if err := c.OrderRepository.Save(context.Background(), &order); err != nil {
		return repositories.OrderOutputDTO{}, err
	}
	dto := repositories.OrderOutputDTO{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
		IssueDate:  order.IssueDate,
		DeleteAt:   order.DeleteAt,
	}

	c.OrderCreated.SetPayload(dto)
	c.EventDispatcher.Dispatch(c.OrderCreated)

	return dto, nil
}
