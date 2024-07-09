package apis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
	"golang.org/x/sync/semaphore"
	"golang.org/x/sync/singleflight"
)

var imageContentTypes = []string{"image/png", "image/jpg", "image/jpeg", "image/gif"}
var defaultThumbSizes = []string{"100x100"}

// bindFileApi registers the file api endpoints and the corresponding handlers.
func bindFileApi(app core.App, rg *echo.Group) {
	api := fileApi{
		app:             app,
		thumbGenSem:     semaphore.NewWeighted(int64(runtime.NumCPU() + 2)), // the value is arbitrary chosen and may change in the future
		thumbGenPending: new(singleflight.Group),
		thumbGenMaxWait: 60 * time.Second,
	}

	subGroup := rg.Group("/files", ActivityLogger(app))
	subGroup.POST("/token", api.fileToken)
	subGroup.HEAD("/:collection/:recordId/:filename", api.download, LoadCollectionContext(api.app))
	subGroup.GET("/:collection/:recordId/:filename", api.download, LoadCollectionContext(api.app))
}

type fileApi struct {
	app core.App

	// thumbGenSem is a semaphore to prevent too much concurrent
	// requests generating new thumbs at the same time.
	thumbGenSem *semaphore.Weighted

	// thumbGenPending represents a group of currently pending
	// thumb generation processes.
	thumbGenPending *singleflight.Group

	// thumbGenMaxWait is the maximum waiting time for starting a new
	// thumb generation process.
	thumbGenMaxWait time.Duration
}

func (api *fileApi) fileToken(c echo.Context) error {
	event := new(core.FileTokenEvent)
	event.HttpContext = c

	if admin, _ := c.Get(ContextAdminKey).(*models.Admin); admin != nil {
		event.Model = admin
		event.Token, _ = tokens.NewAdminFileToken(api.app, admin)
	} else if record, _ := c.Get(ContextAuthRecordKey).(*models.Record); record != nil {
		event.Model = record
		event.Token, _ = tokens.NewRecordFileToken(api.app, record)
	}

	return api.app.OnFileBeforeTokenRequest().Trigger(event, func(e *core.FileTokenEvent) error {
		if e.Model == nil || e.Token == "" {
			return NewBadRequestError("Failed to generate file token.", nil)
		}

		return api.app.OnFileAfterTokenRequest().Trigger(event, func(e *core.FileTokenEvent) error {
			if e.HttpContext.Response().Committed {
				return nil
			}

			return e.HttpContext.JSON(http.StatusOK, map[string]string{
				"token": e.Token,
			})
		})
	})
}

