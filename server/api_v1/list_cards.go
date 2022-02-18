package api_v1

import (
	"net/http"

	"github.com/oka-is/alice/server/engine"
	"github.com/oka-is/alice/server/mapper_v1"
)

func ListCards(ctx *engine.Context) {
	workspaceID := ctx.Param(paramWorkspaceID)
	cards, err := ctx.GetStore().ListCardsByWorkspace(ctx.Context, workspaceID)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.ProtoBuf(http.StatusOK, mapper_v1.MapListCardsResponse(cards))
}
