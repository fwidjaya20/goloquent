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
		"CREATE INDEX %s_indexes ON %s (%s);\n",
		blueprint.name,
		blueprint.name,
		strings.Join(blueprint.indexes, ","),
	)

	return query
}
