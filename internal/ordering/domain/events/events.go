package events

import (
	"time"

	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/value_objects"
	"golang_ddd_architecture_as_poetry_clean_project/internal/shared"
)

type DomainEvent interface {
	EventName() string
}

type OrderCreated struct {
	OrderID    valueobjects.OrderID
	CustomerID valueobjects.CustomerID
	OccurredAt time.Time
}

func (OrderCreated) EventName() string { return "OrderCreated" }

type OrderLineAdded struct {
	OrderID    valueobjects.OrderID
	SKU        valueobjects.ProductSKU
	Quantity   valueobjects.Quantity
	UnitPrice  shared.Money
	OccurredAt time.Time
}

func (OrderLineAdded) EventName() string { return "OrderLineAdded" }

type OrderSubmitted struct {
	OrderID    valueobjects.OrderID
	CustomerID valueobjects.CustomerID
	Total      shared.Money
	OccurredAt time.Time
}

func (OrderSubmitted) EventName() string { return "OrderSubmitted" }

type OrderCancelled struct {
	OrderID    valueobjects.OrderID
	Reason     string
	OccurredAt time.Time
}

func (OrderCancelled) EventName() string { return "OrderCancelled" }
