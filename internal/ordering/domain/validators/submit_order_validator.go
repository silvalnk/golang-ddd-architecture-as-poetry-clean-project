package validators

import (
	"fmt"

	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/aggregate"
	"golang_ddd_architecture_as_poetry_clean_project/internal/shared"
)

func ValidateSubmitOrder(order *aggregate.Order) shared.Notification {
	notification := shared.NewNotification()

	if !order.Status().CanSubmit() {
		notification.AddField("status", fmt.Sprintf("must be 'draft', got '%s'", order.Status()))
	}

	if len(order.Lines()) == 0 {
		notification.AddField("lines", "order must have at least one line")
	}

	subtotal, err := order.Subtotal()
	if err != nil {
		notification.Add(err.Error())
		return notification
	}
	if subtotal.Cents() == 0 {
		notification.AddField("subtotal", "subtotal must be greater than zero")
	}

	return notification
}
