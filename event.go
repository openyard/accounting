package accounting

import (
	"fmt"
	"log"
	"time"

	"github.com/fgrid/uuid"
)

type Event interface{}

type EventBase struct {
	ID               string
	event            Event
	eventType        EventType
	whenOccurred     time.Time
	whenNoticed      time.Time
	customer         *Customer
	resultingEntries map[string]*Entry
	secondaryEvents  []*EventBase
	adjustedEvent    *EventBase
	replacementEvent *EventBase
	isProcessed      bool
}

func NewEvent(event Event, eventType EventType, whenOccurred, whenNoticed time.Time, customer *Customer) *EventBase {
	return &EventBase{
		ID:               uuid.NewV4().String(),
		event:            event,
		eventType:        eventType,
		whenOccurred:     whenOccurred,
		whenNoticed:      whenNoticed,
		customer:         customer,
		resultingEntries: make(map[string]*Entry, 0),
		secondaryEvents:  make([]*EventBase, 0),
		isProcessed:      false,
	}
}

func NewEventAdjustment(event Event, eventType EventType, whenOccurred, whenNoticed time.Time, adjustedEvent *EventBase) *EventBase {
	if adjustedEvent.hasBeenAdjusted() {
		panic(fmt.Sprintf("the adjusted event %T is already adjusted", adjustedEvent))
	}
	e := &EventBase{
		ID:               uuid.NewV4().String(),
		event:            event,
		eventType:        eventType,
		whenOccurred:     whenOccurred,
		whenNoticed:      whenNoticed,
		customer:         adjustedEvent.customer,
		resultingEntries: make(map[string]*Entry, 0),
		secondaryEvents:  make([]*EventBase, 0),
		adjustedEvent:    adjustedEvent,
		isProcessed:      false,
	}
	e.adjustedEvent.replacementEvent = e
	return e
}

func (e *EventBase) AddResultingEntry(entry *Entry) {
	e.resultingEntries[entry.ID] = entry
}

func (e *EventBase) AllResultingEntries() map[string]*Entry {
	result := e.ResultingEntries()
	for _, event := range e.secondaryEvents {
		for ID, entry := range event.resultingEntries {
			result[ID] = entry
		}
	}
	return result
}

func (e *EventBase) ResultingEntries() map[string]*Entry {
	return e.resultingEntries
}

func (e *EventBase) Process() {
	if e.isProcessed {
		log.Printf("cannot process event twice")
		return
	}
	if e.adjustedEvent != nil {
		e.adjustedEvent.reverse()
	}
	rule := e.findRule()
	rule.Process(e, e.event)
	e.isProcessed = true
}

func (e *EventBase) findRule() PostingRule {
	return e.customer.serviceAgreement.getPostingRule(e.eventType, e.whenOccurred)
}

func (e *EventBase) friendAddSecondaryEvent(event *EventBase) {
	e.secondaryEvents = append(e.secondaryEvents, event)
}

func (e *EventBase) reverse() {
	resultingEntries := make(map[string]*Entry, 0)
	for _, entry := range e.resultingEntries {
		reversingEntry := NewEntry(e.ID, e.whenNoticed, entry.Type, entry.Amount.Inv())
		e.customer.AddEntry(reversingEntry, entry.Type)
		resultingEntries[reversingEntry.ID] = reversingEntry
	}
	e.resultingEntries = resultingEntries
	e.reverseSecondaryEvents()
}

func (e *EventBase) reverseSecondaryEvents() {
	for _, event := range e.secondaryEvents {
		event.reverse()
	}
}

func (e *EventBase) hasBeenAdjusted() bool {
	return e.replacementEvent != nil
}
