package accounting

import "time"

type ServiceAgreement struct {
	Rate         float64
	postingRules map[EventType]*TemporalCollection
}

func NewServiceAgreement(rate float64) *ServiceAgreement {
	return &ServiceAgreement{
		Rate:         rate,
		postingRules: make(map[EventType]*TemporalCollection, 0),
	}
}

func (sa *ServiceAgreement) AddPostingRule(eventType EventType, rule PostingRule, date time.Time) {
	sa.temporalCollection(eventType).Put(date, rule)
}

func (sa *ServiceAgreement) getPostingRule(eventType EventType, when time.Time) PostingRule {
	pr, _ := sa.temporalCollection(eventType).Get(when)
	return pr.(PostingRule)
}

func (sa *ServiceAgreement) temporalCollection(eventType EventType) *TemporalCollection {
	_, found := sa.postingRules[eventType]
	if !found {
		sa.postingRules[eventType] = NewTemporalCollection()
	}
	return sa.postingRules[eventType]
}
