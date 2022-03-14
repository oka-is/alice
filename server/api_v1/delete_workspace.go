package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/server/engine"
)

func DeleteWorkspace(ctx *engine.Context) {
	workspaceID := ctx.Param(paramWorkspaceID)
	// TODO: delete just a linked record if workspaces is shared
	err := ctx.GetStore().DeleteWorkspace(ctx, workspaceID)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.Done()
}
