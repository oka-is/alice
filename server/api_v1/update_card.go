package api_v1

import (
	"net/http"

	"github.com/oka-is/alice/desc/alice_v1"
	"github.com/oka-is/alice/pkg/domain"
	"github.com/oka-is/alice/pkg/validator"
	"github.com/oka-is/alice/server/engine"
	"github.com/oka-is/alice/server/mapper_v1"
)

func UpdateCard(ctx *engine.Context) {
	req := new(alice_v1.UpsertCardRequest)
	if err := ctx.MustBindProto(req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	card, items := mapper_v1.BindUpsertCard(req)
	card.ID = domain.NewEmptyString(ctx.Param(paramCardID))
	err := ctx.GetStore().UpdateCardWithItems(ctx.Context, &card, items)
	switch {
	case validator.IsInvalid(err):
		_ = ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	case err != nil:
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.UpsertCardResponse{
		Card: mapper_v1.MapCard(card),
	})
}
