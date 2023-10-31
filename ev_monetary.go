package accounting

import (
	"time"

	"github.com/fgrid/money"
)

type MonetaryEvent struct {
	Amount *money.Money
}

func NewMonetaryEvent(amount *money.Money, eventType EventType, whenOccurred, whenNoticed time.Time, customer *Customer) *EventBase {
	return NewEvent(&MonetaryEvent{Amount: amount}, eventType, whenOccurred, whenNoticed, customer)
}
