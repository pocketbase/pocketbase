package router_test

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/tools/router"
)

type unwrapTester struct {
	http.ResponseWriter
}

func (ut unwrapTester) Unwrap() http.ResponseWriter {
	return ut.ResponseWriter
}

func TestEventWritten(t *testing.T) {
	t.Parallel()

	res1 := httptest.NewRecorder()

	res2 := httptest.NewRecorder()
	res2.Write([]byte("test"))

	res3 := &router.ResponseWriter{ResponseWriter: unwrapTester{httptest.NewRecorder()}}

	res4 := &router.ResponseWriter{ResponseWriter: unwrapTester{httptest.NewRecorder()}}
	res4.Write([]byte("test"))

	scenarios := []struct {
		name     string
		response http.ResponseWriter
		expected bool
	}{
		{
			name:     "non-written non-WriteTracker",
			response: res1,
			expected: false,
		},
		{
			name:     "written non-WriteTracker",
			response: res2,
			expected: false,
		},
		{
			name:     "non-written WriteTracker",
			response: res3,
			expected: false,
		},
		{
			name:     "written WriteTracker",
			response: res4,
			expected: true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			event := router.Event{
				Response: s.response,
			}

			result := event.Written()

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestEventStatus(t *testing.T) {
	t.Parallel()

	res1 := httptest.NewRecorder()

	res2 := httptest.NewRecorder()
	res2.WriteHeader(123)

	res3 := &router.ResponseWriter{ResponseWriter: unwrapTester{httptest.NewRecorder()}}

	res4 := &router.ResponseWriter{ResponseWriter: unwrapTester{httptest.NewRecorder()}}
	res4.WriteHeader(123)

	scenarios := []struct {
		name     string
		response http.ResponseWriter
		expected int
	}{
		{
			name:     "non-written non-StatusTracker",
			response: res1,
			expected: 0,
		},
		{
			name:     "written non-StatusTracker",
			response: res2,
			expected: 0,
		},
		{
			name:     "non-written StatusTracker",
			response: res3,
			expected: 0,
		},
		{
			name:     "written StatusTracker",
			response: res4,
			expected: 123,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			event := router.Event{
				Response: s.response,
			}

			result := event.Status()

			if result != s.expected {
				t.Fatalf("Expected %d, got %d", s.expected, result)
			}
		})
	}
}

func TestEventIsTLS(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	event := router.Event{Request: req}

	// without TLS
	if event.IsTLS() {
		t.Fatalf("Expected IsTLS false")
	}

	// dummy TLS state
	req.TLS = new(tls.ConnectionState)

	// with TLS
	if !event.IsTLS() {
		t.Fatalf("Expected IsTLS true")
	}
}

func TestEventSetCookie(t *testing.T) {
	t.Parallel()

	event := router.Event{
		Response: httptest.NewRecorder(),
	}

	cookie := event.Response.Header().Get("set-cookie")
	if cookie != "" {
		t.Fatalf("Expected empty cookie string, got %q", cookie)
	}

	event.SetCookie(&http.Cookie{Name: "test", Value: "a"})

	expected := "test=a"

	cookie = event.Response.Header().Get("set-cookie")
	if cookie != expected {
		t.Fatalf("Expected cookie %q, got %q", expected, cookie)
	}
}

func TestEventRemoteIP(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		remoteAddr string
		expected   string
	}{
		{"", "invalid IP"},
		{"1.2.3.4", "invalid IP"},
		{"1.2.3.4:8090", "1.2.3.4"},
		{"[0000:0000:0000:0000:0000:0000:0000:0002]:80", "0000:0000:0000:0000:0000:0000:0000:0002"},
		{"[::2]:80", "0000:0000:0000:0000:0000:0000:0000:0002"}, // should always return the expanded version
	}

	for _, s := range scenarios {
		t.Run(s.remoteAddr, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.RemoteAddr = s.remoteAddr

			event := router.Event{Request: req}

			ip := event.RemoteIP()

			if ip != s.expected {
				t.Fatalf("Expected IP %q, got %q", s.expected, ip)
			}
		})
	}
}

