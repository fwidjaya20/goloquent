package goloquent

// Condition is a struct that is used for store query logic or condition
type Condition struct {
	Connector       Connector
	Column          string
	Operator        Operator
	Value           interface{}
	IsCompareColumn bool
	ColumnCompare   string
}

func newCondition(connector Connector, column string, operator Operator, value interface{}) *Condition {
	return &Condition{
		Connector:       connector,
		Column:          column,
		Operator:        operator,
		Value:           value,
		IsCompareColumn: false,
	}
}

func newCompareColumnCondition(connector Connector, targetColumn string, operator Operator, sourceColumn string) *Condition {
	return &Condition{
		Connector:       connector,
		Column:          targetColumn,
		Operator:        operator,
		Value:           nil,
		IsCompareColumn: true,
		ColumnCompare:   sourceColumn,
	}
}
