package cmd

import (
	"fmt"

	"github.com/oka-is/alice/pkg/storage"
	"github.com/oka-is/alice/pkg/validator"
	"github.com/urfave/cli/v2"
)

var (
	FlagPostgresDSN = &cli.StringFlag{
		Name:    "pg-dsn",
		Aliases: []string{"psn"},
		EnvVars: []string{"PG_DSN"},
		Value:   "postgres://localhost:5432/alice?sslmode=disable&timezone=utc",
	}

	FlagSseKey = &cli.StringFlag{
		Name:    "sse-key",
		EnvVars: []string{"SSE_KEY"},
		Value:   "c72bd6f32260c5c19ac4d5e2617bcc097ca40470c26c5c0bd3c7ff8c0297476e547905043cdf8c90637d72b41086135fccd727bf16bc622b12dff46593ab4ecd",
	}
)

func BeforeStore(ctx *cli.Context) error {
	db, err := storage.Connect(ctx.String(FlagPostgresDSN.Name))
	if err != nil {
		return fmt.Errorf("filed to connect to store: %w", err)
	}

	store := storage.NewStorage(db, validator.New(), []byte(ctx.String(FlagSseKey.Name)))
	Ctx(ctx).SetStore(store)

	return nil
}
