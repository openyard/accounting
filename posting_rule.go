package accounting

import "github.com/fgrid/money"

type CalculateAmountFunc func(event Event) *money.Money
type IsTaxableFunc func() bool

type PostingRule interface {
	Process(base *EventBase, event Event)
}

type PostingRuleBase struct {
	PostingRule
	entryType       EntryType
	calculateAmount CalculateAmountFunc
	isTaxable       IsTaxableFunc
}

func NewPostingRule(entryType EntryType) *PostingRuleBase {
	return &PostingRuleBase{entryType: entryType}
}

func (pr *PostingRuleBase) WithCalculateAmountFunc(f CalculateAmountFunc) *PostingRuleBase {
	return &PostingRuleBase{
		entryType:       pr.entryType,
		calculateAmount: f,
		isTaxable:       pr.isTaxable,
	}
}

func (pr *PostingRuleBase) WithIsTaxableFunc(f IsTaxableFunc) *PostingRuleBase {
	return &PostingRuleBase{
		entryType:       pr.entryType,
		calculateAmount: pr.calculateAmount,
		isTaxable:       f,
	}
}

func (pr *PostingRuleBase) Process(base *EventBase, event Event) {
	pr.makeEntry(base, pr.calculateAmount(event))
	if pr.isTaxable() {
		NewTaxEvent(base, pr.calculateAmount(event)).Process()
	}
}

func (pr *PostingRuleBase) makeEntry(base *EventBase, amount *money.Money) {
	entry := NewEntry(base.ID, base.whenOccurred, pr.entryType, amount)
	base.customer.AddEntry(entry, pr.entryType)
	base.resultingEntries[entry.ID] = entry
}
