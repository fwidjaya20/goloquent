package migration

import "github.com/fwidjaya20/goloquent/pkg/goloquent"

// Migration1 .
var Migration1 = goloquent.Migration{
	Schema: []*goloquent.Schema{
		goloquent.Create("genres", func(table *goloquent.Schema) {
			table.Serial("id").AutoIncrement()
			table.String("name").Unique()
			table.Timestamp()
			table.Index("id", "name")
		}),
		goloquent.Create("movies", func(table *goloquent.Schema) {
			table.Serial("id").AutoIncrement()
			table.String("title")
			table.SmallInteger("year")
			table.Serial("genre_id")
			table.SmallInteger("duration")
			table.String("director")

			table.Foreign("genre_id").
				Reference("id").
				On("genres").
				OnUpdate(goloquent.RA_CASCADE).
				OnDelete(goloquent.RA_RESTRICT)

			table.Index("id", "title", "year", "director")
		}),
	},
}
