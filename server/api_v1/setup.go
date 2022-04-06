package api_v1

import (
	"fmt"
	"net/http"

	"github.com/wault-pw/alice/server/engine"
)

func Extend(router *engine.Engine) *engine.Engine {
	router.POST("/v1/register", engine.WrapAction(Register))
	router.POST("/v1/login/auth0", engine.WrapAction(LoginAuth0))

	// BEGIN:AUTHENTICATED
	auth := router.Group("/", engine.WrapAction(useAuth))
	auth.POST("/v1/login/auth1", engine.WrapAction(LoginAuth1))
	auth.POST("/v1/login/otp", engine.WrapAction(LoginOtp))
	// END:AUTHENTICATED

	// BEGIN:AUTHENTICATED+OTP
	otp := router.Group("/", engine.WrapAction(useAuth), engine.WrapAction(useOtpCheck))
	otp.POST("/v1/whoami", engine.WrapAction(WhoAmI))
	otp.POST("/v1/users/find", engine.WrapAction(FindUser))
	otp.POST("/v1/otp/issue", engine.WrapAction(OtpIssue))
	otp.POST("/v1/otp/enable", engine.WrapAction(OtpEnable))
	otp.POST("/v1/otp/disable", engine.WrapAction(OtpDisable))
	otp.POST("/v1/credentials/update", engine.WrapAction(UpdateCredentials))
	otp.POST("/v1/sessions/delete", engine.WrapAction(DeleteSession))
	otp.POST("/v1/terminate", engine.WrapAction(Terminate))
	otp.GET("/v1/backups/create", engine.WrapAction(CreateBackup))
	otp.POST("/v1/workspaces/list", engine.WrapAction(ListWorkspaces))
	otp.POST("/v1/workspaces/create", engine.WrapAction(CreateWorkspace))
	otp.POST(fmt.Sprintf("/v1/shares/:%s/delete", paramShareID), engine.WrapAction(DeleteShare))
	otp.POST(fmt.Sprintf("/v1/workspaces/:%s/cards", paramWorkspaceID), engine.WrapAction(ListCards))
	otp.POST(fmt.Sprintf("/v1/workspaces/:%s/delete", paramWorkspaceID), engine.WrapAction(DeleteWorkspace))
	otp.POST(fmt.Sprintf("/v1/workspaces/:%s/update", paramWorkspaceID), engine.WrapAction(UpdateWorkspace))
	otp.POST(fmt.Sprintf("/v1/workspaces/:%s/shares/create", paramWorkspaceID), engine.WrapAction(CreateShare))
	otp.POST(fmt.Sprintf("/v1/workspaces/:%s/shares/list", paramWorkspaceID), engine.WrapAction(ListShares))
	otp.POST(fmt.Sprintf("/v1/workspaces/:%s/cards/:%s/items", paramWorkspaceID, paramCardID), engine.WrapAction(ListCardItems))
	otp.POST(fmt.Sprintf("/v1/workspaces/:%s/cards/:%s/archive", paramWorkspaceID, paramCardID), engine.WrapAction(ArchiveCard))
	otp.POST(fmt.Sprintf("/v1/workspaces/:%s/cards/create", paramWorkspaceID), engine.WrapAction(CreateCard))
	otp.POST(fmt.Sprintf("/v1/workspaces/:%s/cards/:%s/update", paramWorkspaceID, paramCardID), engine.WrapAction(UpdateCard))
	otp.POST(fmt.Sprintf("/v1/workspaces/:%s/cards/:%s/delete", paramWorkspaceID, paramCardID), engine.WrapAction(DeleteCard))
	otp.POST(fmt.Sprintf("/v1/workspaces/:%s/cards/:%s/clone", paramWorkspaceID, paramCardID), engine.WrapAction(CloneCard))
	// END:AUTHENTICATED+OTP

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

func useOtpCheck(ctx *engine.Context) {
	session := ctx.MustGetSession()

	if !session.OtpSucceed.Bool {
		_ = ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("OTP is not succeed"))
		return
	}

	ctx.Next()
}
