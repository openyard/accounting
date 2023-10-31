package accounting

type EventType int

const (
	_ EventType = iota
	EventTypeUsage
	EventTypeServiceCall
	EventTypeTax
)

func (et EventType) String() string {
	return [...]string{"", "Usage", "Service Call", "Tax"}[et]
}
