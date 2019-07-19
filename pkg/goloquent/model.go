package goloquent

import (
	"reflect"
	"time"
)

// IModel .
type IModel interface {
	GetTableName() string
	GetPK() string
	GetColumns(v IModel) []string
	IsAutoIncrement() bool
	IsUuid() bool
	IsTimestamp() bool
	IsSoftDelete() bool
	MapToPayload(v IModel) map[string]interface{}
}

// Model .
type Model struct {
	Table         string
	PrimaryKey    string
	AutoIncrement bool
	Uuid          bool
	Timestamp     bool
	SoftDelete    bool
	CreatedAt     *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deleted_at"`
}

// AutoIncrementModel is a factory method for creating Model with AutoIncrement TRUE
func AutoIncrementModel(table string, pk string, isTimestamp bool, isSoftDelete bool) Model {
	return Model{
		Table:         table,
		PrimaryKey:    pk,
		AutoIncrement: true,
		Uuid:          false,
		Timestamp:     isTimestamp,
		SoftDelete:    isSoftDelete,
	}
}

// GetTableName .
func (m *Model) GetTableName() string {
	return m.Table
}

// GetPK .
func (m *Model) GetPK() string {
	return m.PrimaryKey
}

// GetColumns .
func (m *Model) GetColumns(v IModel) []string {
	var columns []string

	typeOf := reflect.TypeOf(v)

	for i := 0; i < typeOf.Elem().NumField(); i++ {
		column := typeOf.Elem().Field(i)

		if "Model" != column.Name {
			tag := typeOf.Elem().Field(i).Tag.Get("db")

			columns = append(columns, tag)
		}
	}

	return columns
}

// IsAutoIncrement .
func (m *Model) IsAutoIncrement() bool {
	return m.AutoIncrement
}

// IsUuid .
func (m *Model) IsUuid() bool {
	return m.Uuid
}

// IsTimestamp .
func (m *Model) IsTimestamp() bool {
	return m.Timestamp
}

// IsSoftDelete .
func (m *Model) IsSoftDelete() bool {
	return m.SoftDelete
}

// MapToPayload .
func (m *Model) MapToPayload(v IModel) map[string]interface{} {
	var payload = make(map[string]interface{})

	model := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	if reflect.Ptr == model.Kind() {
		model = model.Elem()
	}

	if reflect.Ptr == value.Kind() {
		value = value.Elem()
	}

	for i := 0; i < model.NumField(); i++ {
		if "Model" != model.Field(i).Name {
			key := model.Field(i).Tag.Get("db")
			val := value.Field(i).Interface()

			payload[key] = val
		} else {
			model := value.Field(i).Interface().(Model)
			if model.IsTimestamp() {
				payload[CREATED_AT] = model.CreatedAt
				payload[UPDATED_AT] = model.UpdatedAt
			}
			if model.IsSoftDelete() {
				payload[DELETED_AT] = model.DeletedAt
			}
		}
	}

	return payload
}
