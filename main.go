package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/wault-pw/alice/cmd"
)

var (
	Version = ":VERSION:"
)

func main() {
	app := &cli.App{
		Name:      "alice",
		Usage:     "password manager backend",
		Version:   Version,
		Copyright: "© 2022 Wault OÜ",
		Flags: []cli.Flag{
			cmd.FlagProduction,
		},
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "start API server",
				Flags: []cli.Flag{
					cmd.FlagPostgresDSN,
					cmd.FlagSseKey,
					cmd.FlagServerAddress,
					cmd.FlagServerAllowOrigin,
					cmd.FlagServerCookieDomain,
					cmd.FlagServerCookieSecure,
					cmd.FlagServerJwtKey,
					cmd.FlagServerBackupUrl,
					cmd.FlagServerVer666,
					cmd.FlagServerMountCypress,
					cmd.FlagServerSentryDsn,
				},
				Before: cmd.BeforeAll(
					cmd.BeforeStoreProduction,
					cmd.BeforeServerProduction,
					cmd.BeforeStore,
				),
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
