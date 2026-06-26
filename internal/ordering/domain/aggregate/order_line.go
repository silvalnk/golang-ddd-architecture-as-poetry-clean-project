package aggregate

import (
	valueobjects "golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/value_objects"
	"golang_ddd_architecture_as_poetry_clean_project/internal/shared"
)

type OrderLine struct {
	sku       valueobjects.ProductSKU
	quantity  valueobjects.Quantity
	unitPrice shared.Money
}

func NewOrderLine(sku valueobjects.ProductSKU, quantity valueobjects.Quantity, unitPrice shared.Money) OrderLine {
	return OrderLine{
		sku:       sku,
		quantity:  quantity,
		unitPrice: unitPrice,
	}
}

func (l OrderLine) SKU() valueobjects.ProductSKU {
	return l.sku
}

func (l OrderLine) Quantity() valueobjects.Quantity {
	return l.quantity
}

func (l OrderLine) UnitPrice() shared.Money {
	return l.unitPrice
}

func (l OrderLine) LineTotal() shared.Money {
	return l.unitPrice.Mul(l.quantity.Value())
}
