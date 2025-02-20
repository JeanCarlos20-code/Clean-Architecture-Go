package service

import (
	"context"
	"time"

	"github.com/JeanCarlos20-code/CleanArchitecture/internal/controller/grpc/pb"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/repositories"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/use-cases/order"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	CreateOrderUseCase order.CreateOrderUseCase
	ListOrderUseCase   order.ListOrderUseCase
}

func NewOrderService(createOrderUseCase order.CreateOrderUseCase, listOrderUseCase order.ListOrderUseCase) *OrderService {
	return &OrderService{
		CreateOrderUseCase: createOrderUseCase,
		ListOrderUseCase:   listOrderUseCase,
	}
}

func ConvertProtoTimestampToTime(protoTs *timestamppb.Timestamp) time.Time {
	if protoTs == nil {
		return time.Time{}
	}
	return protoTs.AsTime()
}

func ConvertTimeToProtoTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderInput) (*pb.OrderOutput, error) {
	issueDate, err := time.Parse(time.RFC3339, in.IssueDate)
	if err != nil {
		return nil, err
	}

	dto := repositories.OrderInputDTO{
		Price:     float64(in.Price),
		Tax:       float64(in.Tax),
		IssueDate: issueDate,
	}

	dto.TypeRequisition = "gRPC"
	output, err := s.CreateOrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.OrderOutput{
		Id:              output.ID,
		Price:           float32(output.Price),
		Tax:             float32(output.Tax),
		FinalPrice:      float32(output.FinalPrice),
		IssueDate:       ConvertTimeToProtoTimestamp(&output.IssueDate),
		TypeRequisition: output.TypeRequisition,
		DeleteAt:        ConvertTimeToProtoTimestamp(output.DeleteAt),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.ListOrdersInput) (*pb.ListOrdersOutput, error) {
	output, err := s.ListOrderUseCase.Execute(int(in.Page), int(in.Limit), string(in.Sort))
	if err != nil {
		return nil, err
	}
	var pbOrders []*pb.OrderOutput
	for _, o := range output {
		pbOrders = append(pbOrders, &pb.OrderOutput{
			Id:              o.ID,
			Price:           float32(o.Price),
			Tax:             float32(o.Tax),
			FinalPrice:      float32(o.FinalPrice),
			IssueDate:       ConvertTimeToProtoTimestamp(&o.IssueDate),
			TypeRequisition: o.TypeRequisition,
			DeleteAt:        ConvertTimeToProtoTimestamp(o.DeleteAt),
		})
	}
	return &pb.ListOrdersOutput{
		Orders:     pbOrders,
		TotalCount: int32(len(pbOrders)),
	}, nil
}
