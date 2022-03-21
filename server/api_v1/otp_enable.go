package api_v1

import (
	"fmt"
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/server/engine"
)

func OtpEnable(ctx *engine.Context) {
	req := new(alice_v1.OtpEnableRequest)
	err := ctx.MustBindProto(req)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	user, session := ctx.MustGetUser(), ctx.MustGetSession()
	err = ctx.NewUserPolicy(user).CanWrite()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	if !ctx.IsOtpValid(req.GetPasscode(), string(user.OtpCandidate.Bytea)) {
		_ = ctx.AbortWithError(http.StatusUnprocessableEntity, fmt.Errorf("OTP passcode is invalid"))
		return
	}

	err = ctx.GetStore().EnableUserOtp(ctx.Ctx(), user.ID.String, req.GetIdentity(), user.OtpCandidate.Bytea)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.GetStore().DeleteUserSessionExcept(ctx.Ctx(), user.ID.String, session.Jti.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.Done()
}
