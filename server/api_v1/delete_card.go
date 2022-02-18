package api_v1

import (
	"net/http"

	"github.com/oka-is/alice/server/engine"
)

func DeleteCard(ctx *engine.Context) {
	cardID := ctx.Param(paramCardID)
	err := ctx.GetStore().DeleteCard(ctx.Context, cardID)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.Done()
}
