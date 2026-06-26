package valueobjects

import "github.com/google/uuid"

type OrderID struct {
	value uuid.UUID
}

func NewOrderID() OrderID {
	return OrderID{value: uuid.New()}
}

func OrderIDFromUUID(v uuid.UUID) OrderID {
	return OrderID{value: v}
}

func (o OrderID) UUID() uuid.UUID {
	return o.value
}

func (o OrderID) String() string {
	return o.value.String()
}
