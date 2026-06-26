package shared

import "fmt"

const CurrencyBRL = "BRL"

type Money struct {
	cents    int64
	currency string
}

func NewMoneyBRL(reais int64, centavos uint8) (Money, error) {
	if centavos >= 100 {
		return Money{}, InvalidValue("centavos", "must be less than 100")
	}
	return Money{
		cents:    reais*100 + int64(centavos),
		currency: CurrencyBRL,
	}, nil
}

func NewMoneyFromCents(cents int64) (Money, error) {
	if cents < 0 {
		return Money{}, InvalidValue("cents", "cannot be negative")
	}
	return Money{
		cents:    cents,
		currency: CurrencyBRL,
	}, nil
}

func (m Money) Cents() int64 {
	return m.cents
}

func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, InvalidOperation("cannot combine different currencies")
	}
	return Money{
		cents:    m.cents + other.cents,
		currency: m.currency,
	}, nil
}

func (m Money) Mul(quantity uint32) Money {
	return Money{
		cents:    m.cents * int64(quantity),
		currency: m.currency,
	}
}

func (m Money) IsGreaterThan(other Money) (bool, error) {
	if m.currency != other.currency {
		return false, InvalidOperation("cannot compare different currencies")
	}
	return m.cents > other.cents, nil
}

func (m Money) String() string {
	reais := m.cents / 100
	centavos := m.cents % 100
	if centavos < 0 {
		centavos = -centavos
	}
	return fmt.Sprintf("R$ %d,%02d", reais, centavos)
}
