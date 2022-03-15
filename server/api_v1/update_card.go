package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/domain"
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

	card, items := mapper_v1.BindUpsertCard(req)
	card.ID = domain.NewEmptyString(ctx.Param(paramCardID))

	err = ctx.GetStore().UpdateCardWithItems(ctx.Context, &card, items)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.UpsertCardResponse{
		Card: mapper_v1.MapCard(card),
	})
}
