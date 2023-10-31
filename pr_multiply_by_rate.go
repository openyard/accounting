package accounting

import "github.com/fgrid/money"

type MultiplyByRatePR struct {
	*PostingRuleBase
}

func NewMultiplyByRatePR(entryType EntryType) *MultiplyByRatePR {
	pr := &MultiplyByRatePR{}
	pr.PostingRuleBase = NewPostingRule(entryType).
		WithCalculateAmountFunc(pr.calculateAmount).
		WithIsTaxableFunc(pr.isTaxable)
	return pr
}

func (pr *MultiplyByRatePR) calculateAmount(event Event) *money.Money {
	usage := event.(*Usage)
	return money.EUR(usage.Quantity.Amount).MulFloat64(usage.Rate())
}

func (pr *MultiplyByRatePR) isTaxable() bool {
	return !(pr.entryType == EntryTypeTax)
}
