package goloquent

import (
	"fmt"
	"reflect"
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

// BuildSelect .
func (b *Builder) BuildSelect(model IModel, binding Binding) string {
	var query string

	query = fmt.Sprintf("%sSELECT", query)

	if nil != binding.Aggregate {
		query = fmt.Sprintf("%s %s", query, b.buildSelectAggregate(binding.Aggregate.AggregateFunc, model.GetTableName(), binding.Aggregate.Column))
	} else if len(model.GetColumns(model)) > 0 {
		var selectColumns string

		selectComma := true
		for _, col := range model.GetColumns(model) {
			selectColumns, selectComma = b.buildSelectColumns(selectColumns, model.GetTableName(), col, selectComma)
		}

		query = fmt.Sprintf("%s %s", query, selectColumns)
	}

	query = fmt.Sprintf(`%sFROM "%s" `, query, model.GetTableName())

	if len(binding.Conditions) > 0 {
		query = fmt.Sprintf(`%sWHERE %s `, query, b.buildQueryCondition(model.GetTableName(), binding))
	}

	if binding.Limit > 0 {
		query = fmt.Sprintf(`%sLIMIT %d `, query, binding.Limit)
	}

	if binding.Offset > 0 {
		query = fmt.Sprintf(`%sOFFSET %d;`, query, binding.Offset)
	}

	return query
}

// BuildInsert .
func (b *Builder) BuildInsert(model IModel, returning ...string) string {
	var query string

	query = fmt.Sprintf("%sINSERT INTO %s ", query, model.GetTableName())
	query = fmt.Sprintf("%s(%s) ", query, b.buildInsertColumnOrValue(model, b.isAutoIncrementPrimaryKey, b.buildInsertColumns))
	query = fmt.Sprintf("%sVALUES (%s) ", query, b.buildInsertColumnOrValue(model, b.isAutoIncrementPrimaryKey, b.buildInsertValue))

	if len(returning) < 1 {
		returning = append(returning, model.GetPK())
	}

	query = fmt.Sprintf("%sRETURNING \"%s\";\n", query, strings.Join(returning, `", "`))

	return query
}

// BuildBulkInsert .
func (b *Builder) BuildBulkInsert(model IModel, data []interface{}, returning ...string) string {
	var query string

	query = fmt.Sprintf("%sINSERT INTO %s ", query, model.GetTableName())
	query = fmt.Sprintf("%s(%s) ", query, b.buildInsertColumnOrValue(model, b.isAutoIncrementPrimaryKey, b.buildInsertColumns))
	query = fmt.Sprintf("%sVALUES ", query)

	first := true
	for i, v := range data {
		if first {
			first = false
		} else {
			query = fmt.Sprintf("%s, ", query)
		}

		query = b.buildInsertBulkValue(query, v.(IModel), i)
	}

	if len(returning) < 1 {
		returning = append(returning, model.GetPK())
	}

	query = fmt.Sprintf("%s RETURNING \"%s\";\n", query, strings.Join(returning, `", "`))

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

// isAutoIncrementPrimaryKey is a function that will skip id column if model is autoincrement when building query
func (b *Builder) isAutoIncrementPrimaryKey(column string, model IModel) bool {
	return column == model.GetPK() && !model.IsUuid() && model.IsAutoIncrement()
}

// buildInsertColumnOrValue is a decorator function to wrap creational of columns or values
func (b *Builder) buildInsertColumnOrValue(
	model IModel,
	skipFunc func(column string, model IModel) bool,
	generator func(query string, column string, hasComma bool) (string, bool),
) string {
	var query string
	hasComma := true

	columns := model.GetColumns(model)

	if model.IsTimestamp() {
		columns = append(columns, "created_at", "updated_at")
	}

	if model.IsSoftDelete() {
		columns = append(columns, "deleted_at")
	}

	for _, col := range columns {
		if !skipFunc(col, model) {
			query, hasComma = generator(query, col, hasComma)
		}
	}

	return query
}

func (b *Builder) buildInsertColumns(query string, column string, hasComma bool) (string, bool) {
	if hasComma {
		hasComma = false
	} else {
		query = fmt.Sprintf("%s, ", query)
	}

	return fmt.Sprintf(`%s"%s"`, query, column), hasComma
}

func (b *Builder) buildInsertValue(query string, column string, hasComma bool) (string, bool) {
	if hasComma {
		hasComma = false
	} else {
		query = fmt.Sprintf("%s, ", query)
	}

	return fmt.Sprintf(`%s:%s`, query, column), hasComma
}

func (b *Builder) buildInsertBulkValue(query string, model IModel, i int) string {
	columns := model.GetColumns(model)

	if model.IsTimestamp() {
		columns = append(columns, "created_at", "updated_at")
	}

	if model.IsSoftDelete() {
		columns = append(columns, "deleted_at")
	}

	var qCol []string
	for _, col := range columns {
		if !b.isAutoIncrementPrimaryKey(col, model) {
			qCol = append(qCol, fmt.Sprintf(":%d%s", i, col))
		}
	}

	return fmt.Sprintf(`%s(%s)`, query, strings.Join(qCol, ", "))
}

func (b *Builder) buildSelectColumns(query string, table string, column string, hasComma bool) (string, bool) {
	if hasComma {
		hasComma = false
	} else {
		query = fmt.Sprintf("%s, ", query)
	}

	query = fmt.Sprintf(`%s"%s"."%s"`, query, table, column)

	return query, hasComma
}

func (b *Builder) buildSelectAggregate(aggregateFn AggregateFunction, table string, column string) string {
	if "*" == column {
		return fmt.Sprintf(`%v("%s".%s) `, aggregateFn, table, column)
	}

	return fmt.Sprintf(`%v("%s"."%s") `, aggregateFn, table, column)
}

func (b *Builder) buildQueryCondition(table string, binding Binding) string {
	var query string

	for i, w := range binding.Conditions {
		switch w.Operator {
		case IN, NOT_IN:
			if 0 == i {
				query = fmt.Sprintf(`%s"%s"."%s" %s (:%d%s)`, query, table, w.Column, w.Operator, i, w.Column)
			} else {
				query = fmt.Sprintf(`%s%s "%s"."%s" %s (:%d%s)`, query, w.Connector, table, w.Column, w.Operator, i, w.Column)
			}
		default:
			if 0 == i {
				query = fmt.Sprintf(`%s"%s"."%s" %s :%d%s `, query, table, w.Column, w.Operator, i, w.Column)
			} else {
				query = fmt.Sprintf(`%s%s "%s"."%s" %s :%d%s `, query, w.Connector, table, w.Column, w.Operator, i, w.Column)
			}
		}
	}

	return query
}

func (b *Builder) buildWhereInValues(value interface{}) string {
	var query string

	items := reflect.ValueOf(value)

	if reflect.Slice == items.Kind() {
		for i := 0; i < items.Len(); i++ {
			if items.Len()-1 == i {
				query = fmt.Sprintf(`%s'%v'`, query, items.Index(i))
			} else {
				query = fmt.Sprintf(`%s'%v',`, query, items.Index(i))
			}
		}
	}

	return query
}
