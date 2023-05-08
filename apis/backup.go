package apis

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/forms"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

// bindBackupApi registers the file api endpoints and the corresponding handlers.
//
// @todo add hooks once the app hooks api restructuring is finalized
func bindBackupApi(app core.App, rg *echo.Group) {
	api := backupApi{app: app}

	subGroup := rg.Group("/backups", ActivityLogger(app))
	subGroup.GET("", api.list, RequireAdminAuth())
	subGroup.POST("", api.create, RequireAdminAuth())
	subGroup.GET("/:name", api.download)
	subGroup.DELETE("/:name", api.delete, RequireAdminAuth())
	subGroup.POST("/:name/restore", api.restore, RequireAdminAuth())
}

type backupApi struct {
	app core.App
}

type backupItem struct {
	Name     string         `json:"name"`
	Size     int64          `json:"size"`
	Modified types.DateTime `json:"modified"`
}

func (api *backupApi) list(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fsys, err := api.app.NewBackupsFilesystem()
	if err != nil {
		return NewBadRequestError("Failed to load backups filesystem", err)
	}
	defer fsys.Close()

	fsys.SetContext(ctx)

	backups, err := fsys.List("")
	if err != nil {
		return NewBadRequestError("Failed to retrieve backup items. Raw error: \n"+err.Error(), nil)
	}

	result := make([]backupItem, len(backups))

	for i, obj := range backups {
		modified, _ := types.ParseDateTime(obj.ModTime)

		result[i] = backupItem{
			Name:     obj.Key,
			Size:     obj.Size,
			Modified: modified,
		}
	}

	return c.JSON(http.StatusOK, result)
}

func (api *backupApi) create(c echo.Context) error {
	if cast.ToString(api.app.Cache().Get(core.CacheActiveBackupsKey)) != "" {
		return NewBadRequestError("Try again later - another backup/restore process has already been started", nil)
	}

	form := forms.NewBackupCreate(api.app)
	if err := c.Bind(form); err != nil {
		return NewBadRequestError("An error occurred while loading the submitted data.", err)
	}

	return form.Submit(func(next forms.InterceptorNextFunc[string]) forms.InterceptorNextFunc[string] {
		return func(name string) error {
			if err := next(name); err != nil {
				return NewBadRequestError("Failed to create backup", err)
			}

			return c.NoContent(http.StatusNoContent)
		}
	})
}

func (api *backupApi) download(c echo.Context) error {
	fileToken := c.QueryParam("token")

	_, err := api.app.Dao().FindAdminByToken(
		fileToken,
		api.app.Settings().AdminFileToken.Secret,
	)
	if err != nil {
		return NewForbiddenError("Insufficient permissions to access the resource.", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	fsys, err := api.app.NewBackupsFilesystem()
	if err != nil {
		return NewBadRequestError("Failed to load backups filesystem", err)
	}
	defer fsys.Close()

	fsys.SetContext(ctx)

	name := c.PathParam("name")

	br, err := fsys.GetFile(name)
	if err != nil {
		return NewBadRequestError("Failed to retrieve backup item. Raw error: \n"+err.Error(), nil)
	}
	defer br.Close()

	return fsys.Serve(
		c.Response(),
		c.Request(),
		name,
		filepath.Base(name), // without the path prefix (if any)
	)
}

func (api *backupApi) restore(c echo.Context) error {
	if cast.ToString(api.app.Cache().Get(core.CacheActiveBackupsKey)) != "" {
		return NewBadRequestError("Try again later - another backup/restore process has already been started", nil)
	}

	name := c.PathParam("name")

	existsCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fsys, err := api.app.NewBackupsFilesystem()
	if err != nil {
		return NewBadRequestError("Failed to load backups filesystem", err)
	}
	defer fsys.Close()

	fsys.SetContext(existsCtx)

	if exists, err := fsys.Exists(name); !exists {
		return NewNotFoundError("Missing or invalid backup file", err)
	}

	go func() {
		// wait max 10 minutes
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		// give some optimistic time to write the response
		time.Sleep(1 * time.Second)

		if err := api.app.RestoreBackup(ctx, name); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	}()

	return c.NoContent(http.StatusNoContent)
}

func (api *backupApi) delete(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fsys, err := api.app.NewBackupsFilesystem()
	if err != nil {
		return NewBadRequestError("Failed to load backups filesystem", err)
	}
	defer fsys.Close()

	fsys.SetContext(ctx)

	name := c.PathParam("name")

	if err := fsys.Delete(name); err != nil {
		return NewBadRequestError("Invalid or already deleted backup file. Raw error: \n"+err.Error(), nil)
	}

	return c.NoContent(http.StatusNoContent)
}
