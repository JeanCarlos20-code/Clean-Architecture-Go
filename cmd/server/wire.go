//go:build wireinject
// +build wireinject

package main

import (
	"database/sql"

	"github.com/JeanCarlos20-code/CleanArchitecture/internal/controller/web"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/repositories"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/use-cases/order"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/event"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/infra/database"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/infra/database/SQLC"
	"github.com/JeanCarlos20-code/CleanArchitecture/pkg/events"
	"github.com/google/wire"
)

var setOrderRepositoryDependency = wire.NewSet(
	database.NewOrderRepository,
	wire.Bind(new(repositories.OrderRepositoryInterface), new(*database.OrderRepository)),
)

var setEventDispatcherDependency = wire.NewSet(
	events.NewEventDispatcher,
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
	wire.Bind(new(events.EventDispatcherInterface), new(*events.EventDispatcher)),
)

var setOrderCreatedEvent = wire.NewSet(
	event.NewOrderCreated,
	wire.Bind(new(events.EventInterface), new(*event.OrderCreated)),
)

func ProvideSQLCQueries(db *sql.DB) *SQLC.Queries {
	return SQLC.New(db)
}

func NewCreateOrderUseCase(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *order.CreateOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		ProvideSQLCQueries,
		order.NewCreateOrderUseCase,
	)
	return &order.CreateOrderUseCase{}
}

func NewListOrderUseCase(db *sql.DB) *order.ListOrderUseCase {
	wire.Build(
		setOrderRepositoryDependency,
		ProvideSQLCQueries,
		order.NewListOrderUseCase,
	)
	return &order.ListOrderUseCase{}
}

func NewWebOrderHandler(db *sql.DB, eventDispatcher events.EventDispatcherInterface) *web.WebOrderHandler {
	wire.Build(
		ProvideSQLCQueries,
		setOrderRepositoryDependency,
		setOrderCreatedEvent,
		web.NewWebOrderHandler,
	)
	return &web.WebOrderHandler{}
}
