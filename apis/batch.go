package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/router"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cast"
)

func bindBatchApi(app core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	sub := rg.Group("/batch")
	sub.POST("", batchTransaction).Unbind(DefaultBodyLimitMiddlewareId) // the body limit is inlined
}

type HandleFunc func(e *core.RequestEvent) error

type BatchActionHandlerFunc func(app core.App, ir *core.InternalRequest, params map[string]string, next func(data any) error) HandleFunc

// ValidBatchActions defines a map with the supported batch InternalRequest actions.
//
// Note: when adding new routes make sure that their middlewares are inlined!
var ValidBatchActions = map[*regexp.Regexp]BatchActionHandlerFunc{
	// "upsert" handler
	regexp.MustCompile(`^PUT /api/collections/(?P<collection>[^\/\?]+)/records(?P<query>\?.*)?$`): func(app core.App, ir *core.InternalRequest, params map[string]string, next func(any) error) HandleFunc {
		var id string
		if len(ir.Body) > 0 && ir.Body["id"] != "" {
			id = cast.ToString(ir.Body["id"])
		}
		if id != "" {
			_, err := app.FindRecordById(params["collection"], id)
			if err == nil {
				// update
				// ---
				params["id"] = id // required for the path value
				ir.Method = "PATCH"
				ir.URL = "/api/collections/" + params["collection"] + "/records/" + id + params["query"]
				return recordUpdate(next)
			}
		}

		// create
		// ---
		ir.Method = "POST"
		ir.URL = "/api/collections/" + params["collection"] + "/records" + params["query"]
		return recordCreate(next)
	},
	regexp.MustCompile(`^POST /api/collections/(?P<collection>[^\/\?]+)/records(\?.*)?$`): func(app core.App, ir *core.InternalRequest, params map[string]string, next func(any) error) HandleFunc {
		return recordCreate(next)
	},
	regexp.MustCompile(`^PATCH /api/collections/(?P<collection>[^\/\?]+)/records/(?P<id>[^\/\?]+)(\?.*)?$`): func(app core.App, ir *core.InternalRequest, params map[string]string, next func(any) error) HandleFunc {
		return recordUpdate(next)
	},
	regexp.MustCompile(`^DELETE /api/collections/(?P<collection>[^\/\?]+)/records/(?P<id>[^\/\?]+)(\?.*)?$`): func(app core.App, ir *core.InternalRequest, params map[string]string, next func(any) error) HandleFunc {
		return recordDelete(next)
	},
}

type BatchRequestResult struct {
	Body   any `json:"body"`
	Status int `json:"status"`
}

type batchRequestsForm struct {
	Requests []*core.InternalRequest `form:"requests" json:"requests"`

	max int
}

func (brs batchRequestsForm) validate() error {
	return validation.ValidateStruct(&brs,
		validation.Field(&brs.Requests, validation.Required, validation.Length(0, brs.max)),
	)
}

// NB! When the request is submitted as multipart/form-data,
// the regular fields data is expected to be submitted as serailized
// json under the @jsonPayload field and file keys need to follow the
// pattern "requests.N.fileField" or  requests[N].fileField.
func batchTransaction(e *core.RequestEvent) error {
	maxRequests := e.App.Settings().Batch.MaxRequests
	if !e.App.Settings().Batch.Enabled || maxRequests <= 0 {
		return e.ForbiddenError("Batch requests are not allowed.", nil)
	}

	txTimeout := time.Duration(e.App.Settings().Batch.Timeout) * time.Second
	if txTimeout <= 0 {
		txTimeout = 3 * time.Second // for now always limit
	}

	maxBodySize := e.App.Settings().Batch.MaxBodySize
	if maxBodySize <= 0 {
		maxBodySize = 128 << 20
	}

	err := applyBodyLimit(e, maxBodySize)
	if err != nil {
		return err
	}

	form := &batchRequestsForm{max: maxRequests}

	// load base requests data
	err = e.BindBody(form)
	if err != nil {
		return e.BadRequestError("Failed to read the submitted batch data.", err)
	}

	// load uploaded files into each request item
	// note: expects the files to be under "requests.N.fileField" or "requests[N].fileField" format
	// 		 (the other regular fields must be put under `@jsonPayload` as serialized json)
	if strings.HasPrefix(e.Request.Header.Get("Content-Type"), "multipart/form-data") {
		for i, ir := range form.Requests {
			iStr := strconv.Itoa(i)

			files, err := extractPrefixedFiles(e.Request, "requests."+iStr+".", "requests["+iStr+"].")
			if err != nil {
				return e.BadRequestError("Failed to read the submitted batch files data.", err)
			}

			for key, files := range files {
				if ir.Body == nil {
					ir.Body = map[string]any{}
				}
				ir.Body[key] = files
			}
		}
	}

	// validate batch request form
	err = form.validate()
	if err != nil {
		return e.BadRequestError("Invalid batch request data.", err)
	}

	event := new(core.BatchRequestEvent)
	event.RequestEvent = e
	event.Batch = form.Requests

	return e.App.OnBatchRequest().Trigger(event, func(e *core.BatchRequestEvent) error {
		bp := batchProcessor{
			app:         e.App,
			baseEvent:   e.RequestEvent,
			infoContext: core.RequestInfoContextBatch,
		}

		if err := bp.Process(e.Batch, txTimeout); err != nil {
			return firstApiError(err, e.BadRequestError("Batch transaction failed.", err))
		}

		return e.JSON(http.StatusOK, bp.results)
	})
}

