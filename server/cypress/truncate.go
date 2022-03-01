package cypress

import (
	"net/http"

	"github.com/oka-is/alice/server/engine"
)

func Truncate(ctx *engine.Context) {
	err := ctx.GetStore().TruncateAll(ctx.Ctx())
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Done()
}
