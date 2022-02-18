package api_v1

import (
	"net/http"

	"github.com/oka-is/alice/desc/alice_v1"
	"github.com/oka-is/alice/pkg/validator"
	"github.com/oka-is/alice/server/engine"
	"github.com/oka-is/alice/server/mapper_v1"
)

func Register(ctx *engine.Context) {
	req := new(alice_v1.RegistrationRequest)
	if err := ctx.MustBindProto(req); err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
	}

	user, userWorkspace, workspace, cards := mapper_v1.BindRegistration(req)
	err := ctx.GetStore().CreateUser(ctx.Context, &user, &userWorkspace, &workspace, cards)
	switch {
	case validator.IsInvalid(err):
		_ = ctx.AbortWithError(http.StatusUnprocessableEntity, err)
	case err != nil:
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.Done()
}