type batchProcessor struct {
	app         core.App
	baseEvent   *core.RequestEvent
	infoContext string
	results     []*BatchRequestResult
	failedIndex int
	errCh       chan error
	stopCh      chan struct{}
}

func (p *batchProcessor) Process(batch []*core.InternalRequest, timeout time.Duration) error {
	p.results = make([]*BatchRequestResult, 0, len(batch))

	if p.stopCh != nil {
		close(p.stopCh)
	}
	p.stopCh = make(chan struct{}, 1)

	if p.errCh != nil {
		close(p.errCh)
	}
	p.errCh = make(chan error, 1)

	return p.app.RunInTransaction(func(txApp core.App) error {
		// used to interupts the recursive processing calls in case of a timeout or connection close
		defer func() {
			p.stopCh <- struct{}{}
		}()

		go func() {
			err := p.process(txApp, batch, 0)

			if err != nil {
				err = validation.Errors{
					"requests": validation.Errors{
						strconv.Itoa(p.failedIndex): &BatchResponseError{
							code:    "batch_request_failed",
							message: "Batch request failed.",
							err:     router.ToApiError(err),
						},
					},
				}
			}

			// note: to avoid copying and due to the process recursion the final results order is reversed
			if err == nil {
				slices.Reverse(p.results)
			}

			p.errCh <- err
		}()

		select {
		case responseErr := <-p.errCh:
			return responseErr
		case <-time.After(timeout):
			// note: we don't return 408 Reques Timeout error because
			// some browsers perform automatic retry behind the scenes
			// which are hard to debug and unnecessary
			return errors.New("batch transaction timeout")
		case <-p.baseEvent.Request.Context().Done():
			return errors.New("batch request interrupted")
		}
	})
}

func (p *batchProcessor) process(activeApp core.App, batch []*core.InternalRequest, i int) error {
	select {
	case <-p.stopCh:
		return nil
	default:
		if len(batch) == 0 {
			return nil
		}

		result, err := processInternalRequest(
			activeApp,
			p.baseEvent,
			batch[0],
			p.infoContext,
			func(_ any) error {
				if len(batch) == 1 {
					return nil
				}

				err := p.process(activeApp, batch[1:], i+1)

				// update the failed batch index (if not already)
				if err != nil && p.failedIndex == 0 {
					p.failedIndex = i + 1
				}

				return err
			},
		)

		if err != nil {
			return err
		}

		p.results = append(p.results, result)

		return nil
	}
}

func processInternalRequest(
	activeApp core.App,
	baseEvent *core.RequestEvent,
	ir *core.InternalRequest,
	infoContext string,
	optNext func(data any) error,
) (*BatchRequestResult, error) {
	handle, params, ok := prepareInternalAction(activeApp, ir, optNext)
	if !ok {
		return nil, errors.New("unknown batch request action")
	}

	// construct a new http.Request
	// ---------------------------------------------------------------
	buf, mw, err := multipartDataFromInternalRequest(ir)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(strings.ToUpper(ir.Method), ir.URL, buf)
	if err != nil {
		return nil, err
	}

	// cleanup multipart temp files
	defer func() {
		if r.MultipartForm != nil {
			if err := r.MultipartForm.RemoveAll(); err != nil {
				activeApp.Logger().Warn("failed to cleanup temp batch files", "error", err)
			}
		}
	}()

	// load batch request path params
	// ---
	for k, v := range params {
		r.SetPathValue(k, v)
	}

	// clone original request
	// ---
	r.RequestURI = r.URL.RequestURI()
	r.Proto = baseEvent.Request.Proto
	r.ProtoMajor = baseEvent.Request.ProtoMajor
	r.ProtoMinor = baseEvent.Request.ProtoMinor
	r.Host = baseEvent.Request.Host
	r.RemoteAddr = baseEvent.Request.RemoteAddr
	r.TLS = baseEvent.Request.TLS

	if s := baseEvent.Request.TransferEncoding; s != nil {
		s2 := make([]string, len(s))
		copy(s2, s)
		r.TransferEncoding = s2
	}

	if baseEvent.Request.Trailer != nil {
		r.Trailer = baseEvent.Request.Trailer.Clone()
	}

	if baseEvent.Request.Header != nil {
		r.Header = baseEvent.Request.Header.Clone()
	}

	// apply batch request specific headers
	// ---
	for k, v := range ir.Headers {
		// individual Authorization header keys don't have affect
		// because the auth state is populated from the base event
		if strings.EqualFold(k, "authorization") {
			continue
		}

		r.Header.Set(k, v)
	}
	r.Header.Set("Content-Type", mw.FormDataContentType())

	// construct a new RequestEvent
	// ---------------------------------------------------------------
	event := &core.RequestEvent{}
	event.App = activeApp
	event.Auth = baseEvent.Auth
	event.SetAll(baseEvent.GetAll())

	// load RequestInfo context
	if infoContext == "" {
		infoContext = core.RequestInfoContextDefault
	}
	event.Set(core.RequestEventKeyInfoContext, infoContext)

	// assign request
	event.Request = r
	event.Request.Body = &router.RereadableReadCloser{ReadCloser: r.Body} // enables multiple reads

	// assign response
	rec := httptest.NewRecorder()
	event.Response = &router.ResponseWriter{ResponseWriter: rec} // enables status and write tracking

	// execute
	// ---------------------------------------------------------------
	if err := handle(event); err != nil {
		return nil, err
	}

	result := rec.Result()
	defer result.Body.Close()

	body, _ := types.ParseJSONRaw(rec.Body.Bytes())

	return &BatchRequestResult{
		Status: result.StatusCode,
		Body:   body,
	}, nil
}

