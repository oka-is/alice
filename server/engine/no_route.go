package engine

import (
	"net/http"
)

func noRoute(ctx *Context) {
	ctx.String(http.StatusNotFound, "[404]")
}
