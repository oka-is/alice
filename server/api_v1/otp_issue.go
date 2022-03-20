package api_v1

import (
	"net/http"

	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/alice/server/mapper_v1"
)

func OtpIssue(ctx *engine.Context) {
	user := ctx.MustGetUser()
	err := ctx.NewUserPolicy(user).CanWrite()
	if err != nil {
		ctx.HandleError(err)
		return
	}

	secret, url, err := ctx.OtpIssue(user)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	err = ctx.GetStore().IssueUserOtp(ctx.Ctx(), user.ID.String, secret)
	if err != nil {
		ctx.HandleError(err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, mapper_v1.MapOtpIssue(secret, url))
}
