package goloquent

import "fmt"

// Schema is a struct that is used to store information about schema table
type Schema struct {
	command     Command
	name        string
	columns     []*Column
	renames     []*Column
	drops       []string
	primaryKeys []string
	references  []*Reference
	uniques     []string
	indexes     []string
}

// Create is a command for create a table
func Create(name string, blueprint func(blueprint *Schema)) *Schema {
	schema := &Schema{
		command: CMD_CREATE,
		name:    name,
	}

	blueprint(schema)

	return schema
}

// Table is a command for alter a table
func Table(name string, blueprint func(blueprint *Schema)) *Schema {
	schema := &Schema{
		command: CMD_ALTER,
		name:    name,
	}

	blueprint(schema)

	return schema
}

// Drop is a command for drop a table
func Drop(name string) *Schema {
	return &Schema{
		command: CMD_DROP,
		name:    name,
	}
}

// SmallInteger is a Schema Command for create Schema Column
func (s *Schema) SmallInteger(name string) *Column {
	col := newColumn(name, DT_SMALLINT)

	return s.addColumn(col)
}

// Integer is a Schema Command for create Schema Column
func (s *Schema) Integer(name string) *Column {
	col := newColumn(name, DT_INT)

	return s.addColumn(col)
}

// BigInteger is a Schema Command for create Schema Column
func (s *Schema) BigInteger(name string) *Column {
	col := newColumn(name, DT_BIGINT)

	return s.addColumn(col)
}

// Decimal is a Schema Command for create Schema Column
func (s *Schema) Decimal(name string) *Column {
	col := newColumn(name, DT_DECIMAL)

	return s.addColumn(col)
}

// Numeric is a Schema Command for create Schema Column
func (s *Schema) Numeric(name string) *Column {
	col := newColumn(name, DT_NUMERIC)

	return s.addColumn(col)
}

// Real is a Schema Command for create Schema Column
func (s *Schema) Real(name string) *Column {
	col := newColumn(name, DT_REAL)

	return s.addColumn(col)
}

// Double is a Schema Command for create Schema Column
func (s *Schema) Double(name string) *Column {
	col := newColumn(name, DT_DOUBLE)

	return s.addColumn(col)
}

// SmallSerial is a Schema Command for create Schema Column
func (s *Schema) SmallSerial(name string) *Column {
	col := newColumn(name, DT_SMALLSERIAL)

	return s.addColumn(col)
}

// Serial is a Schema Command for create Schema Column
func (s *Schema) Serial(name string) *Column {
	col := newColumn(name, DT_SERIAL)

	return s.addColumn(col)
}

// BigSerial is a Schema Command for create Schema Column
func (s *Schema) BigSerial(name string) *Column {
	col := newColumn(name, DT_BIGSERIAL)

	return s.addColumn(col)
}

// String is a Schema Command for create Schema Column
func (s *Schema) String(name string) *Column {
	col := newColumn(name, DT_STRING)

	return s.addColumn(col)
}

// Text is a Schema Command for create Schema Column
func (s *Schema) Text(name string) *Column {
	col := newColumn(name, DT_TEXT)

	return s.addColumn(col)
}

// UUID is a Schema Command for create Schema Column
func (s *Schema) UUID(name string) *Column {
	col := newColumn(name, DT_UUID)

	return s.addColumn(col)
}

// JSON is a Schema Command for create Schema Column
func (s *Schema) JSON(name string) *Column {
	col := newColumn(name, DT_JSON)

	return s.addColumn(col)
}

// Boolean is a Schema Command for create Schema Column
func (s *Schema) Boolean(name string) *Column {
	col := newColumn(name, DT_BOOL)

	return s.addColumn(col)
}

// Date is a Schema Command for create Schema Column
func (s *Schema) Date(name string) *Column {
	col := newColumn(name, DT_DATE)

	return s.addColumn(col)
}

// Time is a Schema Command for create Schema Column
func (s *Schema) Time(name string) *Column {
	col := newColumn(name, DT_TIME)

	return s.addColumn(col)
}

// Primary is a Schema Command for add multiple primary keys
func (s *Schema) Primary(columns ...string) {
	s.primaryKeys = append(s.primaryKeys, columns...)
}

// Foreign is a Schema Command for add table references
func (s *Schema) Foreign(key string) *Reference {
	ref := &Reference{
		key:      key,
		onUpdate: RA_NO_ACTION,
		onDelete: RA_NO_ACTION,
	}

	return s.addReference(ref)
}

// Unique is a Schema Command for add multiple unique columns
func (s *Schema) Unique(columns ...string) {
	s.uniques = append(s.uniques, columns...)
}

// Index is a Schema Command for add multiple index columns
func (s *Schema) Index(columns ...string) {
	s.indexes = append(s.indexes, columns...)
}

// Timestamp is a Schema Command for create Schema Column
func (s *Schema) Timestamp() {
	s.addColumn(newColumn("created_at", DT_TIMESTAMP))
	s.addColumn(newColumn("updated_at", DT_TIMESTAMP))
}

// TimestampTz is a Schema Command for create Schema Column
func (s *Schema) TimestampTz() {
	s.addColumn(newColumn("created_at", DT_TIMESTAMPTZ))
	s.addColumn(newColumn("updated_at", DT_TIMESTAMPTZ))
}

// SoftDelete is a Schema Command for create Schema 'deleted_at' olumn
func (s *Schema) SoftDelete() {
	s.addColumn(newColumn("deleted_at", DT_TIMESTAMP))
}

// SoftDeleteTz is a Schema Command for create Schema 'deleted_at' olumn
func (s *Schema) SoftDeleteTz() {
	s.addColumn(newColumn("deleted_at", DT_TIMESTAMPTZ))
}

// Rename is a Schema Command for Renaming Schema Column
func (s *Schema) Rename(from string, to string) {
	col := renameColumn(from, to)

	s.addRenameColumn(col)
}

// Drop is a Schema Command for Drop Schema columns
func (s *Schema) Drop(column ...string) {
	s.drops = column
}

// DropTimestamp is a Schema Command for Drop Schema Timestamp
func (s *Schema) DropTimestamp() {
	s.drops = append(s.drops, "created_at", "updated_at")
}

// DropSoftDelete is a Schema Command for Drop Schema Timestamp
func (s *Schema) DropSoftDelete() {
	s.drops = append(s.drops, "deleted_at")
}

// Verbose is a function for print schema detail
func (s *Schema) Verbose() {
	fmt.Printf("Name        : %v\n", s.name)
	fmt.Printf("PrimaryKeys : %v\n", s.primaryKeys)
	fmt.Printf("Unique      : %v\n", s.uniques)
	fmt.Printf("Indexes     : %v\n", s.indexes)
	fmt.Printf("References  : [\n")
	for _, v := range s.references {
		fmt.Println("  {")
		v.Verbose()
		fmt.Println("  }")
	}
	fmt.Println("]")
	fmt.Println("----------------------------------------")
	for _, v := range s.columns {
		v.Verbose()
		fmt.Println("........................................")
	}
}

func (s *Schema) addColumn(col *Column) *Column {
	s.columns = append(s.columns, col)

	return col
}

func (s *Schema) addRenameColumn(col *Column) *Column {
	s.renames = append(s.renames, col)

	return col
}

func (s *Schema) addReference(ref *Reference) *Reference {
	s.references = append(s.references, ref)

	return ref
}