func multipartDataFromInternalRequest(ir *core.InternalRequest) (*bytes.Buffer, *multipart.Writer, error) {
	buf := &bytes.Buffer{}

	mw := multipart.NewWriter(buf)

	regularFields := map[string]any{}
	fileFields := map[string][]*filesystem.File{}

	// separate regular fields from files
	// ---
	for k, rawV := range ir.Body {
		switch v := rawV.(type) {
		case *filesystem.File:
			fileFields[k] = append(fileFields[k], v)
		case []*filesystem.File:
			fileFields[k] = append(fileFields[k], v...)
		default:
			regularFields[k] = v
		}
	}

	// submit regularFields as @jsonPayload
	// ---
	rawBody, err := json.Marshal(regularFields)
	if err != nil {
		return nil, nil, errors.Join(err, mw.Close())
	}

	jsonPayload, err := mw.CreateFormField("@jsonPayload")
	if err != nil {
		return nil, nil, errors.Join(err, mw.Close())
	}
	_, err = jsonPayload.Write(rawBody)
	if err != nil {
		return nil, nil, errors.Join(err, mw.Close())
	}

	// submit fileFields as multipart files
	// ---
	for key, files := range fileFields {
		for _, file := range files {
			part, err := mw.CreateFormFile(key, file.Name)
			if err != nil {
				return nil, nil, errors.Join(err, mw.Close())
			}

			fr, err := file.Reader.Open()
			if err != nil {
				return nil, nil, errors.Join(err, mw.Close())
			}

			_, err = io.Copy(part, fr)
			if err != nil {
				return nil, nil, errors.Join(err, fr.Close(), mw.Close())
			}

			err = fr.Close()
			if err != nil {
				return nil, nil, errors.Join(err, mw.Close())
			}
		}
	}

	return buf, mw, mw.Close()
}

func extractPrefixedFiles(request *http.Request, prefixes ...string) (map[string][]*filesystem.File, error) {
	if request.MultipartForm == nil {
		if err := request.ParseMultipartForm(router.DefaultMaxMemory); err != nil {
			return nil, err
		}
	}

	result := make(map[string][]*filesystem.File)

	for k, fhs := range request.MultipartForm.File {
		for _, p := range prefixes {
			if strings.HasPrefix(k, p) {
				resultKey := strings.TrimPrefix(k, p)

				for _, fh := range fhs {
					file, err := filesystem.NewFileFromMultipart(fh)
					if err != nil {
						return nil, err
					}

					result[resultKey] = append(result[resultKey], file)
				}
			}
		}
	}

	return result, nil
}

func prepareInternalAction(activeApp core.App, ir *core.InternalRequest, optNext func(data any) error) (HandleFunc, map[string]string, bool) {
	full := strings.ToUpper(ir.Method) + " " + ir.URL

	for re, actionFactory := range ValidBatchActions {
		params, ok := findNamedMatches(re, full)
		if ok {
			return actionFactory(activeApp, ir, params, optNext), params, true
		}
	}

	return nil, nil, false
}

func findNamedMatches(re *regexp.Regexp, str string) (map[string]string, bool) {
	match := re.FindStringSubmatch(str)
	if match == nil {
		return nil, false
	}

	result := map[string]string{}

	names := re.SubexpNames()

	for i, m := range match {
		if names[i] != "" {
			result[names[i]] = m
		}
	}

	return result, true
}

// -------------------------------------------------------------------

var (
	_ router.SafeErrorItem     = (*BatchResponseError)(nil)
	_ router.SafeErrorResolver = (*BatchResponseError)(nil)
)

type BatchResponseError struct {
	err     *router.ApiError
	code    string
	message string
}

func (e *BatchResponseError) Error() string {
	return e.message
}

func (e *BatchResponseError) Code() string {
	return e.code
}

func (e *BatchResponseError) Resolve(errData map[string]any) any {
	errData["response"] = e.err
	return errData
}

func (e BatchResponseError) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"message":  e.message,
		"code":     e.code,
		"response": e.err,
	})
}
