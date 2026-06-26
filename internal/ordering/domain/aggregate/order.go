package aggregate

import (
	"fmt"
	"time"

	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/events"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/services"
	valueobjects "golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/value_objects"
	"golang_ddd_architecture_as_poetry_clean_project/internal/shared"
)

type Order struct {
	id            valueobjects.OrderID
	customerID    valueobjects.CustomerID
	status        valueobjects.OrderStatus
	lines         []OrderLine
	createdAt     time.Time
	pendingEvents []events.DomainEvent
}

func NewDraftOrder(id valueobjects.OrderID, customerID valueobjects.CustomerID) *Order {
	createdAt := time.Now().UTC()
	order := &Order{
		id:            id,
		customerID:    customerID,
		status:        valueobjects.OrderStatusDraft,
		lines:         make([]OrderLine, 0),
		createdAt:     createdAt,
		pendingEvents: make([]events.DomainEvent, 0),
	}
	order.record(events.OrderCreated{
		OrderID:    id,
		CustomerID: customerID,
		OccurredAt: createdAt,
	})
	return order
}

func (o *Order) ID() valueobjects.OrderID {
	return o.id
}

func (o *Order) CustomerID() valueobjects.CustomerID {
	return o.customerID
}

func (o *Order) Status() valueobjects.OrderStatus {
	return o.status
}

func (o *Order) Lines() []OrderLine {
	out := make([]OrderLine, len(o.lines))
	copy(out, o.lines)
	return out
}

func (o *Order) LineCount() int {
	return len(o.lines)
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

func (o *Order) Subtotal() (shared.Money, error) {
	total, err := shared.NewMoneyFromCents(0)
	if err != nil {
		return shared.Money{}, err
	}
	for _, line := range o.lines {
		total, err = total.Add(line.LineTotal())
		if err != nil {
			return shared.Money{}, err
		}
	}
	return total, nil
}

func (o *Order) AddLine(sku valueobjects.ProductSKU, quantity valueobjects.Quantity, unitPrice shared.Money) error {
	if !o.status.CanAddLine() {
		return shared.InvalidOperation(fmt.Sprintf("cannot add lines to order in '%s' state", o.status))
	}

	for _, line := range o.lines {
		if line.SKU().String() == sku.String() {
			return shared.InvariantViolation("SKU already present in order")
		}
	}

	line := NewOrderLine(sku, quantity, unitPrice)
	o.lines = append(o.lines, line)

	o.record(events.OrderLineAdded{
		OrderID:    o.id,
		SKU:        sku,
		Quantity:   quantity,
		UnitPrice:  unitPrice,
		OccurredAt: time.Now().UTC(),
	})

	return nil
}

func (o *Order) Submit() error {
	if !o.status.CanSubmit() {
		return shared.InvalidOperation(fmt.Sprintf("order cannot be submitted in '%s' state", o.status))
	}
	if len(o.lines) == 0 {
		return shared.InvariantViolation("order must have at least one line")
	}

	o.status = valueobjects.OrderStatusSubmitted
	total, err := services.ApplyVolumeDiscount(o)
	if err != nil {
		return err
	}

	o.record(events.OrderSubmitted{
		OrderID:    o.id,
		CustomerID: o.customerID,
		Total:      total,
		OccurredAt: time.Now().UTC(),
	})
	return nil
}

func (o *Order) Cancel(reason string) error {
	if !o.status.CanCancel() {
		return shared.InvalidOperation(fmt.Sprintf("order cannot be cancelled in '%s' state", o.status))
	}

	o.status = valueobjects.OrderStatusCancelled
	o.record(events.OrderCancelled{
		OrderID:    o.id,
		Reason:     reason,
		OccurredAt: time.Now().UTC(),
	})
	return nil
}

func (o *Order) TakePendingEvents() []events.DomainEvent {
	out := make([]events.DomainEvent, len(o.pendingEvents))
	copy(out, o.pendingEvents)
	o.pendingEvents = make([]events.DomainEvent, 0)
	return out
}

func (o *Order) Clone() *Order {
	lines := make([]OrderLine, len(o.lines))
	copy(lines, o.lines)
	pending := make([]events.DomainEvent, len(o.pendingEvents))
	copy(pending, o.pendingEvents)

	return &Order{
		id:            o.id,
		customerID:    o.customerID,
		status:        o.status,
		lines:         lines,
		createdAt:     o.createdAt,
		pendingEvents: pending,
	}
}

func (o *Order) record(event events.DomainEvent) {
	o.pendingEvents = append(o.pendingEvents, event)
}
