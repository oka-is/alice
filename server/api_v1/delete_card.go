package api_v1

import (
	"github.com/wault-pw/alice/server/engine"
)

func DeleteCard(ctx *engine.Context) {
	cardID := ctx.Param(paramCardID)

	user := ctx.MustGetUser()
	err := ctx.NewUserPolicy(user).CanWrite()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.GetStore().DeleteCard(ctx.Context, cardID)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.Done()
}
