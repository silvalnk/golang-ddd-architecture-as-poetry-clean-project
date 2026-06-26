package shared

import "fmt"

const (
	CodeInvariantViolation = "invariant_violation"
	CodeNotFound           = "not_found"
	CodeInvalidOperation   = "invalid_operation"
	CodeInvalidValue       = "invalid_value"
	CodeValidationFailed   = "validation_failed"
)

type DomainError struct {
	Code         string
	Message      string
	Entity       string
	ID           string
	Field        string
	Notification *Notification
}

func (e DomainError) Error() string {
	switch e.Code {
	case CodeNotFound:
		return fmt.Sprintf("entity not found: %s (%s)", e.Entity, e.ID)
	case CodeInvalidValue:
		return fmt.Sprintf("invalid value for %s: %s", e.Field, e.Message)
	case CodeValidationFailed:
		if e.Notification != nil {
			return fmt.Sprintf("validation failed:\n%s", e.Notification.String())
		}
	}
	return e.Message
}

func InvariantViolation(message string) error {
	return DomainError{Code: CodeInvariantViolation, Message: message}
}

func NotFound(entity, id string) error {
	return DomainError{Code: CodeNotFound, Entity: entity, ID: id}
}

func InvalidOperation(message string) error {
	return DomainError{Code: CodeInvalidOperation, Message: message}
}

func InvalidValue(field, reason string) error {
	return DomainError{Code: CodeInvalidValue, Field: field, Message: reason}
}

func ValidationFailed(notification Notification) error {
	copy := notification
	return DomainError{Code: CodeValidationFailed, Notification: &copy}
}
