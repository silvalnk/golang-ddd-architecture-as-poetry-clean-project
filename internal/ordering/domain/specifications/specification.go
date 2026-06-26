package specifications

import (
	"golang_ddd_architecture_as_poetry_clean_project/internal/ordering/domain/aggregate"
	"golang_ddd_architecture_as_poetry_clean_project/internal/shared"
)

type Specification interface {
	IsSatisfiedBy(order *aggregate.Order) (bool, error)
}

type specificationFunc func(order *aggregate.Order) (bool, error)

func (f specificationFunc) IsSatisfiedBy(order *aggregate.Order) (bool, error) {
	return f(order)
}

type NotifyingSpecification interface {
	Specification
	FailureMessage() string
}

func And(left, right Specification) Specification {
	return specificationFunc(func(order *aggregate.Order) (bool, error) {
		leftResult, err := left.IsSatisfiedBy(order)
		if err != nil {
			return false, err
		}
		rightResult, err := right.IsSatisfiedBy(order)
		if err != nil {
			return false, err
		}
		return leftResult && rightResult, nil
	})
}

func Not(spec Specification) Specification {
	return specificationFunc(func(order *aggregate.Order) (bool, error) {
		result, err := spec.IsSatisfiedBy(order)
		if err != nil {
			return false, err
		}
		return !result, nil
	})
}

func ValidateSpecification(spec NotifyingSpecification, order *aggregate.Order, notification *shared.Notification) {
	ok, err := spec.IsSatisfiedBy(order)
	if err != nil {
		notification.Add(err.Error())
		return
	}
	if !ok {
		notification.Add(spec.FailureMessage())
	}
}

type CanSubmitOrder struct{}

func (CanSubmitOrder) IsSatisfiedBy(order *aggregate.Order) (bool, error) {
	return order.Status().CanSubmit() && len(order.Lines()) > 0, nil
}

func (CanSubmitOrder) FailureMessage() string {
	return "order must be in draft state with at least one line"
}

type IsHighValueOrder struct {
	threshold shared.Money
}

func NewIsHighValueOrderAboveBRL500() (IsHighValueOrder, error) {
	threshold, err := shared.NewMoneyBRL(500, 0)
	if err != nil {
		return IsHighValueOrder{}, err
	}
	return IsHighValueOrder{threshold: threshold}, nil
}

func (s IsHighValueOrder) IsSatisfiedBy(order *aggregate.Order) (bool, error) {
	subtotal, err := order.Subtotal()
	if err != nil {
		return false, err
	}
	return subtotal.IsGreaterThan(s.threshold)
}

func (IsHighValueOrder) FailureMessage() string {
	return "order subtotal must exceed BRL 500"
}
