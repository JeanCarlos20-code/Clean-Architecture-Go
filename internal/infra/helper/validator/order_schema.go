package validator

import (
	"time"

	"github.com/JeanCarlos20-code/CleanArchitecture/internal/core/entities"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func ValidateOrder(order *entities.Order) error {
	return validation.ValidateStruct(order,
		validation.Field(&order.Price, validation.Required, validation.By(isFloat)),
		validation.Field(&order.Tax, validation.Required, validation.By(isFloat)),
		validation.Field(&order.IssueDate, validation.Required, validation.By(isTime)),
	)
}

func isFloat(value interface{}) error {
	_, ok := value.(float64)
	if !ok {
		return validation.NewError("validation_float", "O valor deve ser um número decimal")
	}
	return nil
}

func isTime(value interface{}) error {
	_, ok := value.(time.Time)
	if !ok {
		return validation.NewError("validation_time", "O valor deve ser um tipo de data válido")
	}
	return nil
}