func TestFindUploadedFiles(t *testing.T) {
	scenarios := []struct {
		filename        string
		expectedPattern string
	}{
		{"ab.png", `^ab\w{10}_\w{10}\.png$`},
		{"test", `^test_\w{10}\.txt$`},
		{"a b c d!@$.j!@$pg", `^a_b_c_d_\w{10}\.jpg$`},
		{strings.Repeat("a", 150), `^a{100}_\w{10}\.txt$`},
	}

	for _, s := range scenarios {
		t.Run(s.filename, func(t *testing.T) {
			// create multipart form file body
			body := new(bytes.Buffer)
			mp := multipart.NewWriter(body)
			w, err := mp.CreateFormFile("test", s.filename)
			if err != nil {
				t.Fatal(err)
			}
			w.Write([]byte("test"))
			mp.Close()
			// ---

			req := httptest.NewRequest(http.MethodPost, "/", body)
			req.Header.Add("Content-Type", mp.FormDataContentType())

			event := router.Event{Request: req}

			result, err := event.FindUploadedFiles("test")
			if err != nil {
				t.Fatal(err)
			}

			if len(result) != 1 {
				t.Fatalf("Expected 1 file, got %d", len(result))
			}

			if result[0].Size != 4 {
				t.Fatalf("Expected the file size to be 4 bytes, got %d", result[0].Size)
			}

			pattern, err := regexp.Compile(s.expectedPattern)
			if err != nil {
				t.Fatalf("Invalid filename pattern %q: %v", s.expectedPattern, err)
			}
			if !pattern.MatchString(result[0].Name) {
				t.Fatalf("Expected filename to match %s, got filename %s", s.expectedPattern, result[0].Name)
			}
		})
	}
}

func TestFindUploadedFilesMissing(t *testing.T) {
	body := new(bytes.Buffer)
	mp := multipart.NewWriter(body)
	mp.Close()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Add("Content-Type", mp.FormDataContentType())

	event := router.Event{Request: req}

	result, err := event.FindUploadedFiles("test")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result != nil {
		t.Errorf("Expected result to be nil, got %v", result)
	}
}

func TestEventSetGet(t *testing.T) {
	event := router.Event{}

	// get before any set (ensures that doesn't panic)
	if v := event.Get("test"); v != nil {
		t.Fatalf("Expected nil value, got %v", v)
	}

	event.Set("a", 123)
	event.Set("b", 456)

	scenarios := []struct {
		key      string
		expected any
	}{
		{"", nil},
		{"missing", nil},
		{"a", 123},
		{"b", 456},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s", i, s.key), func(t *testing.T) {
			result := event.Get(s.key)
			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}

func TestEventSetAllGetAll(t *testing.T) {
	data := map[string]any{
		"a": 123,
		"b": 456,
	}
	rawData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	event := router.Event{}
	event.SetAll(data)

	// modify the data to ensure that the map was shallow coppied
	data["c"] = 789

	result := event.GetAll()
	rawResult, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}

	if len(rawResult) == 0 || !bytes.Equal(rawData, rawResult) {
		t.Fatalf("Expected\n%v\ngot\n%v", rawData, rawResult)
	}
}

func TestEventString(t *testing.T) {
	scenarios := []testResponseWriteScenario[string]{
		{
			name:            "no explicit content-type",
			status:          234,
			headers:         nil,
			body:            "test",
			expectedStatus:  234,
			expectedHeaders: map[string]string{"content-type": "text/plain; charset=utf-8"},
			expectedBody:    "test",
		},
		{
			name:            "with explicit content-type",
			status:          234,
			headers:         map[string]string{"content-type": "text/test"},
			body:            "test",
			expectedStatus:  234,
			expectedHeaders: map[string]string{"content-type": "text/test"},
			expectedBody:    "test",
		},
	}

	for _, s := range scenarios {
		testEventResponseWrite(t, s, func(e *router.Event) error {
			return e.String(s.status, s.body)
		})
	}
}

