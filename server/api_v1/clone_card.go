package api_v1

import (
	"net/http"

	"github.com/oka-is/alice/desc/alice_v1"
	"github.com/oka-is/alice/pkg/domain"
	"github.com/oka-is/alice/server/engine"
	"github.com/oka-is/alice/server/mapper_v1"
)

func CloneCard(ctx *engine.Context) {
	cardID, workspaceID := ctx.Param(paramCardID), ctx.Param(paramWorkspaceID)
	req := new(alice_v1.CloneCardRequest)
	if err := ctx.MustBindProto(req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
	}

	card := domain.Card{
		WorkspaceID: domain.NewEmptyString(workspaceID),
		TitleEnc:    domain.NewEmptyBytes(req.GetTitleEnc()),
	}
	err := ctx.GetStore().CloneCard(ctx.Context, &card, cardID)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.CloneCardResponse{
		Card: mapper_v1.MapCard(card),
	})
}
