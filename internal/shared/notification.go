package shared

import "strings"

type ValidationMessage struct {
	Field   string
	Message string
}

type Notification struct {
	errors []ValidationMessage
}

func NewNotification() Notification {
	return Notification{errors: make([]ValidationMessage, 0)}
}

func (n *Notification) Add(message string) {
	n.errors = append(n.errors, ValidationMessage{Message: message})
}

func (n *Notification) AddField(field, message string) {
	n.errors = append(n.errors, ValidationMessage{Field: field, Message: message})
}

func (n *Notification) Merge(other Notification) {
	n.errors = append(n.errors, other.errors...)
}

func (n Notification) IsValid() bool {
	return len(n.errors) == 0
}

func (n Notification) HasErrors() bool {
	return !n.IsValid()
}

func (n Notification) Errors() []ValidationMessage {
	out := make([]ValidationMessage, len(n.errors))
	copy(out, n.errors)
	return out
}

func (n Notification) IntoError() error {
	if n.HasErrors() {
		return ValidationFailed(n)
	}
	return nil
}

func (n Notification) String() string {
	var sb strings.Builder
	for _, err := range n.errors {
		if err.Field == "" {
			sb.WriteString("  - " + err.Message + "\n")
			continue
		}
		sb.WriteString("  - " + err.Field + ": " + err.Message + "\n")
	}
	return sb.String()
}
