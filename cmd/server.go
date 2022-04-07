package cmd

import (
	"fmt"

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

	FlagServerJwtKey = &cli.StringFlag{
		Name:    "jwt-key",
		Value:   "60582546f144e604e4a11c927bcf8bd82a0d7bbd4a31eeaa69ce11d69acd0a4a",
		EnvVars: []string{"JWT_KEY"},
	}

	FlagServerAllowOrigin = &cli.StringSliceFlag{
		Name:    "allow-origin",
		Aliases: []string{"ao"},
		EnvVars: []string{"ALLOW_ORIGIN"},
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

	FlagServerSentryDsn = &cli.StringFlag{
		Name:    "sentry-dsn",
		Usage:   "sentry dsn for error monitoring",
		EnvVars: []string{"SENTRY_DSN"},
	}
)

func Server(ctx *cli.Context) error {
	opts := engine.Opts{
		AllowOrigin:  ctx.StringSlice(FlagServerAllowOrigin.Name),
		CookieSecure: ctx.Bool(FlagServerCookieSecure.Name),
		CookieDomain: ctx.String(FlagServerCookieDomain.Name),
		Production:   ctx.Bool(FlagProduction.Name),
		BackupUrl:    ctx.String(FlagServerBackupUrl.Name),
		SentryDsn:    ctx.String(FlagServerSentryDsn.Name),
		JwtKey:       []byte(ctx.String(FlagServerJwtKey.Name)),
		Ver:          pack.NewWer(pack.Ver1),
	}
	opts.SetDefaultPolicies()

	if ctx.Bool(FlagServerVer666.Name) {
		opts.Ver = pack.NewWer(pack.Ver666)
	}

	routes, err := engine.New(Ctx(ctx).GetStore(), opts)
	if err != nil {
		return err
	}

	api_v1.Extend(routes)

	if ctx.Bool(FlagServerMountCypress.Name) {
		cypress.Extend(routes)
	}

	return routes.Run(ctx.String(FlagServerAddress.Name))
}

func BeforeServerProduction(ctx *cli.Context) error {
	if !ctx.Bool(FlagProduction.Name) {
		return nil
	}

	if !isStringFlagSet(ctx, FlagServerJwtKey) {
		return fmt.Errorf("please provide <%s>", FlagServerJwtKey.Name)
	}

	if isBoolFlagSet(ctx, FlagServerVer666) {
		return fmt.Errorf("flag <%s> cant be set in production", FlagServerVer666.Name)
	}

	if isBoolFlagSet(ctx, FlagServerMountCypress) {
		return fmt.Errorf("flag <%s> cant be set in production", FlagServerMountCypress.Name)
	}

	return nil
}
