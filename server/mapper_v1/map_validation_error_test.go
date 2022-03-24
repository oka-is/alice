package mapper_v1

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wault-pw/alice/pkg/validator"
)

func TestMapValidationError(t *testing.T) {
	t.Run("is empty for NULL", func(t *testing.T) {
		got := MapValidationError(nil)
		require.Nil(t, got.Items)
	})

	t.Run("it works", func(t *testing.T) {
		err := validator.NewError("foo")
		got := MapValidationError(err)
		require.Len(t, got.Items, 1)
		require.Equal(t, "BASE", got.GetItems()[0].GetField())
		require.Equal(t, "foo", got.GetItems()[0].GetDescription())
	})
}
