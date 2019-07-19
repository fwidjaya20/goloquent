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
