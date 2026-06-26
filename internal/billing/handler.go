package billing

import (
	"sync"

	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/events"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/value_objects"
	"golang_ddd_architecture_as_poetry_clean_project/internal/shared"
)

type InvoiceStatus string

const (
	InvoiceStatusPending   InvoiceStatus = "pending"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

type Invoice struct {
	OrderID valueobjects.OrderID
	Amount  shared.Money
	Status  InvoiceStatus
}

type EventHandler struct {
	mu       sync.Mutex
	invoices []Invoice
}

func NewEventHandler() *EventHandler {
	return &EventHandler{
		invoices: make([]Invoice, 0),
	}
}

func (h *EventHandler) Invoices() []Invoice {
	h.mu.Lock()
	defer h.mu.Unlock()

	out := make([]Invoice, len(h.invoices))
	copy(out, h.invoices)
	return out
}

func (h *EventHandler) Handle(event events.DomainEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	switch e := event.(type) {
	case events.OrderSubmitted:
		h.invoices = append(h.invoices, Invoice{
			OrderID: e.OrderID,
			Amount:  e.Total,
			Status:  InvoiceStatusPending,
		})
	case events.OrderCancelled:
		for i := range h.invoices {
			if h.invoices[i].OrderID.String() == e.OrderID.String() {
				h.invoices[i].Status = InvoiceStatusCancelled
				return
			}
		}
	}
}
