package apis

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tokens"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cast"
)

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

	handlerErr := api.app.OnFileBeforeTokenRequest().Trigger(event, func(e *core.FileTokenEvent) error {
		if e.Model == nil || e.Token == "" {
			return NewBadRequestError("Failed to generate file token.", nil)
		}

		return e.HttpContext.JSON(http.StatusOK, map[string]string{
			"token": e.Token,
		})
	})

	if handlerErr == nil {
		if err := api.app.OnFileAfterTokenRequest().Trigger(event); err != nil && api.app.IsDebug() {
			log.Println(err)
		}
	}

	return handlerErr
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

		if !api.canAccessRecord(adminOrAuthRecord, record, record.Collection().ViewRule) {
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

			// create a new thumb if it doesn exists
			if exists, _ := fs.Exists(servedPath); !exists {
				if err := fs.CreateThumb(originalPath, servedPath, thumbSize); err != nil {
					servedPath = originalPath // fallback to the original
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
		res := e.HttpContext.Response()
		req := e.HttpContext.Request()
		if err := fs.Serve(res, req, e.ServedPath, e.ServedName); err != nil {
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

// @todo move to a helper and maybe combine with the realtime checks when refactoring the realtime service
func (api *fileApi) canAccessRecord(adminOrAuthRecord models.Model, record *models.Record, accessRule *string) bool {
	admin, _ := adminOrAuthRecord.(*models.Admin)
	if admin != nil {
		// admins can access everything
		return true
	}

	if accessRule == nil {
		// only admins can access this record
		return false
	}

	ruleFunc := func(q *dbx.SelectQuery) error {
		if *accessRule == "" {
			return nil // empty public rule
		}

		// mock request data
		requestData := &models.RequestData{
			Method: "GET",
		}
		requestData.AuthRecord, _ = adminOrAuthRecord.(*models.Record)

		resolver := resolvers.NewRecordFieldResolver(api.app.Dao(), record.Collection(), requestData, true)
		expr, err := search.FilterData(*accessRule).BuildExpr(resolver)
		if err != nil {
			return err
		}
		resolver.UpdateQuery(q)
		q.AndWhere(expr)

		return nil
	}

	foundRecord, err := api.app.Dao().FindRecordById(record.Collection().Id, record.Id, ruleFunc)
	if err == nil && foundRecord != nil {
		return true
	}

	return false
}
