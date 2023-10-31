package accounting

import "github.com/fgrid/money"

type PoorCapPR struct {
	*PostingRuleBase
	rate       float64
	usageLimit *Quantity
}

func NewPoorCapPR(entryType EntryType, rate float64, usageLimit *Quantity) *PoorCapPR {
	pr := &PoorCapPR{rate: rate, usageLimit: usageLimit}
	pr.PostingRuleBase = NewPostingRule(entryType).
		WithCalculateAmountFunc(pr.calculateAmount).
		WithIsTaxableFunc(pr.isTaxable)
	return pr
}

func (pr *PoorCapPR) calculateAmount(event Event) *money.Money {
	usage := event.(*Usage)
	amountUsed := usage.Quantity
	if amountUsed.GT(pr.usageLimit) {
		return money.EUR(amountUsed.Amount).MulFloat64(usage.Rate())
	}
	return money.EUR(amountUsed.Amount).MulFloat64(pr.rate)
}

func (pr *PoorCapPR) isTaxable() bool {
	return false
}
