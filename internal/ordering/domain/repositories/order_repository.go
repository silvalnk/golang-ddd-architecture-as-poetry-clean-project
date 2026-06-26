package repositories

import (
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/aggregate"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/value_objects"
)

type OrderRepository interface {
	Save(order *aggregate.Order) error
	FindByID(id valueobjects.OrderID) (*aggregate.Order, error)
}
