package goloquent

// BeginTransaction .
func (q *Query) BeginTransaction() *Query {
	var err error

	q.Tx, err = q.DB.Beginx()

	if err != nil {
		q.Rollback()
		panic(err)
	}

	return q
}

// Rollback .
func (q *Query) Rollback() *Query {
	q.Tx.Rollback()
	return q
}

// Commit .
func (q *Query) Commit() *Query {
	q.Tx.Commit()
	return q
}

// EndTransaction .
func (q *Query) EndTransaction() *Query {
	q.Tx = nil
	return q
}
