package api_v1

import (
	"net/http"

	"github.com/oka-is/alice/server/engine"
)

func DeleteSession(ctx *engine.Context) {
	err := ctx.GetStore().DeleteSession(ctx.Ctx(), ctx.MustGetSession().Jti.String)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.Done()
}
