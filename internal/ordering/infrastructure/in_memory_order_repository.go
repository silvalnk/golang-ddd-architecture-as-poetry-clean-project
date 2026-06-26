package infrastructure

import (
	"sync"

	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/aggregate"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/value_objects"
)

type InMemoryOrderRepository struct {
	mu    sync.RWMutex
	store map[string]aggregate.Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		store: make(map[string]aggregate.Order),
	}
}

func (r *InMemoryOrderRepository) Save(order *aggregate.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[order.ID().String()] = *order.Clone()
	return nil
}

func (r *InMemoryOrderRepository) FindByID(id valueobjects.OrderID) (*aggregate.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.store[id.String()]
	if !ok {
		return nil, nil
	}
	copy := order.Clone()
	return copy, nil
}
