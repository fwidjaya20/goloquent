package goloquent

// Count is an aggregate function for retrive row count
func (q *Query) Count() int64 {
	defer q.resetBindings()

	defer q.resetBindings()

	q.Binding.Aggregate = newAggregate(COUNT, "*")

	return int64(q.execAggregate())
}

// Max is an aggregate function for retrive column Max malue
func (q *Query) Max(column string) float64 {
	defer q.resetBindings()

	defer q.resetBindings()

	q.Binding.Aggregate = newAggregate(MAX, column)

	return q.execAggregate()
}

// Min is an aggregate function for retrive column Min malue
func (q *Query) Min(column string) float64 {
	defer q.resetBindings()

	defer q.resetBindings()

	q.Binding.Aggregate = newAggregate(MIN, column)

	return q.execAggregate()
}

// Avg is an aggregate function for retrive column Avg malue
func (q *Query) Avg(column string) float64 {
	defer q.resetBindings()

	defer q.resetBindings()

	q.Binding.Aggregate = newAggregate(AVG, column)

	return q.execAggregate()
}

// Sum is an aggregate function for retrive column Sum malue
func (q *Query) Sum(column string) float64 {
	defer q.resetBindings()

	q.Binding.Aggregate = newAggregate(SUM, column)

	return q.execAggregate()
}

func (q *Query) execAggregate() float64 {
	var result float64

	stmt, args, _ := q.prepareNamed(q.ToSQL())

	stmt.Get(&result, args)

	return result
}
