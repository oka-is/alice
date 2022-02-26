package api_v1

import (
	"fmt"
	"net/http"

	"github.com/oka-is/alice/desc/alice_v1"
	"github.com/oka-is/alice/server/engine"
	"github.com/oka-is/srp6ago"
	"google.golang.org/protobuf/proto"
)

func LoginCookie(ctx *engine.Context) {
	_, token, err := ctx.GetStore().IssueSession(ctx.Ctx(), ctx.JwtOpts())
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.SetCookieToken(token)
	ctx.Done()
}

func LoginAuth0(ctx *engine.Context) {
	req := new(alice_v1.Login0Request)
	err := ctx.MustBindProto(req)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

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

	srp.SetClientPublicKey(req.GetMutual())
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
