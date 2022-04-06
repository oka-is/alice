package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func ListShares(ctx *engine.Context) {
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

	shares, err := ctx.GetStore().ListSharedUserWorkspaces(ctx.Ctx(), uw.WorkspaceID.String, uw.OwnerID.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.ListShareResponse{
		Shares: mapper_v1.MapUserWorkspaces(shares),
	})
}
