package accounting

import (
	"fmt"
	"time"

	"github.com/fgrid/money"
)

type Transaction struct {
	stamp     time.Time
	entries   []*Entry
	wasPosted bool
}

func NewTwoLeggedTransaction(amount *money.Money, from, to *Account, stamp time.Time) *Transaction {
	entryFrom := NewTwoLeggedEntry(amount.Debit(), stamp)
	entryTo := NewTwoLeggedEntry(amount.Credit(), stamp)
	from.AddEntry(entryFrom)
	to.AddEntry(entryTo)
	t := &Transaction{
		stamp:     stamp,
		entries:   []*Entry{entryFrom, entryTo},
		wasPosted: true,
	}
	return t
}

func NewMultiLeggedTransaction(stamp time.Time) *Transaction {
	return &Transaction{
		stamp:     stamp,
		wasPosted: false,
	}
}

func (t *Transaction) Add(amount *money.Money, account *Account) error {
	if t.wasPosted {
		return fmt.Errorf("cannot add entry to a transaction that's already posted")
	}
	t.entries = append(t.entries, NewMultiLeggedEntry(amount, t.stamp, account))
	return nil
}

func (t *Transaction) Post() error {
	if !t.canPost() {
		return fmt.Errorf("unable to post")
	}
	for _, entry := range t.entries {
		entry.post()
	}
	t.wasPosted = true
	return nil
}

func (t *Transaction) canPost() bool {
	return isZero(t.balance())
}

func (t *Transaction) balance() *money.Money {
	if len(t.entries) == 0 {
		return money.EUR(0)
	}

	result := money.EUR(0)
	for _, entry := range t.entries {
		result, _ = result.Add(entry.Amount)
	}
	return result
}

func isZero(amount *money.Money) bool {
	return amount.Cents() == 0
}
