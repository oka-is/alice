package api_v1

import (
	"net/http"

	"github.com/oka-is/alice/desc/alice_v1"
	"github.com/oka-is/alice/pkg/validator"
	"github.com/oka-is/alice/server/engine"
	"github.com/oka-is/alice/server/mapper_v1"
)

func CreateWorkspace(ctx *engine.Context) {
	req := new(alice_v1.CreateWorkspaceRequest)
	session := ctx.MustGetSession()
	err := ctx.MustBindProto(req)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
	}

	userWorkspace, workspace := mapper_v1.BindCreateWorkspace(req)
	userWorkspace.OwnerID = session.UserID
	userWorkspace.UserID = session.UserID

	err = ctx.GetStore().CreateWorkspace(ctx.Context, &userWorkspace, &workspace)
	switch {
	case validator.IsInvalid(err):
		_ = ctx.AbortWithError(http.StatusUnprocessableEntity, err)
	case err != nil:
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	userWithWorkspace, err := ctx.GetStore().FindUserWithWorkspace(ctx, userWorkspace.ID.String)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.CreateWorkspaceResponse{
		Workspace: mapper_v1.MapUserWithWorkspace(userWithWorkspace),
	})
}
