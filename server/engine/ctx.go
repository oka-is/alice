package engine

import (
	"context"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/wault-pw/alice/lib/jwt"
	"github.com/wault-pw/alice/pkg/domain"
	"github.com/wault-pw/alice/pkg/pack"
	"github.com/wault-pw/alice/pkg/storage"
	"google.golang.org/protobuf/proto"
)

const (
	keyStore = string(rune(iota))
	keySession
	keyOpts
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
	c.SetCookie("jwt", token, 0, "/", opts.CookieDomain, opts.CookieSecure, true)
}

func (c *Context) GetCookieToken() (string, error) {
	return c.Cookie("jwt")
}

func (c *Context) JwtOpts() jwt.Opts {
	return jwt.Opts{
		Aud:  "WEB",
		Sub:  "API/V1",
		Iss:  "ALICE",
		Jti:  domain.NewUUID(),
		Algo: jwt.HS256,
		Key:  []byte{1, 2, 3, 4},
	}
}

func (c *Context) MustBindProto(m proto.Message) error {
	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	return proto.Unmarshal(buf, m)
}
