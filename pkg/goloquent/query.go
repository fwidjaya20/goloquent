package goloquent

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/jmoiron/sqlx"
)

// Query .
type Query struct {
	Builder *Builder
	DB      *sqlx.DB
	Tx      *sqlx.Tx
	Model   IModel
	Binding Binding
}

// DB .
func DB(db *sqlx.DB) *Query {
	return &Query{
		DB:      db,
		Builder: NewBuilder(),
	}
}

// Use .
func (q *Query) Use(model IModel) *Query {
	q.Model = model

	return q
}

// GroupBy methods may be used to group the query results
func (q *Query) GroupBy(columns ...string) *Query {
	q.Binding.GroupBy = columns

	return q
}

// Skip method used to skip a given number of results in the query
func (q *Query) Skip(amount int) *Query {
	q.Binding.Offset = amount

	return q
}

// Take method is used to limit the number of results returned from the query
func (q *Query) Take(amount int) *Query {
	q.Binding.Limit = amount

	return q
}

// OrderBy method allows you to sort the result of the query by a given column.
// The first argument to the orderBy method should be the column you wish to sort by, while the second argument controls the direction of the sort and may be either asc or desc
func (q *Query) OrderBy(direction OrderDirection, columns ...string) *Query {
	q.Binding.Order = &Order{
		Columns:   columns,
		Direction: direction,
	}

	return q
}

// ToSQL method will generate Statement Binding into SQL Query
func (q *Query) ToSQL() string {
	return q.Builder.BuildSelect(q.Model, q.Binding)
}

func (q *Query) generateInsertColumn() []string {
	var columns []string
	typeOf := reflect.TypeOf(q.Model)

	for i := 0; i < typeOf.Elem().NumField(); i++ {
		column := typeOf.Elem().Field(i)

		if "Model" != column.Name {
			tag := typeOf.Elem().Field(i).Tag.Get("db")

			columns = append(columns, tag)
		}
	}

	return columns
}

func (q *Query) generateBulkInsertColumn(data []interface{}) []string {
	var columns []string

	for i := 0; i < len(data); i++ {
		cols := q.generateInsertColumn()

		columns = append(columns, cols...)
	}

	return columns
}

// bulkPayload is a function that will merge all payload for bulkinsert into sequential slice of map[string]interface{}
func (q *Query) bulkPayload(data []interface{}) map[string]interface{} {
	payloads := map[string]interface{}{}

	for i, v := range data {
		model := v.(IModel)
		payload := model.MapToPayload(model)

		for key, value := range payload {
			payloads[fmt.Sprintf("%d%s", i, key)] = value

			if model.IsTimestamp() {
				payloads[fmt.Sprintf("%d%s", i, CREATED_AT)] = time.Now()
			}
		}
	}

	return payloads
}

func (q *Query) makeSliceOf(sample interface{}) (interface{}, error) {
	if reflect.TypeOf(sample).Kind() != reflect.Ptr {
		return nil, errors.New("sample must be a pointer to reference model")
	}

	element := reflect.TypeOf(sample)

	return reflect.New(reflect.SliceOf(element)).Interface(), nil
}

func (q *Query) makeTypeOf(sample interface{}) (interface{}, error) {
	if reflect.TypeOf(sample).Kind() != reflect.Ptr {
		return nil, errors.New("sample must be a pointer to reference model")
	}

	element := reflect.TypeOf(sample).Elem()

	return reflect.New(element).Interface(), nil
}

func (q *Query) mapToSliceModel(slice interface{}) interface{} {
	return reflect.ValueOf(slice).Elem().Interface()
}

func (q *Query) assignModel(result interface{}, m Model) interface{} {
	value := reflect.ValueOf(result).Interface()

	model := reflect.ValueOf(m)

	assignedModel := reflect.ValueOf(value)
	assignedModel.Elem().FieldByName("Model").Set(model)

	return assignedModel.Interface()
}

func (q *Query) resetBindings() {
	q.Binding = Binding{}
}

func (q *Query) prepareNamed(query string) (*sqlx.NamedStmt, map[string]interface{}, error) {
	stmt, err := q.DB.PrepareNamed(query)

	return stmt, q.mapConditionPayload(), err
}

func (q *Query) mapConditionPayload() map[string]interface{} {
	payload := map[string]interface{}{}

	for i, v := range q.Binding.Conditions {
		switch v.Operator {
		case IN, NOT_IN:
			payload = q.Builder.buildInValue(payload, i, v)
		case BETWEEN, NOT_BETWEEN:
			payload = q.Builder.buildBetweenValue(payload, i, v)
		default:
			if !v.IsCompareColumn {
				key := fmt.Sprintf("%d%s", i, v.Column)

				payload[key] = v.Value
			}
		}
	}

	return payload
}
