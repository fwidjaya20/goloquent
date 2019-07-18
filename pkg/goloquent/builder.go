package goloquent

import (
	"fmt"
	"strings"
	"sync"
)

var once sync.Once
var instance *Builder

// Builder .
type Builder struct{}

// NewBuilder .
func NewBuilder() *Builder {
	once.Do(func() {
		instance = &Builder{}
	})

	return instance
}

// BuildCreateTable .
func (b *Builder) BuildCreateTable(blueprint *Schema) string {
	var query string

	query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", blueprint.name)
	query = fmt.Sprintf("%s (", query)

	var columns []string
	for _, col := range blueprint.columns {
		columns = append(columns, b.buildColumnQuery(col))
	}
	query = fmt.Sprintf("%s %s", query, strings.Join(columns, ","))

	if len(blueprint.primaryKeys) > 0 {
		query = fmt.Sprintf("%s, PRIMARY KEY(%s)", query, strings.Join(blueprint.primaryKeys, ","))
	}

	if len(blueprint.uniques) > 0 {
		query = fmt.Sprintf("%s, UNIQUE(%s)", query, strings.Join(blueprint.uniques, ","))
	}

	if len(blueprint.references) > 0 {
		var references []string
		for _, ref := range blueprint.references {
			references = append(references, b.buildReferenceQuery(ref))
		}
		query = fmt.Sprintf("%s, %s", query, strings.Join(references, ","))
	}

	query = fmt.Sprintf("%s );\n", query)

	if len(blueprint.indexes) > 0 {
		query = fmt.Sprintf("%s%s", query, b.buildIndexQuery(blueprint))
	}

	return query
}

// BuildAlterTable .
func (b *Builder) BuildAlterTable(schema *Schema) string {
	var query []string

	if len(schema.columns) > 0 {
		addColumn := fmt.Sprintf("ALTER TABLE %s %s;", schema.name, b.buildAddColumnQuery(schema.columns...))

		query = append(query, addColumn)
	}

	modifyColumn := b.buildModifyColumnQuery(schema.columns...)
	if "" != modifyColumn {
		modifyColumn = fmt.Sprintf("ALTER TABLE %s %s;", schema.name, modifyColumn)

		query = append(query, modifyColumn)
	}

	if len(schema.renames) > 0 {
		var renameColumn string

		renameColumn = fmt.Sprintf("ALTER TABLE %s %s;", schema.name, b.buildRenameColumnQuery(schema.renames...))

		query = append(query, renameColumn)
	}

	if len(schema.drops) > 0 {
		dropColumn := fmt.Sprintf("ALTER TABLE %s %s;", schema.name, b.buildDropColumnQuery(schema.drops...))

		query = append(query, dropColumn)
	}

	return strings.Join(query, "\n")
}

// BuildDropTable .
func (b *Builder) BuildDropTable(schema *Schema) string {
	var query string

	query = fmt.Sprintf("DROP TABLE IF EXISTS %s;", schema.name)

	fmt.Println(query)

	return query
}

func (b *Builder) buildColumnQuery(column *Column) string {
	var query string

	query = fmt.Sprintf("%s %s", column.name, column.dataType)

	if column.primaryKey {
		query = fmt.Sprintf("%s PRIMARY KEY", query)
	}

	if column.unique {
		query = fmt.Sprintf("%s UNIQUE", query)
	}

	if !column.nullable {
		query = fmt.Sprintf("%s NOT NULL", query)
	}

	if nil != column.defaultValue {
		query = fmt.Sprintf("%s DEFAULT '%v'", query, column.defaultValue)
	}

	return query
}

func (b *Builder) buildReferenceQuery(ref *Reference) string {
	var query string

	query = fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s (%s)", ref.key, ref.relatedTable, ref.relatedKey)

	if "" != ref.onUpdate {
		query = fmt.Sprintf("%s ON UPDATE %s", query, ref.onUpdate)
	}

	if "" != ref.onDelete {
		query = fmt.Sprintf("%s ON DELETE %s", query, ref.onDelete)
	}

	return query
}

func (b *Builder) buildIndexQuery(blueprint *Schema) string {
	var query string

	query = fmt.Sprintf(
		"CREATE INDEX IF NOT EXISTS %s_indexes ON %s (%s);\n",
		blueprint.name,
		blueprint.name,
		strings.Join(blueprint.indexes, ","),
	)

	return query
}

func (b *Builder) buildAddColumnQuery(columns ...*Column) string {
	var query []string

	for _, col := range columns {
		if !col.modified {
			query = append(query, fmt.Sprintf("ADD COLUMN IF NOT EXISTS %s", b.buildColumnQuery(col)))
		}
	}

	return strings.Join(query, ",")
}

func (b *Builder) buildModifyColumnQuery(columns ...*Column) string {
	var query []string

	for _, col := range columns {
		if col.modified {
			query = append(query, fmt.Sprintf("ALTER COLUMN %s DROP DEFAULT", col.name))
			query = append(query, fmt.Sprintf("ALTER COLUMN %s TYPE %s USING %s::%s", col.name, col.dataType, col.name, col.dataType))
		}
	}

	return strings.Join(query, ",")
}

func (b *Builder) buildRenameColumnQuery(columns ...*Column) string {
	var query []string

	for _, col := range columns {
		query = append(query, fmt.Sprintf("RENAME COLUMN %s TO %s", col.oldName, col.name))
	}

	return strings.Join(query, ",")
}

func (b *Builder) buildDropColumnQuery(columns ...string) string {
	var query []string

	for _, col := range columns {
		query = append(query, fmt.Sprintf("DROP COLUMN %s", col))
	}

	return strings.Join(query, ",")
}
