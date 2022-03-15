package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/server/engine"
)

func Terminate(ctx *engine.Context) {
	req := new(alice_v1.TerminateRequest)
	if err := ctx.MustBindProto(req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
	}

	user := ctx.MustGetUser()
	err := ctx.NewUserPolicy(user).CanWrite()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.GetStore().TerminateUser(ctx.Ctx(), req.GetIdentity(), user.ID.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.Done()
}
