package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func CreateShare(ctx *engine.Context) {
	req := new(alice_v1.CreateShareRequest)
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

	share := domain.UserWorkspace{
		WorkspaceID: uw.WorkspaceID,
		UserID:      domain.NewEmptyString(req.GetUserId()),
		AedKeyEnc:   domain.NewEmptyBytes(req.GetAedKeyEnc()),
		OwnerID:     uw.OwnerID,
	}

	err = ctx.GetStore().ShareUserWorkspace(ctx, &share)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.CreateShareResponse{
		UserWorkspace: mapper_v1.MapUserWorkspace(share),
	})
}
