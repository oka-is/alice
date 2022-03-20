package api_v1

import (
	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/server/engine"
)

func OtpDisable(ctx *engine.Context) {
	req := new(alice_v1.OtpEnableRequest)
	err := ctx.MustBindProto(req)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	user := ctx.MustGetUser()
	err = ctx.NewUserPolicy(user).CanWrite()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.GetStore().DisableUserOtp(ctx.Ctx(), user.ID.String)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.Done()
}
