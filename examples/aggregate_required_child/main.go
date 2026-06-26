package main

import "fmt"

type Address struct {
	street string
	city   string
}

func NewAddress(street, city string) Address {
	return Address{street: street, city: city}
}

type GiftWrap struct {
	message  string
	feeCents uint32
}

func NewGiftWrap(message string, feeCents uint32) GiftWrap {
	return GiftWrap{message: message, feeCents: feeCents}
}

type Order struct {
	id             uint32
	customerID     uint32
	billingAddress Address
	giftWrap       *GiftWrap
}

func NewOrder(id uint32, customerID uint32, billingAddress Address) *Order {
	return &Order{
		id:             id,
		customerID:     customerID,
		billingAddress: billingAddress,
	}
}

func (o *Order) AddGiftWrap(message string) error {
	if o.giftWrap != nil {
		return fmt.Errorf("gift wrap already set")
	}
	wrap := NewGiftWrap(message, 500)
	o.giftWrap = &wrap
	return nil
}

func (o *Order) RemoveGiftWrap() error {
	if o.giftWrap == nil {
		return fmt.Errorf("no gift wrap to remove")
	}
	o.giftWrap = nil
	return nil
}

func (o *Order) ChangeBillingCity(city string) {
	o.billingAddress.city = city
}

func (o *Order) TotalCents(itemsTotal uint32) uint32 {
	if o.giftWrap == nil {
		return itemsTotal
	}
	return itemsTotal + o.giftWrap.feeCents
}

type Customer struct {
	id   uint32
	name string
}

func main() {
	fmt.Println("=== Aggregate: relacionamento sem slice ===")

	customer := Customer{id: 42, name: "Alice"}
	order := NewOrder(1, customer.id, NewAddress("Rua A, 100", "Sao Paulo"))

	fmt.Printf("Order #%d\n", order.id)
	fmt.Printf("  customer_id: %d (%s)\n", order.customerID, customer.name)
	fmt.Printf("  billing (1:1): %s, %s\n", order.billingAddress.street, order.billingAddress.city)

	_ = order.AddGiftWrap("Happy birthday!")
	fmt.Println("after add_gift_wrap:")
	if order.giftWrap != nil {
		fmt.Printf("  message: %s\n", order.giftWrap.message)
		fmt.Printf("  fee: %d cents\n", order.giftWrap.feeCents)
	}

	if err := order.AddGiftWrap("Again"); err != nil {
		fmt.Printf("  rejected: %s\n", err)
	}

	order.ChangeBillingCity("Campinas")
	fmt.Printf("billing city updated: %s\n", order.billingAddress.city)
	fmt.Printf("total with gift: %d cents\n", order.TotalCents(10000))

	_ = order.RemoveGiftWrap()
	fmt.Printf("gift_wrap after remove: %v\n", order.giftWrap)
}
