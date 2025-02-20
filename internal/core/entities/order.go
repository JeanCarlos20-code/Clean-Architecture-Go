package entities

import (
	"time"
)

type Order struct {
	ID              string
	Price           float64
	Tax             float64
	FinalPrice      float64
	IssueDate       time.Time
	TypeRequisition string
	DeleteAt        *time.Time
}

func NewOrder(data Order) *Order {
	return &Order{
		ID:              data.ID,
		Price:           data.Price,
		Tax:             data.Tax,
		FinalPrice:      data.FinalPrice,
		IssueDate:       data.IssueDate,
		TypeRequisition: data.TypeRequisition,
		DeleteAt:        data.DeleteAt,
	}

}

func (o *Order) CalculateFinalPrice() float64 {
	o.FinalPrice = o.Price + o.Tax
	return o.FinalPrice
}

func (o *Order) StringToTime(data string) time.Time {
	issueDate, err := time.Parse(time.RFC3339, data)
	if err != nil {
		panic(err)
	}
	return issueDate
}
