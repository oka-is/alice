package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/lib/uuid"
	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func FindUser(ctx *engine.Context) {
	req := new(alice_v1.FindUserRequest)
	err := ctx.MustBindProto(req)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	user, err := ctx.GetStore().FindUser(ctx.Ctx(), uuid.Safe(req.GetId()))
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.FindUserResponse{
		User: mapper_v1.MapPubUser(user),
	})
}
