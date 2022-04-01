package api_v1

import (
	"fmt"
	"io"
	"net/http"

	"github.com/wault-pw/alice/pkg/backup"
	"github.com/wault-pw/alice/server/engine"
)

func CreateBackup(ctx *engine.Context) {
	flusher := ctx.Writer.(http.Flusher)

	ctx.Header("content-type", "text/html; charset=UTF-8")
	ctx.Header("content-disposition", fmt.Sprintf(`attachment; filename="%s"`, "wault.html"))

	err := before(ctx, flusher)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	err = dump(ctx, flusher)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	err = after(ctx, flusher)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	err = pre(ctx, flusher)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
	}

	ctx.Done()
}

func pre(ctx *engine.Context, flusher http.Flusher) error {
	resp, err := http.Get(ctx.MustGetOpts().BackupUrl)
	if err != nil {
		return fmt.Errorf("failed to get a backup url: %w", err)
	}

	defer resp.Body.Close()

	_, err = io.Copy(ctx.Writer, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to copy a backup: %w", err)
	}

	flusher.Flush()
	return nil
}

func before(ctx *engine.Context, flusher http.Flusher) error {
	_, err := ctx.Writer.Write([]byte("<script>window.__DUMP__ = new Uint8Array(["))
	flusher.Flush()
	return err
}

func after(ctx *engine.Context, flusher http.Flusher) error {
	_, err := ctx.Writer.Write([]byte("]);</script>"))
	flusher.Flush()
	return err
}

func dump(ctx *engine.Context, flusher http.Flusher) error {
	err := backup.NewJsBackup(ctx.GetStore(), ctx.Writer, flusher).
		DumpV1(ctx.MustGetSession().UserID.String)
	if err != nil {
		return err
	}

	return nil
}
