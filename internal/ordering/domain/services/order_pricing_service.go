package services

import (
	"golang_ddd_architecture_as_poetry_clean_project/internal/shared"
)

type PricedOrder interface {
	LineCount() int
	Subtotal() (shared.Money, error)
}

func ApplyVolumeDiscount(order PricedOrder) (shared.Money, error) {
	subtotal, err := order.Subtotal()
	if err != nil {
		return shared.Money{}, err
	}

	if order.LineCount() >= 3 {
		return shared.NewMoneyFromCents(subtotal.Cents() * 90 / 100)
	}

	return subtotal, nil
}
