package main

import (
	"errors"
	"fmt"

	"golang_ddd_architecture_as_poetry_clean_project/internal/billing"
	"golang_ddd_architecture_as_poetry_clean_project/internal/catalog"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/application"
	valueobjects "golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/value_objects"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/infrastructure"
	"golang_ddd_architecture_as_poetry_clean_project/internal/shared"
)

func main() {
	fmt.Println("=== DDD in Go — Arquitetura como Poesia ===")

	catalogAdapter := catalog.NewInMemoryCatalogAdapterWithSampleProducts()
	billingHandler := billing.NewEventHandler()

	app := application.NewOrderApplicationService(
		infrastructure.NewInMemoryOrderRepository(),
		catalogAdapter,
		[]application.DomainEventHandler{billingHandler},
	)

	customer := valueobjects.NewCustomerID()
	orderID, err := app.CreateOrder(customer)
	if err != nil {
		panic(err)
	}

	err = app.AddProduct(orderID, "x", 0)
	var domainErr shared.DomainError
	if errors.As(err, &domainErr) && domainErr.Code == shared.CodeValidationFailed {
		fmt.Println("Notification de add_product:")
		fmt.Println(domainErr.Notification.String())
	}

	fmt.Printf("Order criada: %s\n", orderID.String())

	emptyOrderID, err := app.CreateOrder(valueobjects.NewCustomerID())
	if err != nil {
		panic(err)
	}

	err = app.SubmitOrder(emptyOrderID)
	if errors.As(err, &domainErr) && domainErr.Code == shared.CodeValidationFailed {
		fmt.Println("Notification de submit:")
		fmt.Println(domainErr.Notification.String())
	}

	for _, line := range []struct {
		SKU string
		Qty uint32
	}{
		{SKU: "rust-book", Qty: 1},
		{SKU: "ddd-book", Qty: 2},
		{SKU: "ferro-book", Qty: 1},
	} {
		if err := app.AddProduct(orderID, line.SKU, line.Qty); err != nil {
			panic(err)
		}
	}

	quote, err := app.QuoteFinalPrice(orderID)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Quote final: %s\n", quote)

	if err := app.SubmitOrder(orderID); err != nil {
		panic(err)
	}

	invoices := billingHandler.Invoices()
	fmt.Printf("Invoices geradas: %d\n", len(invoices))
	if len(invoices) > 0 {
		fmt.Printf("Invoice status: %s, total: %s\n", invoices[0].Status, invoices[0].Amount)
	}

	fmt.Println("\n--- Patterns demonstrated ---")
	for _, pattern := range []string{
		"Aggregate Root (Order)",
		"Value Objects (Money, Quantity, ProductSku...)",
		"Domain Events (OrderSubmitted -> Billing)",
		"Domain Service (OrderPricingService)",
		"Specification (CanSubmitOrder, IsHighValueOrder)",
		"Notification (aggregated validation errors)",
		"Factory (OrderFactory)",
		"Repository (interface + InMemoryOrderRepository)",
		"Application Service (OrderApplicationService)",
		"Anti-Corruption Layer (CatalogPort)",
		"Context Mapping (Shared Kernel, ACL, Published Language)",
		"SOLID (DIP via interfaces, SRP per layer, OCP in Specifications)",
	} {
		fmt.Printf("  • %s\n", pattern)
	}
}
