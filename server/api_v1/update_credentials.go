package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func UpdateCredentials(ctx *engine.Context) {
	req := new(alice_v1.UpdateCredentialsRequest)
	err := ctx.MustBindProto(req)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	user, session := ctx.MustGetUser(), ctx.MustGetSession()
	err = ctx.NewUserPolicy(user).CanWrite()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	oldIdentity, newUser := mapper_v1.BindUpdateCredentials(req)
	err = ctx.GetStore().UpdateCredentials(ctx.Ctx(), user.ID.String, oldIdentity, newUser)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.GetStore().DeleteUserSessionExcept(ctx.Ctx(), user.ID.String, session.Jti.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	user, err = ctx.GetStore().FindUser(ctx.Ctx(), user.ID.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.UpdateCredentialsResponse{
		User: mapper_v1.MapPrivUser(user),
	})
}
