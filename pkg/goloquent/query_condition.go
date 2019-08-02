package goloquent

// Where methods will compare the column with a value.
func (q *Query) Where(column string, op Operator, value interface{}) *Query {
	cond := newCondition(AND, column, op, value)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// OrWhere methods will compare the column with a value.
func (q *Query) OrWhere(column string, op Operator, value interface{}) *Query {
	cond := newCondition(OR, column, op, value)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// WhereIn method verifies that a given column's value is contained within the given array
func (q *Query) WhereIn(column string, value interface{}) *Query {
	cond := newCondition(AND, column, IN, value)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// OrWhereIn method verifies that a given column's value is contained within the given array
func (q *Query) OrWhereIn(column string, value interface{}) *Query {
	cond := newCondition(OR, column, IN, value)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// Except method verifies that a given column's value is not contained within the given array
func (q *Query) Except(column string, value interface{}) *Query {
	cond := newCondition(AND, column, NOT_IN, value)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// OrExcept method verifies that a given column's value is not contained within the given array
func (q *Query) OrExcept(column string, value interface{}) *Query {
	cond := newCondition(OR, column, NOT_IN, value)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// WhereBetween method verifies that a column's value is between two values
func (q *Query) WhereBetween(column string, firstValue interface{}, lastValue interface{}) *Query {
	cond := newCondition(AND, column, BETWEEN, []interface{}{firstValue, lastValue})

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// OrWhereBetween method verifies that a column's value is between two values
func (q *Query) OrWhereBetween(column string, firstValue interface{}, lastValue interface{}) *Query {
	cond := newCondition(OR, column, BETWEEN, []interface{}{firstValue, lastValue})

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// WhereNotBetween method verifies that a column's value lies outside of two values
func (q *Query) WhereNotBetween(column string, firstValue interface{}, lastValue interface{}) *Query {
	cond := newCondition(AND, column, NOT_BETWEEN, []interface{}{firstValue, lastValue})

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// OrWhereNotBetween method verifies that a column's value lies outside of two values
func (q *Query) OrWhereNotBetween(column string, firstValue interface{}, lastValue interface{}) *Query {
	cond := newCondition(OR, column, NOT_BETWEEN, []interface{}{firstValue, lastValue})

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// WhereNull method verifies that the value of the given column is NULL
func (q *Query) WhereNull(column string) *Query {
	cond := newCondition(AND, column, IS_NULL, nil)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// OrWhereNull method verifies that the value of the given column is NULL
func (q *Query) OrWhereNull(column string) *Query {
	cond := newCondition(OR, column, IS_NULL, nil)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// WhereNotNull method verifies that the value of the given column is NOT NULL
func (q *Query) WhereNotNull(column string) *Query {
	cond := newCondition(AND, column, IS_NOT_NULL, nil)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// OrWhereNotNull method verifies that the value of the given column is NOT NULL
func (q *Query) OrWhereNotNull(column string) *Query {
	cond := newCondition(OR, column, IS_NOT_NULL, nil)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// WhereColumn method may be used to verify that two columns are equal
func (q *Query) WhereColumn(target string, op Operator, source string) *Query {
	cond := newCompareColumnCondition(AND, target, op, source)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// OrWhereColumn method may be used to verify that two columns are equal
func (q *Query) OrWhereColumn(target string, op Operator, source string) *Query {
	cond := newCompareColumnCondition(OR, target, op, source)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}
