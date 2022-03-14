package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func ListWorkspaces(ctx *engine.Context) {
	session := ctx.MustGetSession()
	records, err := ctx.GetStore().ListUserWithWorkspaces(ctx, session.UserID.String)

	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, mapper_v1.MapListUserWorkspaceResponse(records))
}
