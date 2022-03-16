package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func UpdateWorkspace(ctx *engine.Context) {
	req := new(alice_v1.UpdateWorkspaceRequest)
	err := ctx.MustBindProto(req)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	user := ctx.MustGetUser()
	err = ctx.NewUserPolicy(user).CanWrite()
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

	err = ctx.GetStore().UpdateWorkspace(ctx, uw.WorkspaceID.String, req.GetTitleEnc())
	if err != nil {
		ctx.HandleError(err)
		return
	}

	userWithWorkspace, err := ctx.GetStore().FindUserWithWorkspace(ctx, uw.ID.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.UpdateWorkspaceResponse{
		Workspace: mapper_v1.MapUserWithWorkspace(userWithWorkspace),
	})
}
