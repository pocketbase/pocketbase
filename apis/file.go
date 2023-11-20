package apis

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
	"golang.org/x/sync/semaphore"
)

// maxThumbConcurr specifies number of thumbnails that will be generated concurrently.
//
// It is advantages to limit concurrency to avoid high memory usage. This happens due to the
// fact that full sized images have to be kept in memory to be resized. In case many thumbnails
// are requested at the same time (e.g. listing collection in Web-UI) resizing many at the
// same time will slow down the resizing process, causing full sized images to be kept in memory
// even longer and making the problem even worse.
//
// The default value is 2x the number of CPUs. This is a pretty conservative value, assuming
// that it is better to have a slower thumbnail generation than to run out of memory.
// This value could be increased if high S3 latency is expected, full-size images are expected
// to be small and/or the server has plenty of memory.
var maxThumbConcurr = runtime.NumCPU() * 2
var thumbGenSem = semaphore.NewWeighted(int64(maxThumbConcurr))

// thumbGenTimeout specifies the maximum wait-time a request will be blocked waiting for the
// thumbnail to be generated if it does not exist already.
var thumbGenTimeout = 30 * time.Second

var imageContentTypes = []string{"image/png", "image/jpg", "image/jpeg", "image/gif"}
var defaultThumbSizes = []string{"100x100"}

// bindFileApi registers the file api endpoints and the corresponding handlers.
func bindFileApi(app core.App, rg *echo.Group) {
	api := fileApi{app: app}

	subGroup := rg.Group("/files", ActivityLogger(app))
	subGroup.POST("/token", api.fileToken)
	subGroup.HEAD("/:collection/:recordId/:filename", api.download, LoadCollectionContext(api.app))
	subGroup.GET("/:collection/:recordId/:filename", api.download, LoadCollectionContext(api.app))
}

type fileApi struct {
	app core.App
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
		return NewBadRequestError("", errors.New("Failed to load file options."))
	}

	// check whether the request is authorized to view the protected file
	if options.Protected {
		token := c.QueryParam("token")

		adminOrAuthRecord, _ := api.findAdminOrAuthRecordByFileToken(token)

		// create a copy of the cached request data and adjust it for the current auth model
		requestInfo := *RequestInfo(c)
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

	fs, err := api.app.NewFilesystem()
	if err != nil {
		return NewBadRequestError("Filesystem initialization failure.", err)
	}
	defer fs.Close()

	originalPath := baseFilesPath + "/" + filename
	servedPath := originalPath
	servedName := filename

	// check for valid thumb size param
	thumbSize := c.QueryParam("thumb")
	if thumbSize != "" && (list.ExistInSlice(thumbSize, defaultThumbSizes) || list.ExistInSlice(thumbSize, options.Thumbs)) {
		// extract the original file meta attributes and check it existence
		oAttrs, oAttrsErr := fs.Attributes(originalPath)
		if oAttrsErr != nil {
			return NewNotFoundError("", err)
		}

		// check if it is an image
		if list.ExistInSlice(oAttrs.ContentType, imageContentTypes) {
			// add thumb size as file suffix
			servedName = thumbSize + "_" + filename
			servedPath = baseFilesPath + "/thumbs_" + filename + "/" + servedName

			// create a new thumb if it doesn't exist (with limited concurrency since this is a very resource intensive operation)
			if exists, _ := fs.Exists(servedPath); !exists {
				ctx, cancel := context.WithTimeout(c.Request().Context(), thumbGenTimeout)
				defer cancel()
				err := thumbGenSem.Acquire(ctx, 1)
				if err != nil {
					return NewApiError(503, "Service Unavailable", errors.New("timeout error: thumbnail generation canceled due to high load, try again later"))
				}
				defer thumbGenSem.Release(1)

				if exists, _ := fs.Exists(servedPath); !exists { // check again if thumbnail was generated in the meantime already by another request
					if err := fs.CreateThumb(originalPath, servedPath, thumbSize); err != nil {
						servedPath = originalPath // fallback to the original
					}
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

		if err := fs.Serve(e.HttpContext.Response(), e.HttpContext.Request(), e.ServedPath, e.ServedName); err != nil {
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
