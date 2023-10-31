package accounting

import "github.com/fgrid/money"

type TaxEvent struct {
	*EventBase
}

func NewTaxEvent(base *EventBase, taxableAmount *money.Money) *EventBase {
	if base.eventType == EventTypeTax {
		panic("probable endless recursion")
	}
	t := &TaxEvent{
		EventBase: NewMonetaryEvent(taxableAmount, EventTypeTax, base.whenOccurred, base.whenNoticed, base.customer),
	}
	base.friendAddSecondaryEvent(t.EventBase)
	return t.EventBase
}
