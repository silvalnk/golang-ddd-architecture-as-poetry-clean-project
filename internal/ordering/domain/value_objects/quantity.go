package valueobjects

import "golang_ddd_architecture_as_poetry_clean_project/internal/shared"

const MaxQuantity uint32 = 999

type Quantity struct {
	value uint32
}

func NewQuantity(value uint32) (Quantity, error) {
	if value == 0 {
		return Quantity{}, shared.InvalidValue("quantity", "must be greater than zero")
	}
	if value > MaxQuantity {
		return Quantity{}, shared.InvalidValue("quantity", "maximum allowed is 999")
	}
	return Quantity{value: value}, nil
}

func (q Quantity) Value() uint32 {
	return q.value
}
