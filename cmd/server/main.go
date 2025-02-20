package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/JeanCarlos20-code/CleanArchitecture/configs"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/controller/graph/graph"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/controller/grpc/pb"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/controller/grpc/service"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/controller/web/webserver"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/event/handler"
	"github.com/JeanCarlos20-code/CleanArchitecture/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.GetDBDriver(), fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.GetDBUser(), configs.GetDBPassword(), configs.GetDBHost(), configs.GetDBPort(), configs.GetDBName()))
	if err != nil {
		panic(err)
	}

	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrderUseCase(db)

	webserver := webserver.NewWebServer(configs.GetWebServerPort())
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webserver.AddHandler("/order", webOrderHandler.List)
	fmt.Println("Starting web server on port", configs.GetWebServerPort())
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GetGRPCServerPort())
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GetGRPCServerPort()))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GetGraphQLServerPort())
	http.ListenAndServe(":"+configs.GetGraphQLServerPort(), nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
