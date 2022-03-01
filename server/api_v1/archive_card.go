package api_v1

import (
	"net/http"

	"github.com/oka-is/alice/desc/alice_v1"
	"github.com/oka-is/alice/server/engine"
)

func ArchiveCard(ctx *engine.Context) {
	cardID, _ := ctx.Param(paramCardID), ctx.Param(paramWorkspaceID)
	archived, err := ctx.GetStore().ArchiveCard(ctx.Context, cardID)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.ArchiveCardResponse{
		Archived: archived,
	})
}
