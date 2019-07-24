package goloquent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAggregate_Count(t *testing.T) {
	aggregate := newAggregate(COUNT, "*")

	t.Run("TestAggregate_COUNT", func(t *testing.T) {
		require.Equal(t, aggregate, &Aggregate{
			AggregateFunc: COUNT,
			Column:        "*",
		})

		require.NotEqual(t, aggregate, &Aggregate{
			AggregateFunc: MAX,
			Column:        "*",
		})
	})
}

func TestAggregate_Max(t *testing.T) {
	aggregate := newAggregate(MAX, "price")

	t.Run("TestAggregate_MAX", func(t *testing.T) {
		require.Equal(t, aggregate, &Aggregate{
			AggregateFunc: MAX,
			Column:        "price",
		})

		require.NotEqual(t, aggregate, &Aggregate{
			AggregateFunc: MAX,
			Column:        "*",
		})
	})
}

func TestAggregate_Min(t *testing.T) {
	aggregate := newAggregate(MIN, "price")

	t.Run("TestAggregate_MAX", func(t *testing.T) {
		require.Equal(t, aggregate, &Aggregate{
			AggregateFunc: MIN,
			Column:        "price",
		})

		require.NotEqual(t, aggregate, &Aggregate{
			AggregateFunc: MIN,
			Column:        "*",
		})
	})
}

func TestAggregate_Avg(t *testing.T) {
	aggregate := newAggregate(AVG, "price")

	t.Run("TestAggregate_MAX", func(t *testing.T) {
		require.Equal(t, aggregate, &Aggregate{
			AggregateFunc: AVG,
			Column:        "price",
		})

		require.NotEqual(t, aggregate, &Aggregate{
			AggregateFunc: AVG,
			Column:        "*",
		})
	})
}

func TestAggregate_Sum(t *testing.T) {
	aggregate := newAggregate(SUM, "price")

	t.Run("TestAggregate_SUM", func(t *testing.T) {
		require.Equal(t, aggregate, &Aggregate{
			AggregateFunc: SUM,
			Column:        "price",
		})

		require.NotEqual(t, aggregate, &Aggregate{
			AggregateFunc: SUM,
			Column:        "*",
		})
	})
}
