package goloquent

// Condition is a struct that is used for store query logic or condition
type Condition struct {
	Connector Connector
	Column    string
	Operator  Operator
	Value     interface{}
}

func newCondition(connector Connector, column string, operator Operator, value interface{}) *Condition {
	return &Condition{
		Connector: connector,
		Column:    column,
		Operator:  operator,
		Value:     value,
	}
}
