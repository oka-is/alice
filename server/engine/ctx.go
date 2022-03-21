package engine

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/wault-pw/alice/lib/jwt"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/pkg/pack"
	"github.com/wault-pw/alice/pkg/storage"
	"github.com/wault-pw/alice/pkg/validator"
	"github.com/wault-pw/alice/server/policy"
	"google.golang.org/protobuf/proto"
)

const (
	keyStore = string(rune(iota))
	keySession
	keyOpts
)

const (
	jwtCookie = "jwt"
)

type Context struct {
	*gin.Context
}

func Ctx(ctx *gin.Context) *Context {
	return &Context{ctx}
}

func (c *Context) SetStore(store storage.IStore) *Context {
	c.Set(keyStore, store)
	return c
}

func (c *Context) GetStore() storage.IStore {
	store, _ := c.Get(keyStore)
	return store.(storage.IStore)
}

func (c *Context) SetOpts(opts Opts) *Context {
	c.Set(keyOpts, opts)
	return c
}

func (c *Context) MustGetOpts() Opts {
	opts, exists := c.Get(keyOpts)
	if !exists {
		panic("engine opts missing from context")
	}
	return opts.(Opts)
}

func (c *Context) GetSession() (domain.Session, bool) {
	session, exists := c.Get(keySession)
	return session.(domain.Session), exists
}

func (c *Context) MustGetSession() domain.Session {
	session, exists := c.GetSession()
	if !exists {
		panic("session not exists in context")
	}
	return session
}

// MustGetUser returns session user
func (c *Context) MustGetUser() domain.User {
	user, err := c.GetStore().FindUser(c.Ctx(), c.MustGetSession().UserID.String)
	if err != nil {
		panic(fmt.Errorf("failed to get session user: %w", err))
	}
	return user
}

func (c *Context) MustGetVer() *pack.Ver {
	ver := c.MustGetOpts().Ver
	if ver == nil {
		panic("pack ver not exists in context")
	}

	return ver
}

func (c *Context) SetSession(session domain.Session) {
	c.Set(keySession, session)
}

func (c *Context) Ctx() context.Context {
	return context.Background()
}

func (c *Context) SetCookieToken(token string) {
	opts := c.MustGetOpts()

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     jwtCookie,
		Value:    url.QueryEscape(token),
		MaxAge:   0,
		Path:     "/",
		Domain:   opts.CookieDomain,
		SameSite: http.SameSiteStrictMode,
		Secure:   opts.CookieSecure,
		HttpOnly: true,
	})
}

func (c *Context) GetCookieToken() (string, error) {
	return c.Cookie(jwtCookie)
}

func (c *Context) JwtOpts() jwt.Opts {
	return jwt.Opts{
		Aud:  "WEB",
		Sub:  "API/V1",
		Iss:  "ALICE",
		Jti:  domain.NewUUID(),
		Key:  c.MustGetOpts().JwtKey,
		Algo: jwt.HS256,
	}
}

func (c *Context) MustBindProto(m proto.Message) error {
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	return proto.Unmarshal(buf, m)
}

func (c *Context) HandleError(err error) {
	switch {
	case validator.IsInvalid(err):
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	case errors.Is(err, policy.ErrDenied):
		_ = c.AbortWithError(http.StatusForbidden, err)
		return
	case errors.Is(err, storage.ErrNotFound):
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	default:
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (c *Context) NewUserPolicy(user domain.User) policy.IUserPolicy {
	return c.MustGetOpts().UserPolicy.Wrap(user)
}

func (c *Context) NewWorkspacePolicy(user domain.User, uw domain.UserWorkspace) policy.IWorkspacePolicy {
	return c.MustGetOpts().WorkspacePolicy.Wrap(user, uw)
}

func (c *Context) OtpIssue(user domain.User) (secret string, url string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:     "wault",
		Period:     30,
		SecretSize: 20,
		Digits:     otp.DigitsSix,
		// SHA1 is designed by spec
		Algorithm:   otp.AlgorithmSHA1,
		AccountName: user.ID.String,
	})

	if err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), err
}

func (c *Context) IsOtpValid(passcode string, secret string) bool {
	return totp.Validate(passcode, secret)
}
