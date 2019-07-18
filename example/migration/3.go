package migration

import "github.com/fwidjaya20/goloquent/pkg/goloquent"

var Migration3 = goloquent.Migration{
	Schema: []*goloquent.Schema{
		goloquent.Table("movie_ratings", func(table *goloquent.Schema) {
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
		goloquent.Drop("movie_ratings"),
	},
}
