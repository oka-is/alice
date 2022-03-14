package cmd

import (
	"context"

	"github.com/urfave/cli/v2"
	"github.com/wault-pw/alice/pkg/storage"
)

type key int

const (
	keyStore = key(iota)
)

type Context struct {
	*cli.Context
}

func Ctx(ctx *cli.Context) *Context {
	return &Context{ctx}
}

func (c *Context) SetStore(store *storage.Storage) {
	c.Context.Context = context.WithValue(c.Context.Context, keyStore, store)
}

func (c *Context) GetStore() *storage.Storage {
	return c.Context.Context.Value(keyStore).(*storage.Storage)
}
