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
)

func BeforeStore(ctx *cli.Context) error {
	store, err := storage.Connect(ctx.String(FlagPostgresDSN.Name), validator.New())
	if err != nil {
		return fmt.Errorf("filed to connect to store: %w", err)
	}

	Ctx(ctx).SetStore(store)
	return nil
}
