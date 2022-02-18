package api_v1

import (
	"net/http"

	"github.com/oka-is/alice/desc/alice_v1"
	"github.com/oka-is/alice/pkg/domain"
	"github.com/oka-is/alice/pkg/validator"
	"github.com/oka-is/alice/server/engine"
	"github.com/oka-is/alice/server/mapper_v1"
)

func CreateCard(ctx *engine.Context) {
	req := new(alice_v1.CreateCardRequest)
	if err := ctx.MustBindProto(req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
	}

	card, items := mapper_v1.BindCreateCard(req)
	card.WorkspaceID = domain.NewEmptyString(ctx.Param(paramWorkspaceID))
	err := ctx.GetStore().CreateCardWithItems(ctx.Context, &card, items)
	switch {
	case validator.IsInvalid(err):
		_ = ctx.AbortWithError(http.StatusUnprocessableEntity, err)
	case err != nil:
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.CreateCardResponse{
		Card: mapper_v1.MapCard(card),
	})
}