func TestEventHTML(t *testing.T) {
	scenarios := []testResponseWriteScenario[string]{
		{
			name:            "no explicit content-type",
			status:          234,
			headers:         nil,
			body:            "test",
			expectedStatus:  234,
			expectedHeaders: map[string]string{"content-type": "text/html; charset=utf-8"},
			expectedBody:    "test",
		},
		{
			name:            "with explicit content-type",
			status:          234,
			headers:         map[string]string{"content-type": "text/test"},
			body:            "test",
			expectedStatus:  234,
			expectedHeaders: map[string]string{"content-type": "text/test"},
			expectedBody:    "test",
		},
	}

	for _, s := range scenarios {
		testEventResponseWrite(t, s, func(e *router.Event) error {
			return e.HTML(s.status, s.body)
		})
	}
}

func TestEventJSON(t *testing.T) {
	body := map[string]any{"a": 123, "b": 456, "c": "test"}
	expectedPickedBody := `{"a":123,"c":"test"}` + "\n"
	expectedFullBody := `{"a":123,"b":456,"c":"test"}` + "\n"

	scenarios := []testResponseWriteScenario[any]{
		{
			name:            "no explicit content-type",
			status:          200,
			headers:         nil,
			body:            body,
			expectedStatus:  200,
			expectedHeaders: map[string]string{"content-type": "application/json"},
			expectedBody:    expectedPickedBody,
		},
		{
			name:            "with explicit content-type (200)",
			status:          200,
			headers:         map[string]string{"content-type": "application/test"},
			body:            body,
			expectedStatus:  200,
			expectedHeaders: map[string]string{"content-type": "application/test"},
			expectedBody:    expectedPickedBody,
		},
		{
			name:            "with explicit content-type (400)", // no fields picker
			status:          400,
			headers:         map[string]string{"content-type": "application/test"},
			body:            body,
			expectedStatus:  400,
			expectedHeaders: map[string]string{"content-type": "application/test"},
			expectedBody:    expectedFullBody,
		},
	}

	for _, s := range scenarios {
		testEventResponseWrite(t, s, func(e *router.Event) error {
			e.Request.URL.RawQuery = "fields=a,c" // ensures that the picker is invoked
			return e.JSON(s.status, s.body)
		})
	}
}

func TestEventXML(t *testing.T) {
	scenarios := []testResponseWriteScenario[string]{
		{
			name:            "no explicit content-type",
			status:          234,
			headers:         nil,
			body:            "test",
			expectedStatus:  234,
			expectedHeaders: map[string]string{"content-type": "application/xml; charset=utf-8"},
			expectedBody:    xml.Header + "<string>test</string>",
		},
		{
			name:            "with explicit content-type",
			status:          234,
			headers:         map[string]string{"content-type": "text/test"},
			body:            "test",
			expectedStatus:  234,
			expectedHeaders: map[string]string{"content-type": "text/test"},
			expectedBody:    xml.Header + "<string>test</string>",
		},
	}

	for _, s := range scenarios {
		testEventResponseWrite(t, s, func(e *router.Event) error {
			return e.XML(s.status, s.body)
		})
	}
}

func TestEventStream(t *testing.T) {
	scenarios := []testResponseWriteScenario[string]{
		{
			name:            "stream",
			status:          234,
			headers:         map[string]string{"content-type": "text/test"},
			body:            "test",
			expectedStatus:  234,
			expectedHeaders: map[string]string{"content-type": "text/test"},
			expectedBody:    "test",
		},
	}

	for _, s := range scenarios {
		testEventResponseWrite(t, s, func(e *router.Event) error {
			return e.Stream(s.status, s.headers["content-type"], strings.NewReader(s.body))
		})
	}
}

func TestEventBlob(t *testing.T) {
	scenarios := []testResponseWriteScenario[[]byte]{
		{
			name:            "blob",
			status:          234,
			headers:         map[string]string{"content-type": "text/test"},
			body:            []byte("test"),
			expectedStatus:  234,
			expectedHeaders: map[string]string{"content-type": "text/test"},
			expectedBody:    "test",
		},
	}

	for _, s := range scenarios {
		testEventResponseWrite(t, s, func(e *router.Event) error {
			return e.Blob(s.status, s.headers["content-type"], s.body)
		})
	}
}

