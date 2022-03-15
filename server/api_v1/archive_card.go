package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/server/engine"
)

func ArchiveCard(ctx *engine.Context) {
	cardID, _ := ctx.Param(paramCardID), ctx.Param(paramWorkspaceID)

	user := ctx.MustGetUser()
	err := ctx.NewUserPolicy(user).CanWrite()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	archived, err := ctx.GetStore().ArchiveCard(ctx.Context, cardID)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, &alice_v1.ArchiveCardResponse{
		Archived: archived,
	})
}
