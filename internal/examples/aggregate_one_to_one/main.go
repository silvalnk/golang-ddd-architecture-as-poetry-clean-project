package main

import "fmt"

type Money struct {
	cents uint32
}

func NewMoneyBRL(reais uint32, centavos uint8) Money {
	return Money{cents: reais*100 + uint32(centavos)}
}

type ContractTerms struct {
	planName   string
	monthlyFee Money
	minMonths  uint32
}

func StarterContractTerms() ContractTerms {
	return ContractTerms{
		planName:   "Starter",
		monthlyFee: NewMoneyBRL(49, 90),
		minMonths:  12,
	}
}

func (t ContractTerms) EarlyExitPenalty(monthsLeft uint32) Money {
	return Money{cents: t.monthlyFee.cents * monthsLeft / 2}
}

type AccountHolder struct {
	fullName string
	document string
}

func NewAccountHolder(fullName, document string) AccountHolder {
	return AccountHolder{fullName: fullName, document: document}
}

type Subscription struct {
	id         uint32
	holder     AccountHolder
	terms      ContractTerms
	monthsPaid uint32
}

func NewSubscription(id uint32, holder AccountHolder, terms ContractTerms) *Subscription {
	return &Subscription{id: id, holder: holder, terms: terms}
}

func (s *Subscription) RegisterPayment() error {
	if s.monthsPaid >= s.terms.minMonths {
		return fmt.Errorf("minimum contract period already fulfilled")
	}
	s.monthsPaid++
	return nil
}

func (s *Subscription) RenameHolder(fullName string) {
	s.holder.fullName = fullName
}

func (s *Subscription) MonthsRemaining() uint32 {
	if s.monthsPaid >= s.terms.minMonths {
		return 0
	}
	return s.terms.minMonths - s.monthsPaid
}

func main() {
	fmt.Println("=== Composicao 1:1 obrigatoria (sem Option) ===")

	holder := NewAccountHolder("Maria Silva", "123.456.789-00")
	terms := StarterContractTerms()
	sub := NewSubscription(100, holder, terms)

	fmt.Printf("Subscription #%d\n", sub.id)
	fmt.Printf("  holder: %s (%s)\n", sub.holder.fullName, sub.holder.document)
	fmt.Printf("  plan: %s - R$ %.2f/month (min %d months)\n", sub.terms.planName, float64(sub.terms.monthlyFee.cents)/100.0, sub.terms.minMonths)

	_ = sub.RegisterPayment()
	_ = sub.RegisterPayment()

	fmt.Printf("months paid: %d\n", sub.monthsPaid)
	fmt.Printf("months remaining: %d\n", sub.MonthsRemaining())
	fmt.Printf("early exit penalty now: R$ %.2f\n", float64(sub.terms.EarlyExitPenalty(sub.MonthsRemaining()).cents)/100.0)

	sub.RenameHolder("Maria Silva Santos")
	fmt.Printf("holder renamed: %s\n", sub.holder.fullName)
}
