package goloquent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReference_RelatedKey(t *testing.T) {
	ref := Reference{}

	ref.Reference("id").On("roles").OnUpdate(RA_CASCADE).OnDelete(RA_RESTRICT)

	t.Run("TestReference_RELATED_KEY", func(t *testing.T) {
		require.Equal(t, ref, Reference{
			relatedKey:   "id",
			relatedTable: "roles",
			onUpdate:     "CASCADE",
			onDelete:     "RESTRICT",
		})

		require.NotEqual(t, ref.relatedKey, "role_id")
	})
}

func TestReference_RelatedTable(t *testing.T) {
	ref := Reference{}

	ref.Reference("id").On("roles").OnUpdate(RA_CASCADE).OnDelete(RA_RESTRICT)

	t.Run("TestReference_RELATED_TABLE", func(t *testing.T) {
		require.Equal(t, ref, Reference{
			relatedKey:   "id",
			relatedTable: "roles",
			onUpdate:     "CASCADE",
			onDelete:     "RESTRICT",
		})

		require.NotEqual(t, ref.relatedTable, "divisions")
	})
}

func TestReference_OnUpdate(t *testing.T) {
	ref := Reference{}

	ref.Reference("id").On("roles").OnUpdate(RA_CASCADE).OnDelete(RA_RESTRICT)

	t.Run("TestReference_ON_UPDATE", func(t *testing.T) {
		require.Equal(t, ref, Reference{
			relatedKey:   "id",
			relatedTable: "roles",
			onUpdate:     "CASCADE",
			onDelete:     "RESTRICT",
		})

		require.NotEqual(t, ref.onUpdate, RA_NO_ACTION)
	})
}

func TestReference_OnDelete(t *testing.T) {
	ref := Reference{}

	ref.Reference("id").On("roles").OnUpdate(RA_CASCADE).OnDelete(RA_RESTRICT)

	t.Run("TestReference_ON_DELETE", func(t *testing.T) {
		require.Equal(t, ref, Reference{
			relatedKey:   "id",
			relatedTable: "roles",
			onUpdate:     "CASCADE",
			onDelete:     "RESTRICT",
		})

		require.NotEqual(t, ref.onDelete, RA_CASCADE)
	})
}
