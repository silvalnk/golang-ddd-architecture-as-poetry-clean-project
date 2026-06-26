package application_test

import (
	"errors"
	"testing"

	"golang_ddd_architecture_as_poetry_clean_project/internal/billing"
	"golang_ddd_architecture_as_poetry_clean_project/internal/catalog"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/application"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/value_objects"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/infrastructure"
	"golang_ddd_architecture_as_poetry_clean_project/internal/shared"
)

func TestAddProductAggregatesValidationErrors(t *testing.T) {
	app := application.NewOrderApplicationService(
		infrastructure.NewInMemoryOrderRepository(),
		catalog.NewInMemoryCatalogAdapterWithSampleProducts(),
		nil,
	)

	orderID, err := app.CreateOrder(valueobjects.NewCustomerID())
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	err = app.AddProduct(orderID, "x", 0)
	if err == nil {
		t.Fatal("expected validation error")
	}

	var domainErr shared.DomainError
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected domain error, got %T", err)
	}
	if domainErr.Code != shared.CodeValidationFailed {
		t.Fatalf("expected validation_failed, got %s", domainErr.Code)
	}
	if len(domainErr.Notification.Errors()) != 2 {
		t.Fatalf("expected 2 validation messages, got %d", len(domainErr.Notification.Errors()))
	}
}

func TestSubmitOrderPublishesInvoice(t *testing.T) {
	billingHandler := billing.NewEventHandler()
	app := application.NewOrderApplicationService(
		infrastructure.NewInMemoryOrderRepository(),
		catalog.NewInMemoryCatalogAdapterWithSampleProducts(),
		[]application.DomainEventHandler{billingHandler},
	)

	orderID, err := app.CreateOrder(valueobjects.NewCustomerID())
	if err != nil {
		t.Fatalf("create order: %v", err)
	}

	if err := app.AddProduct(orderID, "rust-book", 1); err != nil {
		t.Fatalf("add product: %v", err)
	}
	if err := app.SubmitOrder(orderID); err != nil {
		t.Fatalf("submit order: %v", err)
	}

	invoices := billingHandler.Invoices()
	if len(invoices) != 1 {
		t.Fatalf("expected one invoice, got %d", len(invoices))
	}
	if invoices[0].Status != billing.InvoiceStatusPending {
		t.Fatalf("expected pending invoice, got %s", invoices[0].Status)
	}
}
