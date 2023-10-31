package accounting

import (
	"time"

	"github.com/fgrid/money"
	"github.com/fgrid/uuid"
)

type Entry struct {
	ID      string
	EventID string
	Stamp   time.Time
	Type    EntryType
	Amount  *money.Money
	account *Account
}

func NewEntry(eventID string, stamp time.Time, entryType EntryType, amount *money.Money) *Entry {
	return &Entry{
		ID:      uuid.NewV4().String(),
		EventID: eventID,
		Stamp:   stamp,
		Type:    entryType,
		Amount:  amount,
	}
}

func NewTwoLeggedEntry(amount *money.Money, stamp time.Time) *Entry {
	return &Entry{
		ID:     uuid.NewV4().String(),
		Amount: amount,
		Stamp:  stamp,
	}
}

func NewMultiLeggedEntry(amount *money.Money, stamp time.Time, account *Account) *Entry {
	return &Entry{
		ID:      uuid.NewV4().String(),
		Stamp:   stamp,
		Amount:  amount,
		account: account,
	}
}

func (e *Entry) post() {
	e.account.AddEntry(e)
}
