package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/pkg/validator"
	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func Register(ctx *engine.Context) {
	req := new(alice_v1.RegistrationRequest)
	if err := ctx.MustBindProto(req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, userWorkspace, workspace, cards := mapper_v1.BindRegistration(req)
	err := ctx.GetStore().CreateUser(ctx.Context, &user, &userWorkspace, &workspace, cards)
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
