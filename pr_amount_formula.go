package accounting

import "github.com/fgrid/money"

type AmountFormulaPR struct {
	*PostingRuleBase
	multiplier float64
	fixedFee   *money.Money
}

func NewAmountFormulaPR(entryType EntryType, multiplier float64, fixedFee *money.Money) *AmountFormulaPR {
	pr := &AmountFormulaPR{multiplier: multiplier, fixedFee: fixedFee}
	pr.PostingRuleBase = NewPostingRule(entryType).
		WithCalculateAmountFunc(pr.calculateAmount).
		WithIsTaxableFunc(pr.isTaxable)
	return pr
}

func (pr *AmountFormulaPR) calculateAmount(event Event) *money.Money {
	amount := event.(*MonetaryEvent).Amount
	result, _ := amount.MulFloat64(pr.multiplier).Add(pr.fixedFee)
	return result
}

func (pr *AmountFormulaPR) isTaxable() bool {
	return false
}