func TestEventNoContent(t *testing.T) {
	s := testResponseWriteScenario[any]{
		name:            "no content",
		status:          234,
		headers:         map[string]string{"content-type": "text/test"},
		body:            nil,
		expectedStatus:  234,
		expectedHeaders: map[string]string{"content-type": "text/test"},
		expectedBody:    "",
	}

	testEventResponseWrite(t, s, func(e *router.Event) error {
		return e.NoContent(s.status)
	})
}

func TestEventFlush(t *testing.T) {
	rec := httptest.NewRecorder()

	event := &router.Event{
		Response: unwrapTester{&router.ResponseWriter{ResponseWriter: rec}},
	}
	event.Response.Write([]byte("test"))
	event.Flush()

	if !rec.Flushed {
		t.Fatal("Expected response to be flushed")
	}
}

func TestEventRedirect(t *testing.T) {
	scenarios := []testResponseWriteScenario[any]{
		{
			name:           "non-30x status",
			status:         200,
			expectedStatus: 200,
			expectedError:  router.ErrInvalidRedirectStatusCode,
		},
		{
			name:            "30x status",
			status:          302,
			headers:         map[string]string{"location": "test"}, // should be overwritten with the argument
			expectedStatus:  302,
			expectedHeaders: map[string]string{"location": "example"},
		},
	}

	for _, s := range scenarios {
		testEventResponseWrite(t, s, func(e *router.Event) error {
			return e.Redirect(s.status, "example")
		})
	}
}

