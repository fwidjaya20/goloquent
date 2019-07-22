package goloquent

import (
	"database/sql"
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

// All .
func (q *Query) All() (interface{}, error) {
	q.resetBindings()

	return q.Get()
}

// Get .
func (q *Query) Get() (interface{}, error) {
	defer q.resetBindings()

	results, err := q.makeSliceOf(q.Model)

	if nil != err {
		return nil, err
	}

	stmt, args, err := q.prepareNamed(q.ToSQL())

	fmt.Println(q.ToSQL())

	err = stmt.Select(results, args)

	return q.mapToSliceModel(results), err
}

// First .
func (q *Query) First() (interface{}, error) {
	defer q.resetBindings()

	q.Binding.Limit = 1

	result, err := q.makeTypeOf(q.Model)

	if nil != err {
		return nil, err
	}

	stmt, args, err := q.prepareNamed(q.ToSQL())

	err = stmt.Get(result, args)

	return reflect.ValueOf(result).Interface(), err
}

// Paginate .
func (q *Query) Paginate(page int, limit ...int) (interface{}, error) {
	defer q.resetBindings()

	if len(limit) > 0 {
		q.Binding.Limit = limit[0]
	} else {
		q.Binding.Limit = 50
	}

	q.Binding.Offset = (page - 1) * q.Binding.Limit

	return q.Get()
}

// Insert .
func (q *Query) Insert(returning ...string) (sql.Result, error) {
	var result sql.Result
	var err error

	query := q.Builder.BuildInsert(q.Model, returning...)

	payload := q.Model.MapToPayload(q.Model)

	if q.Model.IsTimestamp() {
		payload[CREATED_AT] = time.Now()
	}

	if nil != q.Tx {
		result, err = q.Tx.NamedExec(query, payload)
	} else {
		result, err = q.DB.NamedExec(query, payload)
	}

	return result, err
}

// BulkInsert .
func (q *Query) BulkInsert(data interface{}, returning ...string) (sql.Result, error) {
	var result sql.Result
	var err error
	var value reflect.Value

	value = reflect.ValueOf(data)

	if reflect.Slice != value.Kind() {
		return nil, errors.New("data must be a slice")
	}

	slice := make([]interface{}, value.Len())

	for i := 0; i < value.Len(); i++ {
		slice[i] = value.Index(i).Interface()
	}

	query := q.Builder.BuildBulkInsert(q.Model, slice, returning...)
	payloads := q.bulkPayload(slice)

	if nil != q.Tx {
		result, err = q.Tx.NamedExec(query, payloads)
	} else {
		result, err = q.DB.NamedExec(query, payloads)
	}

	return result, err
}

// RawCommand .
func (q *Query) RawCommand(query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	var err error

	if nil != q.Tx {
		result, err = q.Tx.Exec(query, args...)
	} else {
		result, err = q.DB.Exec(query, args...)
	}

	return result, err
}

// RawQuery .
func (q *Query) RawQuery(dest IModel, query string, args ...interface{}) error {
	var err error

	err = q.DB.Select(dest, query, args...)

	return err
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
	payload := make(map[string]interface{}, len(q.Binding.Conditions))

	for i, v := range q.Binding.Conditions {
		key := fmt.Sprintf("%d%s", i, v.Column)

		payload[key] = v.Value
	}

	return payload
}
