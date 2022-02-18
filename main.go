package main

import (
	"fmt"
	"os"

	"github.com/oka-is/alice/cmd"
	"github.com/urfave/cli/v2"
)

var (
	Version = ":VERSION:"
)

func main() {
	app := &cli.App{
		Name:      "alice",
		Usage:     "password manager backend",
		Version:   Version,
		Copyright: "© 2022 OKA OÜ",
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "start API server",
				Flags: []cli.Flag{
					cmd.FlagPostgresDSN,
					cmd.FlagServerAddress,
					cmd.FlagServerAllowOrigin,
					cmd.FlagServerCookieDomain,
					cmd.FlagServerCookieSecure,
					cmd.FlagServerBackupUrl,
				},
				Before: cmd.BeforeAll(cmd.BeforeStore),
				Action: cmd.Server,
			},
			{
				Name:   "goose",
				Usage:  "database migrations",
				Flags:  []cli.Flag{cmd.FlagPostgresDSN},
				Before: cmd.BeforeAll(cmd.BeforeStore, cmd.BeforeGoose),
				Subcommands: []*cli.Command{
					{
						Name:   "status",
						Usage:  "migration status",
						Action: cmd.GooseStatus,
					}, {
						Name:   "up",
						Usage:  "run pending migrations",
						Action: cmd.GooseUP,
					}, {
						Name:   "down",
						Usage:  "rollback previous migration",
						Action: cmd.GooseDOWN,
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)
	}
}
