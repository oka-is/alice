package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func WhoAmI(ctx *engine.Context) {
	user := ctx.MustGetUser()
	ctx.ProtoBuf(http.StatusOK, mapper_v1.MapWhoAmI(user))
}
