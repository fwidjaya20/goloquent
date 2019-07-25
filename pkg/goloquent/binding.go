package goloquent

// Binding .
type Binding struct {
	Aggregate  *Aggregate
	Conditions []*Condition
	Limit      int
	Offset     int
	GroupBy    []string
	Order      *Order
}
