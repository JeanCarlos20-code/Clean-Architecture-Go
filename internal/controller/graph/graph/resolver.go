package graph

import "github.com/JeanCarlos20-code/CleanArchitecture/internal/core/use-cases/order"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	CreateOrderUseCase order.CreateOrderUseCase
	ListOrderUseCase   order.ListOrderUseCase
}
