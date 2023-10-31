package accounting

import (
	"time"

	"github.com/fgrid/money"
	"github.com/fgrid/uuid"
)

type Account struct {
	ID       string
	currency *money.Currency
	entries  []*Entry
}

func NewAccount(cur *money.Currency) *Account {
	return &Account{
		ID:       uuid.NewV4().String(),
		currency: cur,
		entries:  make([]*Entry, 0),
	}
}

func (a *Account) WithID(ID string) *Account {
	return &Account{
		ID:       ID,
		currency: a.currency,
		entries:  make([]*Entry, 0),
	}
}

func (a *Account) AddEntry(entry *Entry) {
	a.post(entry)
}

func (a *Account) Withdraw(amount *money.Money, target *Account, stamp time.Time) {
	_ = NewTwoLeggedTransaction(amount, a, target, stamp)
}

func (a *Account) Balance() *money.Money {
	return a.BalanceFor(UpTo(time.Now()))
}

func (a *Account) BalanceFor(period *TimeRange) *money.Money {
	result := money.EUR(0)
	for _, entry := range a.entries {
		if period.Includes(entry.Stamp) {
			result, _ = result.Add(entry.Amount)
		}
	}
	return result
}

func (a *Account) BalanceAt(timepoint time.Time) *money.Money {
	return a.BalanceFor(UpTo(timepoint))
}

func (a *Account) post(entry *Entry) {
	a.entries = append(a.entries, entry)
}
