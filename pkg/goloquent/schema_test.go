package goloquent

import (
	"reflect"
	"testing"
)

func TestCreateSchema(t *testing.T) {
	tt := []struct {
		Schema   *Schema
		Expected *Schema
	}{
		{
			Schema: Create("genres", func(table *Schema) {
				table.Serial("id").AutoIncrement()
				table.String("name").Unique()
				table.Timestamp()
				table.Index("id", "name")
			}),
			Expected: &Schema{
				columns: []*Column{
					&Column{
						name:       "id",
						dataType:   DT_SERIAL,
						primaryKey: true,
						unique:     true,
					},
					&Column{
						name:     "name",
						dataType: DT_STRING,
						unique:   true,
						nullable: true,
					}, &Column{
						name:     "created_at",
						dataType: DT_TIMESTAMP,
						nullable: true,
					},
					&Column{
						name:     "updated_at",
						dataType: DT_TIMESTAMP,
						nullable: true,
					},
				},
				command: CMD_CREATE,
				indexes: []string{
					"id",
					"name",
				},
				name: "genres",
			},
		},
		{
			Schema: Create("movies", func(table *Schema) {
				table.Serial("id").AutoIncrement()
				table.String("title")
				table.SmallInteger("year")
				table.Serial("genre_id")
				table.SmallInteger("duration")
				table.String("director")
				table.Foreign("genre_id").
					Reference("id").
					On("genres").
					OnUpdate(RA_CASCADE).
					OnDelete(RA_RESTRICT)
				table.Index("id", "title", "year", "director")
			}),
			Expected: &Schema{
				columns: []*Column{
					&Column{
						name:       "id",
						dataType:   DT_SERIAL,
						primaryKey: true,
						unique:     true,
					},
					&Column{
						name:     "title",
						dataType: DT_STRING,
						nullable: true,
					}, &Column{
						name:     "year",
						dataType: DT_SMALLINT,
						nullable: true,
					},
					&Column{
						name:     "genre_id",
						dataType: DT_SERIAL,
						nullable: true,
					},
					&Column{
						name:     "duration",
						dataType: DT_SMALLINT,
						nullable: true,
					},
					&Column{
						name:     "director",
						dataType: DT_STRING,
						nullable: true,
					},
				},
				command: CMD_CREATE,
				indexes: []string{
					"id",
					"title",
					"year",
					"director",
				},
				name: "movies",
				references: []*Reference{
					&Reference{
						onUpdate:     RA_CASCADE,
						onDelete:     RA_RESTRICT,
						key:          "genre_id",
						relatedKey:   "id",
						relatedTable: "genres",
					},
				},
			},
		},
		{
			Schema: Create("movie_ratings", func(table *Schema) {
				table.UUID("id").PrimaryKey()
				table.Serial("movie_id")
				table.SmallInteger("rate")

				table.Timestamp()

				table.Foreign("movie_id").
					Reference("id").
					On("movies").
					OnUpdate(RA_CASCADE).
					OnDelete(RA_RESTRICT)

				table.Index("id", "movie_id", "rate")
			}),
			Expected: &Schema{
				columns: []*Column{
					&Column{
						name:       "id",
						dataType:   DT_UUID,
						primaryKey: true,
						nullable:   true,
					},
					&Column{
						name:     "movie_id",
						dataType: DT_SERIAL,
						nullable: true,
					}, &Column{
						name:     "rate",
						dataType: DT_SMALLINT,
						nullable: true,
					},
					&Column{
						name:     "created_at",
						dataType: DT_TIMESTAMP,
						nullable: true,
					},
					&Column{
						name:     "updated_at",
						dataType: DT_TIMESTAMP,
						nullable: true,
					},
				},
				command: CMD_CREATE,
				indexes: []string{
					"id",
					"movie_id",
					"rate",
				},
				name: "movie_ratings",
				references: []*Reference{
					&Reference{
						onUpdate:     RA_CASCADE,
						onDelete:     RA_RESTRICT,
						key:          "movie_id",
						relatedKey:   "id",
						relatedTable: "movies",
					},
				},
			},
		},

		{
			Schema: Table("movie_ratings", func(table *Schema) {
				table.Text("hello").DefaultValue("-")
				table.Text("world").Unique()
				table.Numeric("sekali").DefaultValue(1000)
				table.Text("lagi")
				table.Text("ini")
				table.Text("testing")

				table.String("testing").Change()
				table.Date("lagi").Change()

				table.Rename("testing", "testing_testing")

				table.Drop("hello", "world", "sekali", "lagi", "ini")

				table.SoftDeleteTz()
			}),
			Expected: &Schema{
				columns: []*Column{
					&Column{
						name:         "hello",
						dataType:     DT_TEXT,
						defaultValue: "-",
						nullable:     true,
					},
					&Column{
						name:     "world",
						dataType: DT_TEXT,
						unique:   true,
						nullable: true,
					}, &Column{
						name:         "sekali",
						dataType:     DT_NUMERIC,
						defaultValue: 1000,
						nullable:     true,
					},
					&Column{
						name:     "lagi",
						dataType: DT_TEXT,
						nullable: true,
					}, &Column{
						name:     "ini",
						dataType: DT_TEXT,
						nullable: true,
					},
					&Column{
						name:     "testing",
						dataType: DT_TEXT,
						nullable: true,
					},
					&Column{
						name:     "testing",
						dataType: DT_STRING,
						nullable: true,
						modified: true,
					},
					&Column{
						name:     "lagi",
						dataType: DT_DATE,
						nullable: true,
						modified: true,
					},
					&Column{
						name:     "deleted_at",
						dataType: DT_TIMESTAMPTZ,
						nullable: true,
					},
				},
				command: CMD_ALTER,
				drops: []string{
					"hello",
					"world",
					"sekali",
					"lagi",
					"ini",
				},
				renames: []*Column{
					&Column{
						name:    "testing_testing",
						oldName: "testing",
					},
				},
				name: "movie_ratings",
			},
		},
		{
			Schema: Drop("movie_ratings"),
			Expected: &Schema{
				name:    "movie_ratings",
				command: CMD_DROP,
			},
		},
	}

	for _, v := range tt {
		if !reflect.DeepEqual(*v.Expected, *v.Schema) {
			for i := range v.Schema.columns {
				if !reflect.DeepEqual(*v.Expected.columns[i], *v.Schema.columns[i]) {
					t.Errorf("expected column: %v", *v.Expected.columns[i])
					t.Errorf("schema column: %v", *v.Schema.columns[i])
				}
			}
			for i := range v.Schema.references {
				if !reflect.DeepEqual(*v.Expected.references[i], *v.Schema.references[i]) {
					t.Errorf("expected references: %v", *v.Expected.references[i])
					t.Errorf("schema references: %v", *v.Schema.references[i])
				}
			}
			for i := range v.Schema.renames {
				if !reflect.DeepEqual(*v.Expected.renames[i], *v.Schema.renames[i]) {
					t.Errorf("expected renames column: %v", *v.Expected.renames[i])
					t.Errorf("schema renames column: %v", *v.Schema.renames[i])
				}
			}
			t.Errorf("expected schema %v", v.Expected)
			t.Errorf("got schema %v", v.Expected)
		}
	}
}
