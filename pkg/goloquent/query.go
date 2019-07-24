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

// Where .
func (q *Query) Where(key string, op Operator, value interface{}) *Query {
	cond := newCondition(AND, key, op, value)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// OrWhere .
func (q *Query) OrWhere(key string, op Operator, value interface{}) *Query {
	cond := newCondition(OR, key, op, value)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// WhereIn .
func (q *Query) WhereIn(key string, value interface{}) *Query {
	cond := newCondition(AND, key, IN, value)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// Except .
func (q *Query) Except(key string, value interface{}) *Query {
	cond := newCondition(AND, key, NOT_IN, value)

	q.Binding.Conditions = append(q.Binding.Conditions, cond)

	return q
}

// ToSQL .
func (q *Query) ToSQL() string {
	return q.Builder.BuildSelect(q.Model, q.Binding)
}

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
	fmt.Println(q.Tx)

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
		default:
			key := fmt.Sprintf("%d%s", i, v.Column)

			payload[key] = v.Value
		}
	}

	return payload
}
