package cmd

import (
	"github.com/urfave/cli/v2"
)

func BeforeAll(actions ...cli.BeforeFunc) cli.BeforeFunc {
	return func(ctx *cli.Context) error {
		for _, action := range actions {
			if err := action(ctx); err != nil {
				return err
			}
		}

		return nil
	}
}
