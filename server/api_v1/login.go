package api_v1

import (
	"fmt"
	"net/http"

	"github.com/wault-pw/alice/desc/alice_v1"
	"github.com/wault-pw/alice/server/engine"
	"github.com/wault-pw/srp6ago"
	"google.golang.org/protobuf/proto"
)

func LoginAuth0(ctx *engine.Context) {
	req := new(alice_v1.Login0Request)
	err := ctx.MustBindProto(req)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	session, token, err := ctx.GetStore().IssueSession(ctx.Ctx(), ctx.JwtOpts())
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.SetCookieToken(token)
	ctx.SetSession(session)

	message, err := auth0(ctx, req)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, message)
}

func LoginAuth1(ctx *engine.Context) {
	req := new(alice_v1.Login1Request)
	err := ctx.MustBindProto(req)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	message, err := auth1(ctx, req)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, message)
}

func LoginOtp(ctx *engine.Context) {
	req := new(alice_v1.LoginOtpRequest)
	err := ctx.MustBindProto(req)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var message proto.Message
	if req.GetPasscode() == "" {
		message, err = loginOtpPre(ctx)
	} else {
		message, err = loginOtp(ctx, req)
	}

	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.ProtoBuf(http.StatusOK, message)
}

func loginOtpPre(ctx *engine.Context) (proto.Message, error) {
	user := ctx.MustGetUser()
	session := ctx.MustGetSession()

	if user.IsOtpDisabled() {
		return &alice_v1.LoginOtpResponse{Required: false},
			ctx.GetStore().OtpSessionSucceed(ctx.Ctx(), session.Jti.String)
	}

	return &alice_v1.LoginOtpResponse{Required: true}, nil
}

func loginOtp(ctx *engine.Context, req *alice_v1.LoginOtpRequest) (proto.Message, error) {
	user := ctx.MustGetUser()
	session := ctx.MustGetSession()

	if ctx.IsOtpValid(req.GetPasscode(), string(user.OtpSecret.Bytea)) {
		return &alice_v1.LoginOtpResponse{Required: false},
			ctx.GetStore().OtpSessionSucceed(ctx.Ctx(), session.Jti.String)
	}

	return &alice_v1.LoginOtpResponse{Required: true}, nil
}

func auth0(ctx *engine.Context, req *alice_v1.Login0Request) (proto.Message, error) {
	user, err := ctx.GetStore().FindUserIdentity(ctx.Ctx(), req.GetIdentity())
	if err != nil {
		return nil, fmt.Errorf("find user identity failed: %w", err)
	}

	srp := ctx.MustGetVer().NewSrpServer(user.Verifier.Bytea, user.SrpSalt.Bytea)
	mutual, err := srp.PublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get srp publick key: %w", err)
	}

	session := ctx.MustGetSession()
	err = ctx.GetStore().CandidateSession(ctx.Ctx(), session.Jti.String, user.ID.String, srp.Marshal())
	if err != nil {
		return nil, fmt.Errorf("failed to candidate session: %w", err)
	}

	return &alice_v1.Login0Response{
		Mutual: mutual,
		Salt:   user.SrpSalt.Bytea,
	}, nil
}

func auth1(ctx *engine.Context, req *alice_v1.Login1Request) (proto.Message, error) {
	session := ctx.MustGetSession()
	srp, err := srp6ago.UnmarshalServer(session.SrpState.Bytea)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall srp server: %w", err)
	}

	err = srp.SetClientPublicKey(req.GetMutual())
	if err != nil {
		return nil, fmt.Errorf("server aborts")
	}

	if !srp.IsProofValid(req.GetProof()) {
		return nil, fmt.Errorf("invalid credentials")
	}

	err = ctx.GetStore().NominateSession(ctx.Ctx(), ctx.MustGetSession().Jti.String)
	if err != nil {
		return nil, fmt.Errorf("failed to candidate session: %w", err)
	}

	return &alice_v1.Login1Response{
		Proof: srp.Proof(),
	}, nil
}
