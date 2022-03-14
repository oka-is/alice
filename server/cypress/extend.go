package cypress

import (
	"github.com/wault-pw/alice/server/engine"
)

func Extend(router *engine.Engine) *engine.Engine {
	group := router.Group("/__cypress__")
	group.POST("/truncate", engine.WrapAction(Truncate))
	return router
}
