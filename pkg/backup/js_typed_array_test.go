package backup

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJsTypedArray_Write(t *testing.T) {
	t.Run("it works", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		writer := NewJsTypedArray(buf)

		_, err := writer.Write([]byte{1, 2})
		require.NoError(t, err)

		_, err = writer.Write([]byte{3, 4})
		require.NoError(t, err)

		require.Equal(t, []byte("1,2,3,4,"), buf.Bytes())
	})
}
