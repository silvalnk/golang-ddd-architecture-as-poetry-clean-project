package catalog

import (
	"strings"
	"sync"

	"golang_ddd_architecture_as_poetry_clean_project/internal/catalog/domain"
	"golang_ddd_architecture_as_poetry_clean_project/internal/shared"
)

type ProductInfo struct {
	SKU       string
	Name      string
	UnitPrice shared.Money
}

type CatalogPort interface {
	FindBySKU(sku string) (*ProductInfo, error)
}

type InMemoryCatalogAdapter struct {
	mu       sync.RWMutex
	products map[string]domain.CatalogProduct
}

func NewInMemoryCatalogAdapterWithSampleProducts() *InMemoryCatalogAdapter {
	products := map[string]domain.CatalogProduct{
		"RUST-BOOK":  {SKU: "RUST-BOOK", Name: "The Rust Programming Language", UnitPriceCents: 25000},
		"DDD-BOOK":   {SKU: "DDD-BOOK", Name: "Domain-Driven Design", UnitPriceCents: 18000},
		"FERRO-BOOK": {SKU: "FERRO-BOOK", Name: "Programming Rust", UnitPriceCents: 22000},
	}
	return &InMemoryCatalogAdapter{products: products}
}

func (a *InMemoryCatalogAdapter) FindBySKU(sku string) (*ProductInfo, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	product, ok := a.products[strings.ToUpper(sku)]
	if !ok {
		return nil, nil
	}

	unitPrice, err := product.UnitPrice()
	if err != nil {
		return nil, err
	}

	return &ProductInfo{
		SKU:       product.SKU,
		Name:      product.Name,
		UnitPrice: unitPrice,
	}, nil
}
