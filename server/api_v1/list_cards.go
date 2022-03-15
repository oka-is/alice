package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func ListCards(ctx *engine.Context) {
	user := ctx.MustGetUser()
	uw, err := ctx.GetStore().FindUserWorkspaceLink(ctx.Ctx(), user.ID.String, ctx.Param(paramWorkspaceID))
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.NewWorkspacePolicy(user, uw).CanSeeWorkspace()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	cards, err := ctx.GetStore().ListCardsByWorkspace(ctx.Context, uw.WorkspaceID.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, mapper_v1.MapListCardsResponse(cards))
}
