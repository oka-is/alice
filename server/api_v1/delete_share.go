package api_v1

import (
	"github.com/wault-pw/alice/server/engine"
)

func DeleteShare(ctx *engine.Context) {
	user := ctx.MustGetUser()
	err := ctx.NewUserPolicy(user).CanWrite()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	uw, err := ctx.GetStore().FindUserWorkspace(ctx.Ctx(), ctx.Param(paramShareID))
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.NewWorkspacePolicy(user, uw).CanDeleteShare()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.GetStore().DeleteUserWorkspace(ctx.Ctx(), uw.ID.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.Done()
}
