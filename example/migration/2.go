package migration

import "github.com/fwidjaya20/goloquent/pkg/goloquent"

// Migration2 .
var Migration2 = goloquent.Migration{
	Schema: []*goloquent.Schema{
		goloquent.Create("movie_ratings", func(table *goloquent.Schema) {
			table.UUID("id").PrimaryKey()
			table.Serial("movie_id")
			table.SmallInteger("rate")

			table.Foreign("movie_id").
				Reference("id").
				On("movies").
				OnUpdate(goloquent.RA_CASCADE).
				OnDelete(goloquent.RA_RESTRICT)

			table.Index("id", "movie_id", "rate")
		}),
	},
}
