package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func ListCardItems(ctx *engine.Context) {
	user := ctx.MustGetUser()
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

	err = ctx.NewWorkspacePolicy(user, uw).CanSeeCard(card)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	items, err := ctx.GetStore().ListCardItems(ctx.Context, card.ID.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.ListCardItemsResponse{
		Items: mapper_v1.MapCardItems(items),
	})
}
