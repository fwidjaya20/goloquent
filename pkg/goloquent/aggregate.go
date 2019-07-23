package goloquent

// Aggregate is a struct for wrapping SQL Aggregate Statement
type Aggregate struct {
	AggregateFunc AggregateFunction
	Column        string
}

func newAggregate(aggregateFn AggregateFunction, column string) *Aggregate {
	return &Aggregate{
		AggregateFunc: aggregateFn,
		Column:        column,
	}
}
