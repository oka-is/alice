package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/wault-pw/alice/pkg/storage"
	"github.com/wault-pw/alice/pkg/validator"
)

var (
	FlagPostgresDSN = &cli.StringFlag{
		Name:    "pg-dsn",
		Usage:   "postgres database connection string",
		Aliases: []string{"psn"},
		EnvVars: []string{"PG_DSN"},
		Value:   "postgres://localhost:5432/alice?sslmode=disable&timezone=utc",
	}

	FlagSseKey = &cli.StringFlag{
		Name:    "sse-key",
		Usage:   "key for server-side database column encryption",
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

func BeforeStoreProduction(ctx *cli.Context) error {
	if !ctx.Bool(FlagProduction.Name) {
		return nil
	}

	if !isStringFlagSet(ctx, FlagSseKey) {
		return fmt.Errorf("please provide <%s>", FlagSseKey.Name)
	}

	return nil
}
