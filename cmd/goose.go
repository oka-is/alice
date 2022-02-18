package cmd

import (
	"embed"

	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v2"
)

//go:generate rm -rf ./embed/migrations
//go:generate cp -r ../migrations ./embed
//go:embed embed
var embedMigrations embed.FS

const migrationsDir = "embed/migrations"

func BeforeGoose(ctx *cli.Context) error {
	goose.SetBaseFS(embedMigrations)
	return nil
}

func GooseStatus(ctx *cli.Context) error {
	return goose.Status(Ctx(ctx).GetStore().SqlDB(), migrationsDir)
}

func GooseUP(ctx *cli.Context) error {
	return goose.Up(Ctx(ctx).GetStore().SqlDB(), migrationsDir)
}

func GooseDOWN(ctx *cli.Context) error {
	return goose.Down(Ctx(ctx).GetStore().SqlDB(), migrationsDir)
}
