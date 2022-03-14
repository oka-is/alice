package cmd

import (
	"github.com/urfave/cli/v2"
	"github.com/wault-pw/alice/pkg/pack"
	"github.com/wault-pw/alice/server/api_v1"
	"github.com/wault-pw/alice/server/cypress"
	"github.com/wault-pw/alice/server/engine"
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

	FlagServerVer666 = &cli.BoolFlag{
		Name:    "ver666",
		EnvVars: []string{"VER666"},
	}

	FlagServerMountCypress = &cli.BoolFlag{
		Name:    "mount-cypress",
		Usage:   "mount cypress endpoints for test cases only",
		EnvVars: []string{"MOUNT_CYPRESS"},
	}
)

func Server(ctx *cli.Context) error {
	opts := engine.Opts{
		AllowOrigin:  ctx.StringSlice(FlagServerAllowOrigin.Name),
		CookieSecure: ctx.Bool(FlagServerCookieSecure.Name),
		CookieDomain: ctx.String(FlagServerCookieDomain.Name),
		BackupUrl:    ctx.String(FlagServerBackupUrl.Name),
		Ver:          pack.NewWer(pack.Ver1),
	}

	if ctx.Bool(FlagServerVer666.Name) {
		opts.Ver = pack.NewWer(pack.Ver666)
	}

	routes := engine.New(Ctx(ctx).GetStore(), opts)
	api_v1.Extend(routes)

	if ctx.Bool(FlagServerMountCypress.Name) {
		cypress.Extend(routes)
	}

	return routes.Run(ctx.String(FlagServerAddress.Name))
}
