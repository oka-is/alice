package api_v1

import (
	"github.com/wault-pw/alice/server/engine"
)

func DeleteSession(ctx *engine.Context) {
	err := ctx.GetStore().DeleteSession(ctx.Ctx(), ctx.MustGetSession().Jti.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.Done()
}
