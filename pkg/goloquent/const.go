package goloquent

// DataType is a replica of string type that used for store postgres data types
type DataType string

// Command is a replica of string type that used for store postgres command
type Command string

// ReferenceAction is a replica of string type that used for store relationship action
type ReferenceAction string

// Condition is a replica of string type that used for store condition operator
type Operator string

// Connector is a replica of string type that used for connecting between condition
type Connector string

// AggregateFunction is a replica of string type that used for store Aggregate Function
type AggregateFunction string

const (
	CMD_CREATE Command = "CREATE"
	CMD_ALTER  Command = "ALTER"
	CMD_DROP   Command = "DROP"
	CMD_ADD    Command = "ADD"
	CMD_RENAME Command = "RENAME"
)

const (
	DT_INT         DataType = "INTEGER"
	DT_SMALLINT    DataType = "SMALLINT"
	DT_BIGINT      DataType = "BIGINT"
	DT_DECIMAL     DataType = "DECIMAL"
	DT_NUMERIC     DataType = "NUMERIC"
	DT_REAL        DataType = "REAL"
	DT_DOUBLE      DataType = "DOUBLE"
	DT_SMALLSERIAL DataType = "SMALLSERIAL"
	DT_SERIAL      DataType = "SERIAL"
	DT_BIGSERIAL   DataType = "BIGSERIAL"
	DT_STRING      DataType = "VARCHAR"
	DT_TEXT        DataType = "TEXT"
	DT_UUID        DataType = "UUID"
	DT_JSON        DataType = "JSON"
	DT_BOOL        DataType = "BOOLEAN"
	DT_DATE        DataType = "DATE"
	DT_TIME        DataType = "TIME"
	DT_TIMESTAMP   DataType = "TIMESTAMP"
	DT_TIMESTAMPTZ DataType = "TIMESTAMPTZ"
)

const (
	RA_NO_ACTION ReferenceAction = "NO ACTION"
	RA_RESTRICT  ReferenceAction = "RESTRICT"
	RA_CASCADE   ReferenceAction = "CASCADE"
)

const (
	CREATED_AT = "created_at"
	UPDATED_AT = "updated_at"
	DELETED_AT = "deleted_at"
)

const (
	AND Connector = "AND"
	OR  Connector = "OR"
)

const (
	EQUAL                 Operator = "="
	NOT_EQUAL             Operator = "!="
	LESS_THAN             Operator = "<"
	GREATER_THAN          Operator = ">"
	LESS_THAN_OR_EQUAL    Operator = "<="
	GREATER_THAN_OR_EQUAL Operator = ">="
	LIKE                  Operator = "LIKE"
	ILIKE                 Operator = "ILIKE"
	NOT_LIKE              Operator = "NOT LIKE"
	IN                    Operator = "IN"
	NOT_IN                Operator = "NOT IN"
)

const (
	COUNT AggregateFunction = "COUNT"
	MIN   AggregateFunction = "MIN"
	MAX   AggregateFunction = "MAX"
	AVG   AggregateFunction = "AVG"
	SUM   AggregateFunction = "SUM"
)
