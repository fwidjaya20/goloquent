package goloquent

import "github.com/jmoiron/sqlx"

// Query .
type Query struct {
	DB    *sqlx.DB
	Model IModel
}

// DB .
func DB(db *sqlx.DB) *Query {
	return &Query{
		DB: db,
	}
}

// Use .
func (q *Query) Use(model IModel) *Query {
	q.Model = model

	return q
}
