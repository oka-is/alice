package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/server/engine"
)

func ArchiveCard(ctx *engine.Context) {
	user := ctx.MustGetUser()
	err := ctx.NewUserPolicy(user).CanWrite()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	card, err := ctx.GetStore().FindCard(ctx.Ctx(), ctx.Param(paramCardID))
	if err != nil {
		ctx.HandleError(err)
		return
	}

	uw, err := ctx.GetStore().FindUserWorkspaceLink(ctx.Ctx(), user.ID.String, ctx.Param(paramWorkspaceID))
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.NewWorkspacePolicy(user, uw).CanManageCard(card)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	archived, err := ctx.GetStore().ArchiveCard(ctx.Context, card.ID.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.ArchiveCardResponse{
		Archived: archived,
	})
}
