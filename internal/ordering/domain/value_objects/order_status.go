package valueobjects

type OrderStatus string

const (
	OrderStatusDraft     OrderStatus = "draft"
	OrderStatusSubmitted OrderStatus = "submitted"
	OrderStatusCancelled OrderStatus = "cancelled"
)

func (s OrderStatus) CanAddLine() bool {
	return s == OrderStatusDraft
}

func (s OrderStatus) CanSubmit() bool {
	return s == OrderStatusDraft
}

func (s OrderStatus) CanCancel() bool {
	return s == OrderStatusDraft || s == OrderStatusSubmitted
}
