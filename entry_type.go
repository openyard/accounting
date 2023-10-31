package accounting

type EntryType int

const (
	_ EntryType = iota
	EntryTypeBaseUsage
	EntryTypeService
	EntryTypeTax
)

var EntryTypes = []EntryType{
	EntryTypeBaseUsage,
	EntryTypeService,
	EntryTypeTax,
}

func (et EntryType) String() string {
	return [...]string{"", "Base Usage", "Service Fee", "Tax"}[et]
}
