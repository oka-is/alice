package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func UpdateCard(ctx *engine.Context) {
	req := new(alice_v1.UpsertCardRequest)
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

	card, err := ctx.GetStore().FindCard(ctx.Ctx(), ctx.Param(paramCardID))
	if err != nil {
		ctx.HandleError(err)
		return
	}

	uw, err := ctx.GetStore().FindUserWorkspaceLink(ctx.Ctx(), user.ID.String, ctx.Param(paramWorkspaceID))
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.NewWorkspacePolicy(user, uw).CanManageCard(card)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	clone, items := mapper_v1.BindUpsertCard(req)
	clone.ID = card.ID

	err = ctx.GetStore().UpdateCardWithItems(ctx.Context, &clone, items)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.UpsertCardResponse{
		Card: mapper_v1.MapCard(clone),
	})
}
