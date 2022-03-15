package cmd

import (
	"github.com/urfave/cli/v2"
)

var (
	FlagProduction = &cli.BoolFlag{
		Name:    "production",
		EnvVars: []string{"PRODUCTION"},
	}
)
