package goloquent

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/jmoiron/sqlx"
)

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

	fmt.Println(err)
	fmt.Println(q.ToSQL())
	fmt.Println(args)

	err = stmt.Select(results, args)

	return q.mapToSliceModel(results), err
}

// Find .
func (q *Query) Find(value interface{}) (interface{}, error) {
	defer q.resetBindings()

	q.Where(q.Model.GetPK(), EQUAL, value)
	q.Take(1)

	result, err := q.makeTypeOf(q.Model)

	if nil != err {
		return nil, err
	}

	stmt, args, err := q.prepareNamed(q.ToSQL())

	err = stmt.Get(result, args)

	return q.assignModel(result, q.Model.GetModel()), err
}

// First .
func (q *Query) First() (interface{}, error) {
	defer q.resetBindings()

	q.Take(1)

	result, err := q.makeTypeOf(q.Model)

	if nil != err {
		return nil, err
	}

	stmt, args, err := q.prepareNamed(q.ToSQL())

	err = stmt.Get(result, args)

	return q.assignModel(result, q.Model.GetModel()), err
}

// Paginate .
func (q *Query) Paginate(page int, limit ...int) (map[string]interface{}, error) {
	defer q.resetBindings()

	if len(limit) > 0 {
		q.Take(limit[0])
	} else {
		q.Take(50)
	}

	q.Skip((page - 1) * q.Binding.Limit)

	data, err := q.Get()
	total := q.Count()

	result := map[string]interface{}{
		"data":  data,
		"total": total,
	}

	return result, err
}

// Insert .
func (q *Query) Insert(returning ...string) (interface{}, error) {
	var result *sqlx.Rows
	var err error

	query := q.Builder.BuildInsert(q.Model, returning...)

	q.Model.SetCreated()

	payload := q.Model.MapToPayload(q.Model)

	if nil != q.Tx {
		result, err = q.Tx.NamedQuery(query, payload)
	} else {
		result, err = q.DB.NamedQuery(query, payload)
	}

	if nil != result && result.Next() {
		result.StructScan(q.Model)
	}

	return q.Model, err
}

// Update .
func (q *Query) Update() (bool, error) {
	var err error

	query := q.Builder.BuildUpdate(q.Model)

	q.Model.SetUpdated()

	payload := q.Model.MapToPayload(q.Model)

	if nil != q.Tx {
		_, err = q.Tx.NamedQuery(query, payload)
	} else {
		_, err = q.DB.NamedQuery(query, payload)
	}

	if nil != err {
		return false, err
	}

	return true, nil
}

// Delete .
func (q *Query) Delete() (bool, error) {
	var err error

	query := q.Builder.BuildDelete(q.Model, nil)

	payload := q.Model.MapToPayload(q.Model)

	if q.Model.IsSoftDelete() {
		q.Model.SetDeleted()

		return q.Update()
	}

	if nil != q.Tx {
		_, err = q.Tx.NamedQuery(query, payload)
	} else {
		_, err = q.DB.NamedQuery(query, payload)
	}

	if nil != err {
		return false, err
	}

	return true, nil
}

// BulkInsert .
func (q *Query) BulkInsert(data interface{}, returning ...string) (bool, error) {
	var err error
	var value reflect.Value

	value = reflect.ValueOf(data)

	if reflect.Slice != value.Kind() {
		return false, errors.New("data must be a slice")
	}

	slice := make([]interface{}, value.Len())

	for i := 0; i < value.Len(); i++ {
		slice[i] = value.Index(i).Interface()
	}

	query := q.Builder.BuildBulkInsert(q.Model, slice, returning...)
	payloads := q.bulkPayload(slice)

	if nil != q.Tx {
		_, err = q.Tx.NamedQuery(query, payloads)
	} else {
		_, err = q.DB.NamedQuery(query, payloads)
	}

	if nil != err {
		return false, err
	}

	return true, nil
}

// RawCommand .
func (q *Query) RawCommand(dest IModel, query string, args interface{}) (interface{}, error) {
	var result *sqlx.Rows
	var err error

	if nil != q.Tx {
		result, err = q.Tx.NamedQuery(query, args)
	} else {
		result, err = q.DB.NamedQuery(query, args)
	}

	if nil != result && result.Next() {
		result.StructScan(dest)
	}

	return dest, err
}

// RawQuery .
func (q *Query) RawQuery(dest IModel, query string, args ...interface{}) error {
	var err error

	err = q.DB.Select(dest, query, args...)

	return err
}
