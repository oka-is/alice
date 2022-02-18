package backup

import (
	"context"
	"io"

	"github.com/oka-is/alice/pkg/storage"
)

type Backup struct {
	writer *Writer
	store  storage.IStore
	f      IFlush
	ctx    context.Context
}

func NewJsBackup(store storage.IStore, wr io.Writer, flusher IFlush) *Backup {
	return NewBackup(store, NewJsTypedArray(wr), flusher)
}

func NewBackup(store storage.IStore, wr io.Writer, flusher IFlush) *Backup {
	return &Backup{
		ctx:    context.Background(),
		writer: NewWriter(wr),
		store:  store,
		f:      flusher,
	}
}
