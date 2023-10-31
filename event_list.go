package accounting

type EventList struct {
	events []*EventBase
}

func NewEventList() *EventList {
	return &EventList{events: make([]*EventBase, 0)}
}

func (el *EventList) Process() {
	for _, event := range el.events {
		event.Process()
	}
}

func (el *EventList) Add(event *EventBase) {
	el.events = append(el.events, event)
}
