package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/entities"
	"github.com/stretchr/testify/assert"
)

type TestLimitList struct {
	page  int
	limit int
	sort  string
}

func (m *MockOrderRepository) List(ctx context.Context, page, limit int, sort string) ([]entities.Order, error) {
	args := m.Called()
	return args.Get(0).([]entities.Order), args.Error(1)
}

func TestListOrderUseCase_Execute(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	useCase := NewListOrderUseCase(mockRepo)

	expectedOrders := []entities.Order{
		{
			ID:         "123",
			Price:      100.0,
			Tax:        10.0,
			FinalPrice: 110.0,
			IssueDate:  time.Now(),
			DeleteAt:   nil,
		},
	}

	mockRepo.On("List").Return(expectedOrders, nil)

	limit := TestLimitList{page: 1, limit: 1, sort: "asc"}

	output, err := useCase.Execute(limit.page, limit.limit, limit.sort)

	assert.NoError(t, err)
	assert.Len(t, output, 1)
	assert.Equal(t, expectedOrders[0].ID, output[0].ID)
	assert.Equal(t, expectedOrders[0].Price, output[0].Price)
	assert.Equal(t, expectedOrders[0].Tax, output[0].Tax)
	assert.Equal(t, expectedOrders[0].FinalPrice, output[0].FinalPrice)
	assert.Empty(t, expectedOrders[0].DeleteAt)
	mockRepo.AssertExpectations(t)
}

func TestListOrderUseCase_Execute_ErrorListing(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	useCase := NewListOrderUseCase(mockRepo)

	errorMsg := errors.New("failed to list orders")
	mockRepo.On("List").Return([]entities.Order{}, errorMsg)

	limit := TestLimitList{page: 1, limit: 1, sort: "asc"}

	output, err := useCase.Execute(limit.page, limit.limit, limit.sort)

	assert.Error(t, err)
	assert.Equal(t, "failed to list orders", err.Error())
	assert.Empty(t, output)
	mockRepo.AssertExpectations(t)
}
