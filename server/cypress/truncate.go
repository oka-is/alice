package cypress

import (
	"github.com/oka-is/alice/server/engine"
)

func Truncate(ctx *engine.Context) {
	ctx.Done()
}
