package engine

import (
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/wault-pw/alice/pkg/storage"
)

func New(store storage.IStore, opts Opts) (*Engine, error) {
	if opts.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	err := router.SetTrustedProxies([]string{})
	if err != nil {
		return nil, fmt.Errorf("failed to set trusted procies: %w", err)
	}

	router.NoRoute(WrapAction(noRoute))
	router.Use(useExtendedContext(store, opts))
	err = applySentry(router, opts.SentryDsn)
	if err != nil {
		return nil, fmt.Errorf("failed to apply sentry: %w", err)
	}

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

	return router, nil
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

// @refs https://docs.sentry.io/platforms/go/guides/gin/
func applySentry(app *gin.Engine, dsn string) error {
	if dsn == "" {
		return nil
	}

	err := sentry.Init(sentry.ClientOptions{Dsn: dsn, AttachStacktrace: true})
	if err != nil {
		return err
	}

	app.Use(sentrygin.New(sentrygin.Options{
		Repanic:         true,
		WaitForDelivery: false,
	}))

	return nil
}

func CaptureException(err error) {
	sentry.CaptureException(err)
}
