package api_v1

import (
	"github.com/wault-pw/alice/server/engine"
)

func DeleteWorkspace(ctx *engine.Context) {
	workspaceID := ctx.Param(paramWorkspaceID)

	user := ctx.MustGetUser()
	err := ctx.NewUserPolicy(user).CanWrite()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.GetStore().DeleteWorkspace(ctx, workspaceID)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.Done()
}
