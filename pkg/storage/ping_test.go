package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage_Ping(t *testing.T) {
	t.Run("it works", func(t *testing.T) {
		store, savepoint := MustNewStore(t)
		t.Cleanup(savepoint.Flush)

		err := store.Ping(context.Background())
		require.NoError(t, err)
	})
}
