package valueobjects

import "github.com/google/uuid"

type CustomerID struct {
	value uuid.UUID
}

func NewCustomerID() CustomerID {
	return CustomerID{value: uuid.New()}
}

func CustomerIDFromUUID(v uuid.UUID) CustomerID {
	return CustomerID{value: v}
}

func (c CustomerID) UUID() uuid.UUID {
	return c.value
}

func (c CustomerID) String() string {
	return c.value.String()
}
