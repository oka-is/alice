package cypress

import (
	"github.com/oka-is/alice/server/engine"
)

func Extend(router *engine.Engine) *engine.Engine {
	group := router.Group("/__cypress__")
	group.POST("/truncate", engine.WrapAction(Truncate))
	return router
}
