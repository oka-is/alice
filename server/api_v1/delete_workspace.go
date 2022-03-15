package api_v1

import (
	"github.com/wault-pw/alice/server/engine"
)

func DeleteWorkspace(ctx *engine.Context) {
	user := ctx.MustGetUser()
	err := ctx.NewUserPolicy(user).CanWrite()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	uw, err := ctx.GetStore().FindUserWorkspaceLink(ctx.Ctx(), user.ID.String, ctx.Param(paramWorkspaceID))
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.NewWorkspacePolicy(user, uw).CanManageWorkspace()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.GetStore().DeleteWorkspace(ctx, uw.WorkspaceID.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.Done()
}
