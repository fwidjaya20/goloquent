package goloquent

import "fmt"

// Column is a struct that is used for store information about table column
type Column struct {
	oldName      string // for renaming column
	name         string
	dataType     DataType
	defaultValue interface{}
	primaryKey   bool
	nullable     bool
	unique       bool
	modified     bool
}

func newColumn(name string, dt DataType) *Column {
	return &Column{
		name:         name,
		dataType:     dt,
		defaultValue: nil,
		primaryKey:   false,
		unique:       false,
		nullable:     true,
		modified:     false,
	}
}

func renameColumn(from string, to string) *Column {
	return &Column{
		oldName: from,
		name:    to,
	}
}

// PrimaryKey is a setter for primary key
func (c *Column) PrimaryKey() *Column {
	c.primaryKey = true

	return c
}

// NotNull is a setter for nullable
func (c *Column) NotNull() *Column {
	c.nullable = false

	return c
}

// Unique is a setter for unique
func (c *Column) Unique() *Column {
	c.unique = true

	return c
}

// AutoIncrement is a setter for auto increment
func (c *Column) AutoIncrement() *Column {
	c.dataType = DT_SERIAL

	c.PrimaryKey()
	c.NotNull()
	c.Unique()

	return c
}

// DefaultValue is a setter for default value
func (c *Column) DefaultValue(v interface{}) *Column {
	c.defaultValue = v

	return c
}

// Change is a setter for modify column
func (c *Column) Change() *Column {
	c.modified = true

	return c
}

// Verbose is a function for print column detail
func (c *Column) Verbose() {
	fmt.Printf("  Column       : %v\n", c.name)
	fmt.Printf("  DataType     : %v\n", c.dataType)
	fmt.Printf("  PK           : %v\n", c.primaryKey)
	fmt.Printf("  Nullable     : %v\n", c.nullable)
	fmt.Printf("  Unique       : %v\n", c.unique)
	fmt.Printf("  Defaultvalue : %v\n", c.defaultValue)
}
