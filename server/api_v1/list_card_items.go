package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func ListCardItems(ctx *engine.Context) {
	cardID := ctx.Param(paramCardID)
	items, err := ctx.GetStore().ListCardItems(ctx.Context, cardID)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.ListCardItemsResponse{
		Items: mapper_v1.MapCardItems(items),
	})
}
