package engine

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wault-pw/alice/pkg/pack"
	"github.com/wault-pw/alice/pkg/storage"
)

type Opts struct {
	AllowOrigin  []string
	CookieDomain string
	CookieSecure bool
	BackupUrl    string
	Ver          *pack.Ver
}

func New(store storage.IStore, opts Opts) *Engine {
	router := gin.Default()
	router.NoRoute(WrapAction(noRoute))
	router.Use(useExtendedContext(store, opts))

	router.Use(cors.New(cors.Config{
		AllowOrigins: opts.AllowOrigin,
		AllowMethods: []string{"PUT", "PATCH", "POST", "UPGRADE", "GET", "UPDATE"},
		AllowHeaders: []string{
			"Access-Control-Allow-Headers",
			"Origin",
			"Accept",
			"X-Requested-With",
			"Content-Type",
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
		},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/health", WrapAction(health))

	return router
}

func WrapAction(action ActionFN) gin.HandlerFunc {
	return func(c *gin.Context) {
		action(Ctx(c))
	}
}

func useExtendedContext(store storage.IStore, opts Opts) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		Ctx(ctx).SetStore(store).SetOpts(opts)
		ctx.Next()
	}
}

func health(ctx *Context) {
	if err := ctx.GetStore().Ping(ctx.Ctx()); err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	} else {
		ctx.String(http.StatusOK, "[OK]")
	}
}
