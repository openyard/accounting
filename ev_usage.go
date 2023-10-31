package accounting

import "time"

type Usage struct {
	Quantity *Quantity
	customer *Customer
}

func NewUsageEvent(quantity *Quantity, whenOccurred, whenNoticed time.Time, customer *Customer) *EventBase {
	return NewEvent(&Usage{Quantity: quantity, customer: customer}, EventTypeUsage, whenOccurred, whenNoticed, customer)
}

func NewUsageAdjustment(quantity *Quantity, whenOccurred, whenNoticed time.Time, usage *EventBase) *EventBase {
	return NewEventAdjustment(&Usage{Quantity: quantity, customer: usage.customer}, EventTypeUsage, whenOccurred, whenNoticed, usage)
}

func (e *Usage) Rate() float64 {
	return e.customer.serviceAgreement.Rate
}
