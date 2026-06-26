package specifications_test

import (
	"testing"

	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/factories"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/specifications"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/value_objects"
	"golang_ddd_architecture_as_poetry_clean_project/internal/shared"
)

func TestCanSubmitRequiresAtLeastOneLine(t *testing.T) {
	order := factories.CreateDraftOrder(valueobjects.NewCustomerID())

	ok, err := specifications.CanSubmitOrder{}.IsSatisfiedBy(order)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected can submit to be false for empty draft order")
	}
}

func TestIsHighValueOrder(t *testing.T) {
	order := factories.CreateDraftOrder(valueobjects.NewCustomerID())
	sku, _ := valueobjects.NewProductSKU("RUST-BOOK")
	qty, _ := valueobjects.NewQuantity(2)
	price, _ := shared.NewMoneyBRL(300, 0)

	if err := order.AddLine(sku, qty, price); err != nil {
		t.Fatalf("add line: %v", err)
	}

	spec, err := specifications.NewIsHighValueOrderAboveBRL500()
	if err != nil {
		t.Fatalf("build spec: %v", err)
	}

	ok, err := spec.IsSatisfiedBy(order)
	if err != nil {
		t.Fatalf("evaluate spec: %v", err)
	}
	if !ok {
		t.Fatal("expected order to be high value")
	}
}

func TestCanComposeSpecificationWithAnd(t *testing.T) {
	order := factories.CreateDraftOrder(valueobjects.NewCustomerID())
	sku, _ := valueobjects.NewProductSKU("RUST-BOOK")
	qty, _ := valueobjects.NewQuantity(2)
	price, _ := shared.NewMoneyBRL(300, 0)
	if err := order.AddLine(sku, qty, price); err != nil {
		t.Fatalf("add line: %v", err)
	}

	highValue, err := specifications.NewIsHighValueOrderAboveBRL500()
	if err != nil {
		t.Fatalf("create spec: %v", err)
	}

	spec := specifications.And(specifications.CanSubmitOrder{}, highValue)
	ok, err := spec.IsSatisfiedBy(order)
	if err != nil {
		t.Fatalf("evaluate composed spec: %v", err)
	}
	if !ok {
		t.Fatal("expected composed spec to be true")
	}
}

func TestNotifyingSpecAccumulatesFailures(t *testing.T) {
	order := factories.CreateDraftOrder(valueobjects.NewCustomerID())
	notification := shared.NewNotification()

	specifications.ValidateSpecification(specifications.CanSubmitOrder{}, order, &notification)
	highValue, err := specifications.NewIsHighValueOrderAboveBRL500()
	if err != nil {
		t.Fatalf("create spec: %v", err)
	}
	specifications.ValidateSpecification(highValue, order, &notification)

	if len(notification.Errors()) != 2 {
		t.Fatalf("expected two failures, got %d", len(notification.Errors()))
	}
}
