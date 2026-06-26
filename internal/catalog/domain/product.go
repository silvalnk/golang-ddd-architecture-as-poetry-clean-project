package domain

import "golang_ddd_architecture_as_poetry_clean_project/internal/shared"

type CatalogProduct struct {
	SKU            string
	Name           string
	UnitPriceCents int64
}

func (p CatalogProduct) UnitPrice() (shared.Money, error) {
	return shared.NewMoneyFromCents(p.UnitPriceCents)
}
