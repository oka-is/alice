package api_v1

import (
	"github.com/wault-pw/alice/server/engine"
)

func DeleteCard(ctx *engine.Context) {
	user := ctx.MustGetUser()
	err := ctx.NewUserPolicy(user).CanWrite()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	card, err := ctx.GetStore().FindCard(ctx.Ctx(), ctx.Param(paramCardID))
	if err != nil {
		ctx.HandleError(err)
		return
	}

	uw, err := ctx.GetStore().FindUserWorkspaceLink(ctx.Ctx(), user.ID.String, ctx.Param(paramWorkspaceID))
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.NewWorkspacePolicy(user, uw).CanManageCard(card)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.GetStore().DeleteCard(ctx.Context, card.ID.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.Done()
}