func (api *fileApi) download(c echo.Context) error {
	collection, _ := c.Get(ContextCollectionKey).(*models.Collection)
	if collection == nil {
		return NewNotFoundError("", nil)
	}

	recordId := c.PathParam("recordId")
	if recordId == "" {
		return NewNotFoundError("", nil)
	}

	record, err := api.app.Dao().FindRecordById(collection.Id, recordId)
	if err != nil {
		return NewNotFoundError("", err)
	}

	filename := c.PathParam("filename")

	fileField := record.FindFileFieldByFile(filename)
	if fileField == nil {
		return NewNotFoundError("", nil)
	}

	options, ok := fileField.Options.(*schema.FileOptions)
	if !ok {
		return NewBadRequestError("", errors.New("failed to load file options"))
	}

	// check whether the request is authorized to view the protected file
	if options.Protected {
		token := c.QueryParam("token")

		adminOrAuthRecord, _ := api.findAdminOrAuthRecordByFileToken(token)

		// create a copy of the cached request data and adjust it for the current auth model
		requestInfo := *RequestInfo(c)
		requestInfo.Context = models.RequestInfoContextProtectedFile
		requestInfo.Admin = nil
		requestInfo.AuthRecord = nil
		if adminOrAuthRecord != nil {
			if admin, _ := adminOrAuthRecord.(*models.Admin); admin != nil {
				requestInfo.Admin = admin
			} else if record, _ := adminOrAuthRecord.(*models.Record); record != nil {
				requestInfo.AuthRecord = record
			}
		}

		if ok, _ := api.app.Dao().CanAccessRecord(record, &requestInfo, record.Collection().ViewRule); !ok {
			return NewForbiddenError("Insufficient permissions to access the file resource.", nil)
		}
	}

	baseFilesPath := record.BaseFilesPath()

	// fetch the original view file field related record
	if collection.IsView() {
		fileRecord, err := api.app.Dao().FindRecordByViewFile(collection.Id, fileField.Name, filename)
		if err != nil {
			return NewNotFoundError("", fmt.Errorf("Failed to fetch view file field record: %w", err))
		}
		baseFilesPath = fileRecord.BaseFilesPath()
	}

	fsys, err := api.app.NewFilesystem()
	if err != nil {
		return NewBadRequestError("Filesystem initialization failure.", err)
	}
	defer fsys.Close()

	originalPath := baseFilesPath + "/" + filename
	servedPath := originalPath
	servedName := filename

	// check for valid thumb size param
	thumbSize := c.QueryParam("thumb")
	if thumbSize != "" && (list.ExistInSlice(thumbSize, defaultThumbSizes) || list.ExistInSlice(thumbSize, options.Thumbs)) {
		// extract the original file meta attributes and check it existence
		oAttrs, oAttrsErr := fsys.Attributes(originalPath)
		if oAttrsErr != nil {
			return NewNotFoundError("", err)
		}

		// check if it is an image
		if list.ExistInSlice(oAttrs.ContentType, imageContentTypes) {
			// add thumb size as file suffix
			servedName = thumbSize + "_" + filename
			servedPath = baseFilesPath + "/thumbs_" + filename + "/" + servedName

			// create a new thumb if it doesn't exist
			if exists, _ := fsys.Exists(servedPath); !exists {
				if err := api.createThumb(c, fsys, originalPath, servedPath, thumbSize); err != nil {
					api.app.Logger().Warn(
						"Fallback to original - failed to create thumb "+servedName,
						slog.Any("error", err),
						slog.String("original", originalPath),
						slog.String("thumb", servedPath),
					)

					// fallback to the original
					servedName = filename
					servedPath = originalPath
				}
			}
		}
	}

	event := new(core.FileDownloadEvent)
	event.HttpContext = c
	event.Collection = collection
	event.Record = record
	event.FileField = fileField
	event.ServedPath = servedPath
	event.ServedName = servedName

	// clickjacking shouldn't be a concern when serving uploaded files,
	// so it safe to unset the global X-Frame-Options to allow files embedding
	// (note: it is out of the hook to allow users to customize the behavior)
	c.Response().Header().Del("X-Frame-Options")

	return api.app.OnFileDownloadRequest().Trigger(event, func(e *core.FileDownloadEvent) error {
		if e.HttpContext.Response().Committed {
			return nil
		}

		if err := fsys.Serve(e.HttpContext.Response(), e.HttpContext.Request(), e.ServedPath, e.ServedName); err != nil {
			return NewNotFoundError("", err)
		}

		return nil
	})
}

func (api *fileApi) findAdminOrAuthRecordByFileToken(fileToken string) (models.Model, error) {
	fileToken = strings.TrimSpace(fileToken)
	if fileToken == "" {
		return nil, errors.New("missing file token")
	}

	claims, _ := security.ParseUnverifiedJWT(strings.TrimSpace(fileToken))
	tokenType := cast.ToString(claims["type"])

	switch tokenType {
	case tokens.TypeAdmin:
		admin, err := api.app.Dao().FindAdminByToken(
			fileToken,
			api.app.Settings().AdminFileToken.Secret,
		)
		if err == nil && admin != nil {
			return admin, nil
		}
	case tokens.TypeAuthRecord:
		record, err := api.app.Dao().FindAuthRecordByToken(
			fileToken,
			api.app.Settings().RecordFileToken.Secret,
		)
		if err == nil && record != nil {
			return record, nil
		}
	}

	return nil, errors.New("missing or invalid file token")
}

func (api *fileApi) createThumb(
	c echo.Context,
	fsys *filesystem.System,
	originalPath string,
	thumbPath string,
	thumbSize string,
) error {
	ch := api.thumbGenPending.DoChan(thumbPath, func() (any, error) {
		ctx, cancel := context.WithTimeout(c.Request().Context(), api.thumbGenMaxWait)
		defer cancel()

		if err := api.thumbGenSem.Acquire(ctx, 1); err != nil {
			return nil, err
		}
		defer api.thumbGenSem.Release(1)

		return nil, fsys.CreateThumb(originalPath, thumbPath, thumbSize)
	})

	res := <-ch

	api.thumbGenPending.Forget(thumbPath)

	return res.Err
}
