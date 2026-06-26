package validators

import (
	"errors"

	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/value_objects"
	"golang_ddd_architecture_as_poetry_clean_project/internal/shared"
)

func ValidateAddProductInput(sku string, quantity uint32) shared.Notification {
	notification := shared.NewNotification()

	if err := validateSKU(sku); err != nil {
		notification.AddField("sku", fieldMessage(err))
	}

	if err := validateQuantity(quantity); err != nil {
		notification.AddField("quantity", fieldMessage(err))
	}

	return notification
}

func validateSKU(sku string) error {
	_, err := valueobjects.NewProductSKU(sku)
	return err
}

func validateQuantity(quantity uint32) error {
	_, err := valueobjects.NewQuantity(quantity)
	return err
}

func fieldMessage(err error) string {
	var domainErr shared.DomainError
	if errors.As(err, &domainErr) && domainErr.Code == shared.CodeInvalidValue {
		return domainErr.Message
	}
	return err.Error()
}
