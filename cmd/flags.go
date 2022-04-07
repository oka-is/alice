package cmd

import "github.com/urfave/cli/v2"

// IsStringFlagSet check is flag was set from ENV
// or if flag value IS NOT a default one.
// When flag sets from ENV, cli overrides initial default flag value
func isStringFlagSet(ctx *cli.Context, flag *cli.StringFlag) bool {
	return flag.IsSet() || ctx.String(flag.Name) != flag.Value
}

func isBoolFlagSet(ctx *cli.Context, flag *cli.BoolFlag) bool {
	return ctx.Bool(flag.Name)
}
