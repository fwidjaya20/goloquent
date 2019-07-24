package goloquent

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCondition_Connector(t *testing.T) {
	condition := newCondition(AND, "id", EQUAL, 1)

	t.Run("TestCondition_CONNECTOR", func(t *testing.T) {
		require.Equal(t, condition, &Condition{
			Connector: "AND",
			Column:    "id",
			Operator:  "=",
			Value:     1,
		})

		require.NotEqual(t, condition.Connector, OR)
	})
}

func TestCondition_Column(t *testing.T) {
	condition := newCondition(AND, "name", ILIKE, "%bulk%")

	t.Run("TestCondition_COLUMN", func(t *testing.T) {
		require.Equal(t, condition, &Condition{
			Connector: "AND",
			Column:    "name",
			Operator:  "ILIKE",
			Value:     "%bulk%",
		})

		require.NotEqual(t, condition.Column, "id")
	})
}

func TestCondition_Operator(t *testing.T) {
	condition := newCondition(AND, "name", LIKE, "%bulk%")

	t.Run("TestCondition_OPERATOR", func(t *testing.T) {
		require.Equal(t, condition, &Condition{
			Connector: "AND",
			Column:    "name",
			Operator:  "LIKE",
			Value:     "%bulk%",
		})

		require.NotEqual(t, condition.Operator, ILIKE)
	})
}

func TestCondition_Value(t *testing.T) {
	condition := newCondition(AND, "name", LIKE, "%bulk%")

	t.Run("TestCondition_VALUE", func(t *testing.T) {
		require.Equal(t, condition, &Condition{
			Connector: "AND",
			Column:    "name",
			Operator:  "LIKE",
			Value:     "%bulk%",
		})

		require.NotEqual(t, condition.Value, "bulk%")
	})
}
