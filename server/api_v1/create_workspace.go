package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func CreateWorkspace(ctx *engine.Context) {
	req := new(alice_v1.CreateWorkspaceRequest)
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

	userWorkspace, workspace := mapper_v1.BindCreateWorkspace(req)
	userWorkspace.OwnerID = user.ID
	userWorkspace.UserID = user.ID

	err = ctx.GetStore().CreateWorkspace(ctx.Context, &userWorkspace, &workspace)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	userWithWorkspace, err := ctx.GetStore().FindUserWithWorkspace(ctx, userWorkspace.ID.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.CreateWorkspaceResponse{
		Workspace: mapper_v1.MapUserWithWorkspace(userWithWorkspace),
	})
}
