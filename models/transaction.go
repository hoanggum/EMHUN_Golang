package models

import "fmt"

type Transaction struct {
	Items              []int
	Utilities          []int
	TransactionUtility int
}

func NewTransaction(items []int, utilities []int, transUtility int) *Transaction {
	return &Transaction{
		Items:              items,
		Utilities:          utilities,
		TransactionUtility: transUtility,
	}
}

func (t *Transaction) GetItems() []int {
	return t.Items
}

func (t *Transaction) GetUtilities() []int {
	return t.Utilities
}

func (t *Transaction) GetTransactionUtility() int {
	return t.TransactionUtility
}
func (t *Transaction) String() string {
	return fmt.Sprintf("Các item: %v | Utilities: %v | Transaction Utility: %d", t.Items, t.Utilities, t.TransactionUtility)
}
