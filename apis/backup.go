package apis

import (
	"context"
	"net/http"
	"path/filepath"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/routine"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

// bindBackupApi registers the file api endpoints and the corresponding handlers.
func bindBackupApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	sub := rg.Group("/backups")
	sub.GET("", backupsList).Bind(RequireSuperuserAuth())
	sub.POST("", backupCreate).Bind(RequireSuperuserAuth())
	sub.POST("/upload", backupUpload).Bind(BodyLimit(0), RequireSuperuserAuth())
	sub.GET("/{key}", backupDownload) // relies on superuser file token
	sub.DELETE("/{key}", backupDelete).Bind(RequireSuperuserAuth())
	sub.POST("/{key}/restore", backupRestore).Bind(RequireSuperuserAuth())
}

type backupFileInfo struct {
	Modified types.DateTime `json:"modified"`
	Key      string         `json:"key"`
	Size     int64          `json:"size"`
}

func backupsList(e *core.RequestEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fsys, err := e.App.NewBackupsFilesystem()
	if err != nil {
		return e.BadRequestError("Failed to load backups filesystem.", err)
	}
	defer fsys.Close()

	fsys.SetContext(ctx)

	backups, err := fsys.List("")
	if err != nil {
		return e.BadRequestError("Failed to retrieve backup items. Raw error: \n"+err.Error(), nil)
	}

	result := make([]backupFileInfo, len(backups))

	for i, obj := range backups {
		modified, _ := types.ParseDateTime(obj.ModTime)

		result[i] = backupFileInfo{
			Key:      obj.Key,
			Size:     obj.Size,
			Modified: modified,
		}
	}

	return e.JSON(http.StatusOK, result)
}

func backupDownload(e *core.RequestEvent) error {
	fileToken := e.Request.URL.Query().Get("token")

	authRecord, err := e.App.FindAuthRecordByToken(fileToken, core.TokenTypeFile)
	if err != nil || !authRecord.IsSuperuser() {
		return e.ForbiddenError("Insufficient permissions to access the resource.", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	fsys, err := e.App.NewBackupsFilesystem()
	if err != nil {
		return e.InternalServerError("Failed to load backups filesystem.", err)
	}
	defer fsys.Close()

	fsys.SetContext(ctx)

	key := e.Request.PathValue("key")

	return fsys.Serve(
		e.Response,
		e.Request,
		key,
		filepath.Base(key), // without the path prefix (if any)
	)
}

func backupDelete(e *core.RequestEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fsys, err := e.App.NewBackupsFilesystem()
	if err != nil {
		return e.InternalServerError("Failed to load backups filesystem.", err)
	}
	defer fsys.Close()

	fsys.SetContext(ctx)

	key := e.Request.PathValue("key")

	if key != "" && cast.ToString(e.App.Store().Get(core.StoreKeyActiveBackup)) == key {
		return e.BadRequestError("The backup is currently being used and cannot be deleted.", nil)
	}

	if err := fsys.Delete(key); err != nil {
		return e.BadRequestError("Invalid or already deleted backup file. Raw error: \n"+err.Error(), nil)
	}

	return e.NoContent(http.StatusNoContent)
}

func backupRestore(e *core.RequestEvent) error {
	if e.App.Store().Has(core.StoreKeyActiveBackup) {
		return e.BadRequestError("Try again later - another backup/restore process has already been started.", nil)
	}

	key := e.Request.PathValue("key")

	existsCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fsys, err := e.App.NewBackupsFilesystem()
	if err != nil {
		return e.InternalServerError("Failed to load backups filesystem.", err)
	}
	defer fsys.Close()

	fsys.SetContext(existsCtx)

	if exists, err := fsys.Exists(key); !exists {
		return e.BadRequestError("Missing or invalid backup file.", err)
	}

	routine.FireAndForget(func() {
		// give some optimistic time to write the response before restarting the app
		time.Sleep(1 * time.Second)

		// wait max 10 minutes to fetch the backup
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if err := e.App.RestoreBackup(ctx, key); err != nil {
			e.App.Logger().Error("Failed to restore backup", "key", key, "error", err.Error())
		}
	})

	return e.NoContent(http.StatusNoContent)
}
