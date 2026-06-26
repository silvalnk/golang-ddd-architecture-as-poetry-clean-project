package factories

import (
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/aggregate"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/value_objects"
)

func CreateDraftOrder(customerID valueobjects.CustomerID) *aggregate.Order {
	return aggregate.NewDraftOrder(valueobjects.NewOrderID(), customerID)
}

func Reconstitute(order *aggregate.Order) (*aggregate.Order, error) {
	return order, nil
}
