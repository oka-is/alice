package api_v1

import (
	"net/http"

	"github.com/oka-is/alice/desc/alice_v1"
	"github.com/oka-is/alice/pkg/validator"
	"github.com/oka-is/alice/server/engine"
)

func Terminate(ctx *engine.Context) {
	req := new(alice_v1.TerminateRequest)
	if err := ctx.MustBindProto(req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
	}

	userID := ctx.MustGetSession().UserID.String
	err := ctx.GetStore().TerminateUser(ctx.Ctx(), req.GetIdentity(), userID)
	switch {
	case validator.IsInvalid(err):
		_ = ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	case err != nil:
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Done()
}
