package cmd

import (
	"github.com/oka-is/alice/pkg/pack"
	"github.com/oka-is/alice/server/api_v1"
	"github.com/oka-is/alice/server/engine"
	"github.com/urfave/cli/v2"
)

var (
	FlagServerAddress = &cli.StringFlag{
		Name:    "address",
		Aliases: []string{"a"},
		EnvVars: []string{"ADDRESS"},
		Value:   "0.0.0.0:8080",
	}

	FlagServerCookieDomain = &cli.StringFlag{
		Name:    "cookie-domain",
		EnvVars: []string{"COOKIE_DOMAIN"},
	}

	FlagServerCookieSecure = &cli.BoolFlag{
		Name:    "cookie-secure",
		EnvVars: []string{"COOKIE_SECURE"},
	}

	FlagServerAllowOrigin = &cli.StringSliceFlag{
		Name:    "allow-origin",
		Aliases: []string{"ao"},
		Value: cli.NewStringSlice(
			"http://localhost:3000",
			"http://0.0.0.0:3000",
			"http://192.168.1.2:3000",
		),
	}

	FlagServerBackupUrl = &cli.StringFlag{
		Name:  "backup-url",
		Value: "http://localhost:3000/backup.html",
	}
)

func Server(ctx *cli.Context) error {
	routes := engine.New(Ctx(ctx).GetStore(), engine.Opts{
		AllowOrigin:  ctx.StringSlice(FlagServerAllowOrigin.Name),
		CookieSecure: ctx.Bool(FlagServerCookieSecure.Name),
		CookieDomain: ctx.String(FlagServerCookieDomain.Name),
		BackupUrl:    ctx.String(FlagServerBackupUrl.Name),
		Ver:          pack.NewWer(pack.Ver1),
	})

	api_v1.Extend(routes)
	return routes.Run(ctx.String(FlagServerAddress.Name))
}