func TestEventFileFS(t *testing.T) {
	// stub test files
	// ---
	dir, err := os.MkdirTemp("", "EventFileFS")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	err = os.WriteFile(filepath.Join(dir, "index.html"), []byte("index"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(dir, "test.txt"), []byte("test"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// create sub directory with an index.html file inside it
	err = os.MkdirAll(filepath.Join(dir, "sub1"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(dir, "sub1", "index.html"), []byte("sub1 index"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(filepath.Join(dir, "sub2"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(filepath.Join(dir, "sub2", "test.txt"), []byte("sub2 test"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	// ---

	scenarios := []struct {
		name     string
		path     string
		expected string
	}{
		{"missing file", "", ""},
		{"root with no explicit file", "", ""},
		{"root with explicit file", "test.txt", "test"},
		{"sub dir with no explicit file", "sub1", "sub1 index"},
		{"sub dir with no explicit file (no index.html)", "sub2", ""},
		{"sub dir explicit file", "sub2/test.txt", "sub2 test"},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			rec := httptest.NewRecorder()

			event := &router.Event{
				Request:  req,
				Response: rec,
			}

			err = event.FileFS(os.DirFS(dir), s.path)

			hasErr := err != nil
			expectErr := s.expected == ""
			if hasErr != expectErr {
				t.Fatalf("Expected hasErr %v, got %v (%v)", expectErr, hasErr, err)
			}

			result := rec.Result()

			raw, err := io.ReadAll(result.Body)
			result.Body.Close()
			if err != nil {
				t.Fatal(err)
			}

			if string(raw) != s.expected {
				t.Fatalf("Expected body\n%s\ngot\n%s", s.expected, raw)
			}

			// ensure that the proper file headers are added
			// (aka. http.ServeContent is invoked)
			length, _ := strconv.Atoi(result.Header.Get("content-length"))
			if length != len(s.expected) {
				t.Fatalf("Expected Content-Length %d, got %d", len(s.expected), length)
			}
		})
	}
}

func TestEventError(t *testing.T) {
	err := new(router.Event).Error(123, "message_test", map[string]any{"a": validation.Required, "b": "test"})

	result, _ := json.Marshal(err)
	expected := `{"data":{"a":{"code":"validation_invalid_value","message":"Invalid value."},"b":{"code":"validation_invalid_value","message":"Invalid value."}},"message":"Message_test.","status":123}`

	if string(result) != expected {
		t.Errorf("Expected\n%s\ngot\n%s", expected, result)
	}
}

func TestEventBadRequestError(t *testing.T) {
	err := new(router.Event).BadRequestError("message_test", map[string]any{"a": validation.Required, "b": "test"})

	result, _ := json.Marshal(err)
	expected := `{"data":{"a":{"code":"validation_invalid_value","message":"Invalid value."},"b":{"code":"validation_invalid_value","message":"Invalid value."}},"message":"Message_test.","status":400}`

	if string(result) != expected {
		t.Errorf("Expected\n%s\ngot\n%s", expected, result)
	}
}

func TestEventNotFoundError(t *testing.T) {
	err := new(router.Event).NotFoundError("message_test", map[string]any{"a": validation.Required, "b": "test"})

	result, _ := json.Marshal(err)
	expected := `{"data":{"a":{"code":"validation_invalid_value","message":"Invalid value."},"b":{"code":"validation_invalid_value","message":"Invalid value."}},"message":"Message_test.","status":404}`

	if string(result) != expected {
		t.Errorf("Expected\n%s\ngot\n%s", expected, result)
	}
}

func TestEventForbiddenError(t *testing.T) {
	err := new(router.Event).ForbiddenError("message_test", map[string]any{"a": validation.Required, "b": "test"})

	result, _ := json.Marshal(err)
	expected := `{"data":{"a":{"code":"validation_invalid_value","message":"Invalid value."},"b":{"code":"validation_invalid_value","message":"Invalid value."}},"message":"Message_test.","status":403}`

	if string(result) != expected {
		t.Errorf("Expected\n%s\ngot\n%s", expected, result)
	}
}

func TestEventUnauthorizedError(t *testing.T) {
	err := new(router.Event).UnauthorizedError("message_test", map[string]any{"a": validation.Required, "b": "test"})

	result, _ := json.Marshal(err)
	expected := `{"data":{"a":{"code":"validation_invalid_value","message":"Invalid value."},"b":{"code":"validation_invalid_value","message":"Invalid value."}},"message":"Message_test.","status":401}`

	if string(result) != expected {
		t.Errorf("Expected\n%s\ngot\n%s", expected, result)
	}
}

func TestEventTooManyRequestsError(t *testing.T) {
	err := new(router.Event).TooManyRequestsError("message_test", map[string]any{"a": validation.Required, "b": "test"})

	result, _ := json.Marshal(err)
	expected := `{"data":{"a":{"code":"validation_invalid_value","message":"Invalid value."},"b":{"code":"validation_invalid_value","message":"Invalid value."}},"message":"Message_test.","status":429}`

	if string(result) != expected {
		t.Errorf("Expected\n%s\ngot\n%s", expected, result)
	}
}

func TestEventInternalServerError(t *testing.T) {
	err := new(router.Event).InternalServerError("message_test", map[string]any{"a": validation.Required, "b": "test"})

	result, _ := json.Marshal(err)
	expected := `{"data":{"a":{"code":"validation_invalid_value","message":"Invalid value."},"b":{"code":"validation_invalid_value","message":"Invalid value."}},"message":"Message_test.","status":500}`

	if string(result) != expected {
		t.Errorf("Expected\n%s\ngot\n%s", expected, result)
	}
}

func TestEventBindBody(t *testing.T) {
	type testDstStruct struct {
		A int    `json:"a" xml:"a" form:"a"`
		B int    `json:"b" xml:"b" form:"b"`
		C string `json:"c" xml:"c" form:"c"`
	}

	emptyDst := `{"a":0,"b":0,"c":""}`

	queryDst := `a=123&b=-456&c=test`

	xmlDst := `
		<?xml version="1.0" encoding="UTF-8" ?>
		<root>
			<a>123</a>
			<b>-456</b>
			<c>test</c>
		</root>
	`

	jsonDst := `{"a":123,"b":-456,"c":"test"}`

	// multipart
	mpBody := &bytes.Buffer{}
	mpWriter := multipart.NewWriter(mpBody)
	mpWriter.WriteField("@jsonPayload", `{"a":123}`)
	mpWriter.WriteField("b", "-456")
	mpWriter.WriteField("c", "test")
	if err := mpWriter.Close(); err != nil {
		t.Fatal(err)
	}

	scenarios := []struct {
		contentType string
		body        io.Reader
		expectDst   string
		expectError bool
	}{
		{
			contentType: "",
			body:        strings.NewReader(jsonDst),
			expectDst:   emptyDst,
			expectError: true,
		},
		{
			contentType: "application/rtf", // unsupported
			body:        strings.NewReader(jsonDst),
			expectDst:   emptyDst,
			expectError: true,
		},
		// empty body
		{
			contentType: "application/json;charset=emptybody",
			body:        strings.NewReader(""),
			expectDst:   emptyDst,
		},
		// json
		{
			contentType: "application/json",
			body:        strings.NewReader(jsonDst),
			expectDst:   jsonDst,
		},
		{
			contentType: "application/json;charset=abc",
			body:        strings.NewReader(jsonDst),
			expectDst:   jsonDst,
		},
		// xml
		{
			contentType: "text/xml",
			body:        strings.NewReader(xmlDst),
			expectDst:   jsonDst,
		},
		{
			contentType: "text/xml;charset=abc",
			body:        strings.NewReader(xmlDst),
			expectDst:   jsonDst,
		},
		{
			contentType: "application/xml",
			body:        strings.NewReader(xmlDst),
			expectDst:   jsonDst,
		},
		{
			contentType: "application/xml;charset=abc",
			body:        strings.NewReader(xmlDst),
			expectDst:   jsonDst,
		},
		// x-www-form-urlencoded
		{
			contentType: "application/x-www-form-urlencoded",
			body:        strings.NewReader(queryDst),
			expectDst:   jsonDst,
		},
		{
			contentType: "application/x-www-form-urlencoded;charset=abc",
			body:        strings.NewReader(queryDst),
			expectDst:   jsonDst,
		},
		// multipart
		{
			contentType: mpWriter.FormDataContentType(),
			body:        mpBody,
			expectDst:   jsonDst,
		},
	}

	for _, s := range scenarios {
		t.Run(s.contentType, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost, "/", s.body)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("content-type", s.contentType)

			event := &router.Event{Request: req}

			dst := testDstStruct{}

			err = event.BindBody(&dst)

			hasErr := err != nil
			if hasErr != s.expectError {
				t.Fatalf("Expected hasErr %v, got %v (%v)", s.expectError, hasErr, err)
			}

			dstRaw, err := json.Marshal(dst)
			if err != nil {
				t.Fatal(err)
			}

			if string(dstRaw) != s.expectDst {
				t.Fatalf("Expected dst\n%s\ngot\n%s", s.expectDst, dstRaw)
			}
		})
	}
}

// -------------------------------------------------------------------

type testResponseWriteScenario[T any] struct {
	name            string
	status          int
	headers         map[string]string
	body            T
	expectedStatus  int
	expectedHeaders map[string]string
	expectedBody    string
	expectedError   error
}

func testEventResponseWrite[T any](
	t *testing.T,
	scenario testResponseWriteScenario[T],
	writeFunc func(e *router.Event) error,
) {
	t.Run(scenario.name, func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rec := httptest.NewRecorder()
		event := &router.Event{
			Request:  req,
			Response: &router.ResponseWriter{ResponseWriter: rec},
		}

		for k, v := range scenario.headers {
			event.Response.Header().Add(k, v)
		}

		err = writeFunc(event)
		if (scenario.expectedError != nil || err != nil) && !errors.Is(err, scenario.expectedError) {
			t.Fatalf("Expected error %v, got %v", scenario.expectedError, err)
		}

		result := rec.Result()

		if result.StatusCode != scenario.expectedStatus {
			t.Fatalf("Expected status code %d, got %d", scenario.expectedStatus, result.StatusCode)
		}

		resultBody, err := io.ReadAll(result.Body)
		result.Body.Close()
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		resultBody, err = json.Marshal(string(resultBody))
		if err != nil {
			t.Fatal(err)
		}

		expectedBody, err := json.Marshal(scenario.expectedBody)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(resultBody, expectedBody) {
			t.Fatalf("Expected body\n%s\ngot\n%s", expectedBody, resultBody)
		}

		for k, ev := range scenario.expectedHeaders {
			if v := result.Header.Get(k); v != ev {
				t.Fatalf("Expected %q header to be %q, got %q", k, ev, v)
			}
		}
	})
}
