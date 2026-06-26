package valueobjects

import (
	"strings"

	"golang_ddd_architecture_as_poetry_clean_project/internal/shared"
)

type ProductSKU struct {
	value string
}

func NewProductSKU(raw string) (ProductSKU, error) {
	if len(raw) < 3 || len(raw) > 32 {
		return ProductSKU{}, shared.InvalidValue("sku", "must be between 3 and 32 characters")
	}
	return ProductSKU{value: strings.ToUpper(raw)}, nil
}

func (p ProductSKU) String() string {
	return p.value
}
