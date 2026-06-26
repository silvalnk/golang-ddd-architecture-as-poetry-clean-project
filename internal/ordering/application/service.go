package application

import (
	"fmt"

	"golang_ddd_architecture_as_poetry_clean_project/internal/catalog"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/aggregate"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/events"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/factories"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/repositories"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/services"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/specifications"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/validators"
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/value_objects"
	"golang_ddd_architecture_as_poetry_clean_project/internal/shared"
)

type DomainEventHandler interface {
	Handle(event events.DomainEvent)
}

type OrderApplicationService struct {
	repository    repositories.OrderRepository
	catalog       catalog.CatalogPort
	eventHandlers []DomainEventHandler
}

func NewOrderApplicationService(
	repository repositories.OrderRepository,
	catalogPort catalog.CatalogPort,
	eventHandlers []DomainEventHandler,
) *OrderApplicationService {
	return &OrderApplicationService{
		repository:    repository,
		catalog:       catalogPort,
		eventHandlers: eventHandlers,
	}
}

func (s *OrderApplicationService) CreateOrder(customerID valueobjects.CustomerID) (valueobjects.OrderID, error) {
	order := factories.CreateDraftOrder(customerID)
	id := order.ID()
	if err := s.persistAndDispatch(order); err != nil {
		return valueobjects.OrderID{}, err
	}
	return id, nil
}

func (s *OrderApplicationService) AddProduct(orderID valueobjects.OrderID, sku string, quantity uint32) error {
	inputValidation := validators.ValidateAddProductInput(sku, quantity)
	if inputValidation.HasErrors() {
		return inputValidation.IntoError()
	}

	product, err := s.catalog.FindBySKU(sku)
	if err != nil {
		return err
	}
	if product == nil {
		return shared.NotFound("product", sku)
	}

	order, err := s.load(orderID)
	if err != nil {
		return err
	}

	productSKU, err := valueobjects.NewProductSKU(product.SKU)
	if err != nil {
		return err
	}
	qty, err := valueobjects.NewQuantity(quantity)
	if err != nil {
		return err
	}

	if err := order.AddLine(productSKU, qty, product.UnitPrice); err != nil {
		return err
	}
	return s.persistAndDispatch(order)
}

func (s *OrderApplicationService) SubmitOrder(orderID valueobjects.OrderID) error {
	order, err := s.load(orderID)
	if err != nil {
		return err
	}

	validation := validators.ValidateSubmitOrder(order)
	if validation.HasErrors() {
		return validation.IntoError()
	}

	if err := order.Submit(); err != nil {
		return err
	}
	return s.persistAndDispatch(order)
}

func (s *OrderApplicationService) ValidateSubmit(orderID valueobjects.OrderID) (shared.Notification, error) {
	order, err := s.load(orderID)
	if err != nil {
		return shared.Notification{}, err
	}
	return validators.ValidateSubmitOrder(order), nil
}

func (s *OrderApplicationService) QuoteFinalPrice(orderID valueobjects.OrderID) (string, error) {
	order, err := s.load(orderID)
	if err != nil {
		return "", err
	}

	total, err := services.ApplyVolumeDiscount(order)
	if err != nil {
		return "", err
	}

	highValueSpec, err := specifications.NewIsHighValueOrderAboveBRL500()
	if err != nil {
		return "", err
	}
	highValue, err := highValueSpec.IsSatisfiedBy(order)
	if err != nil {
		return "", err
	}

	if highValue {
		return fmt.Sprintf("%s (high-value order)", total), nil
	}
	return total.String(), nil
}

func (s *OrderApplicationService) load(id valueobjects.OrderID) (*aggregate.Order, error) {
	order, err := s.repository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, shared.NotFound("order", id.String())
	}
	return order, nil
}

func (s *OrderApplicationService) persistAndDispatch(order *aggregate.Order) error {
	if err := s.repository.Save(order); err != nil {
		return err
	}
	eventsToDispatch := order.TakePendingEvents()
	for _, event := range eventsToDispatch {
		for _, handler := range s.eventHandlers {
			handler.Handle(event)
		}
	}
	return nil
}
