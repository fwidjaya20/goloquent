package goloquent

import (
	"fmt"
	"testing"

	"github.com/fwidjaya20/goloquent/pkg/goloquent"
	"github.com/stretchr/testify/require"
)

func TestBuilder_CreateTable(t *testing.T) {
	builder := goloquent.NewBuilder()

	schema := goloquent.Create("books", func(table *goloquent.Schema) {
		table.Numeric("id").AutoIncrement()
		table.String("name").NotNull()
		table.UUID("author_id")

		table.Timestamp()
		table.SoftDelete()

		table.Index("name", "author_id")

		table.Foreign("author_id").
			On("authors").
			Reference("id").
			OnUpdate(goloquent.RA_CASCADE).
			OnDelete(goloquent.RA_RESTRICT)
	})

	t.Run("CreateTable", func(t *testing.T) {
		query := builder.BuildCreateTable(schema)

		var expectedQuery string

		expectedQuery = fmt.Sprintf("%sCREATE TABLE IF NOT EXISTS books ( id SERIAL PRIMARY KEY UNIQUE NOT NULL,name VARCHAR NOT NULL,author_id UUID,created_at TIMESTAMP,updated_at TIMESTAMP,deleted_at TIMESTAMP, FOREIGN KEY (author_id) REFERENCES authors (id) ON UPDATE CASCADE ON DELETE RESTRICT );\n", expectedQuery)
		expectedQuery = fmt.Sprintf("%sCREATE INDEX IF NOT EXISTS books_indexes ON books (name,author_id);\n", expectedQuery)

		require.Equal(t, expectedQuery, query)
	})
}

func TestBuilder_AlterTable(t *testing.T) {
	builder := goloquent.NewBuilder()

	schema := goloquent.Table("authors", func(table *goloquent.Schema) {
		table.String("email").Unique()
		table.Date("born_at")
		table.String("address")

		table.Text("address").Change()

		table.Rename("born_at", "birthday")

		table.DropTimestamp()
		table.DropSoftDelete()
	})

	t.Run("AlterTable", func(t *testing.T) {
		query := builder.BuildAlterTable(schema)

		var expectedQuery string

		expectedQuery = fmt.Sprintf("%sALTER TABLE authors ADD COLUMN IF NOT EXISTS email VARCHAR UNIQUE,ADD COLUMN IF NOT EXISTS born_at DATE,ADD COLUMN IF NOT EXISTS address VARCHAR;\n", expectedQuery)
		expectedQuery = fmt.Sprintf("%sALTER TABLE authors ALTER COLUMN address DROP DEFAULT,ALTER COLUMN address TYPE TEXT USING address::TEXT;\n", expectedQuery)
		expectedQuery = fmt.Sprintf("%sALTER TABLE authors RENAME COLUMN born_at TO birthday;\n", expectedQuery)
		expectedQuery = fmt.Sprintf("%sALTER TABLE authors DROP COLUMN created_at,DROP COLUMN updated_at,DROP COLUMN deleted_at;", expectedQuery)

		require.Equal(t, expectedQuery, query)
	})
}

func TestBuilder_DropTable(t *testing.T) {
	builder := goloquent.NewBuilder()

	schema := goloquent.Drop("authors")

	t.Run("DropTable", func(t *testing.T) {
		query := builder.BuildDropTable(schema)

		var expectedQuery string

		expectedQuery = fmt.Sprintf("%sDROP TABLE IF EXISTS authors;", expectedQuery)

		require.Equal(t, expectedQuery, query)
	})
}
