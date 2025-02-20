package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/entities"
	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/repositories"
	"github.com/JeanCarlos20-code/CleanArchitecture/pkg/events"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEventInterface struct {
	mock.Mock
}

func (m *MockEventInterface) SetPayload(payload interface{}) {
	m.Called(payload)
}

func (m *MockEventInterface) GetDateTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockEventInterface) GetPayload() interface{} {
	return m.Called()
}

func (m *MockEventInterface) GetName() string {
	args := m.Called()
	return args.Get(0).(string)
}

type MockEventDispatcherInterface struct {
	mock.Mock
}

func (m *MockEventDispatcherInterface) Dispatch(event events.EventInterface) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventDispatcherInterface) Register(eventName string, handler events.EventHandlerInterface) error {
	args := m.Called(eventName, handler)
	return args.Error(0)
}

func (m *MockEventDispatcherInterface) Remove(eventName string, handler events.EventHandlerInterface) error {
	args := m.Called(eventName, handler)
	return args.Error(0)
}

func (m *MockEventDispatcherInterface) Has(eventName string, handler events.EventHandlerInterface) bool {
	args := m.Called(eventName, handler)
	return args.Bool(0)
}

func (m *MockEventDispatcherInterface) Clear() {
	m.Called()
}

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Save(ctx context.Context, order *entities.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

var mockRepo = new(MockOrderRepository)
var mockEvent = new(MockEventInterface)
var mockDispatcher = new(MockEventDispatcherInterface)
var useCase = NewCreateOrderUseCase(mockRepo, mockEvent, mockDispatcher)

func TestCreateOrderUseCase_Execute(t *testing.T) {
	issueDate, err := time.Parse(time.RFC3339, "2025-02-19T15:04:05Z")
	assert.NoError(t, err)

	input := repositories.OrderInputDTO{
		Price:     100.0,
		Tax:       10.0,
		IssueDate: issueDate,
	}

	order := entities.Order{
		Price:     input.Price,
		Tax:       input.Tax,
		IssueDate: input.IssueDate,
	}

	order.CalculateFinalPrice()

	expectedOrder := order

	mockRepo.On("Save", &expectedOrder).Return(nil).Once()
	mockEvent.On("SetPayload", mock.Anything).Return().Once()
	mockDispatcher.On("Dispatch", mock.Anything).Return(nil).Once()

	output, err := useCase.Execute(input)

	assert.NoError(t, err)
	assert.Equal(t, order.CalculateFinalPrice(), 110.0)
	assert.Equal(t, input.Price, output.Price)
	assert.Equal(t, input.Tax, output.Tax)
	assert.Equal(t, order.CalculateFinalPrice(), output.FinalPrice)
	assert.Empty(t, output.DeleteAt)

	mockEvent.AssertExpectations(t)
	mockDispatcher.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestCreateOrderUseCase_Execute_ErrorSaving(t *testing.T) {
	issueDate, err := time.Parse(time.RFC3339, "2025-02-19T15:04:05Z")
	assert.NoError(t, err)

	input := repositories.OrderInputDTO{
		Price:     100.0,
		Tax:       10.0,
		IssueDate: issueDate,
	}

	expectedOrder := entities.Order{
		Price:      input.Price,
		Tax:        input.Tax,
		IssueDate:  input.IssueDate,
		FinalPrice: input.Price + input.Tax,
	}

	errorMsg := errors.New("failed to save order")
	mockRepo.On("Save", &expectedOrder).Return(errorMsg)

	output, err := useCase.Execute(input)

	assert.Error(t, err)
	assert.Equal(t, "failed to save order", err.Error())
	assert.Empty(t, output)

	mockEvent.AssertExpectations(t)
	mockDispatcher.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
