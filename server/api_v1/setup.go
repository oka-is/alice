package api_v1

import (
	"fmt"
	"net/http"

	"github.com/oka-is/alice/server/engine"
)

func Extend(router *engine.Engine) *engine.Engine {
	router.POST("/v1/register", engine.WrapAction(Register))
	router.POST("/v1/login/cookie", engine.WrapAction(LoginCookie))

	// BEGIN:AUTHENTICATED
	auth := router.Group("/", engine.WrapAction(useAuth))
	auth.POST("/v1/login/auth0", engine.WrapAction(LoginAuth0))
	auth.POST("/v1/login/auth1", engine.WrapAction(LoginAuth1))
	auth.POST("/v1/whoami", engine.WrapAction(WhoAmI))
	auth.GET("/v1/backups/create", engine.WrapAction(CreateBackup))
	auth.POST("/v1/workspaces/list", engine.WrapAction(ListWorkspaces))
	auth.POST("/v1/workspaces/create", engine.WrapAction(CreateWorkspace))
	auth.POST(fmt.Sprintf("/v1/workspaces/:%s/cards", paramWorkspaceID), engine.WrapAction(ListCards))
	auth.POST(fmt.Sprintf("/v1/workspaces/:%s/delete", paramWorkspaceID), engine.WrapAction(DeleteWorkspace))
	auth.POST(fmt.Sprintf("/v1/workspaces/:%s/cards/:%s/items", paramWorkspaceID, paramCardID), engine.WrapAction(ListCardItems))
	auth.POST(fmt.Sprintf("/v1/workspaces/:%s/cards/create", paramWorkspaceID), engine.WrapAction(CreateCard))
	auth.POST(fmt.Sprintf("/v1/workspaces/:%s/cards/:%s/delete", paramWorkspaceID, paramCardID), engine.WrapAction(DeleteCard))
	auth.POST(fmt.Sprintf("/v1/workspaces/:%s/cards/:%s/clone", paramWorkspaceID, paramCardID), engine.WrapAction(CloneCard))
	// END:AUTHENTICATED

	return router
}

func useAuth(ctx *engine.Context) {
	token, err := ctx.GetCookieToken()
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	session, err := ctx.GetStore().RetrieveSession(ctx.Ctx(), ctx.JwtOpts(), token)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	ctx.SetSession(session)
	ctx.Next()
}
