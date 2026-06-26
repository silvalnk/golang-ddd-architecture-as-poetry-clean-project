package main

import "fmt"

type CartItem struct {
	product  string
	quantity uint32
}

func newCartItem(product string, quantity uint32) CartItem {
	return CartItem{product: product, quantity: quantity}
}

func (i CartItem) lineTotal(unitPrice uint32) uint32 {
	return unitPrice * i.quantity
}

type ShoppingCart struct {
	id    uint32
	items []CartItem
}

func NewShoppingCart(id uint32) *ShoppingCart {
	return &ShoppingCart{id: id, items: make([]CartItem, 0)}
}

func (c *ShoppingCart) AddItem(product string, quantity uint32) error {
	if quantity == 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}
	for _, item := range c.items {
		if item.product == product {
			return fmt.Errorf("'%s' already in cart", product)
		}
	}
	c.items = append(c.items, newCartItem(product, quantity))
	return nil
}

func (c *ShoppingCart) RemoveItem(product string) error {
	filtered := make([]CartItem, 0, len(c.items))
	removed := false
	for _, item := range c.items {
		if item.product == product {
			removed = true
			continue
		}
		filtered = append(filtered, item)
	}
	if !removed {
		return fmt.Errorf("'%s' not found", product)
	}
	c.items = filtered
	return nil
}

func (c *ShoppingCart) Items() []CartItem {
	out := make([]CartItem, len(c.items))
	copy(out, c.items)
	return out
}

func (c *ShoppingCart) TotalUnits() uint32 {
	var total uint32
	for _, item := range c.items {
		total += item.quantity
	}
	return total
}

func main() {
	fmt.Println("=== Aggregate por composicao ===")

	cart := NewShoppingCart(1)
	_ = cart.AddItem("Coffee", 2)
	_ = cart.AddItem("Book", 1)

	fmt.Printf("Cart #%d\n", cart.id)
	fmt.Printf("  lines: %d\n", len(cart.Items()))
	fmt.Printf("  total units: %d\n", cart.TotalUnits())

	for _, item := range cart.Items() {
		fmt.Printf("  - %s x%d (subtotal @10: %d)\n", item.product, item.quantity, item.lineTotal(10))
	}

	if err := cart.AddItem("Coffee", 1); err != nil {
		fmt.Printf("add rejected (invariant): %s\n", err)
	}

	_ = cart.RemoveItem("Book")
	fmt.Printf("after remove: %d line(s) left\n", len(cart.Items()))
}
