package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func ListWorkspaces(ctx *engine.Context) {
	user := ctx.MustGetUser()
	records, err := ctx.GetStore().ListUserWithWorkspaces(ctx, user.ID.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, mapper_v1.MapListUserWorkspaceResponse(records))
}
