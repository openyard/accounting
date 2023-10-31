package accounting

type Quantity struct {
	Amount uint64
	Unit   Unit
}

func NewQuantity(amount uint64, unit Unit) *Quantity {
	return &Quantity{Amount: amount, Unit: unit}
}

func (q Quantity) GT(limit *Quantity) bool {
	return q.Amount > limit.Amount
}
