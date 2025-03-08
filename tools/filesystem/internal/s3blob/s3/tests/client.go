// Package tests contains various tests helpers and utilities to assist
// with the S3 client testing.
package tests

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"sync"
)

// NewClient creates a new test Client loaded with the specified RequestStubs.
func NewClient(stubs ...*RequestStub) *Client {
	return &Client{stubs: stubs}
}

type RequestStub struct {
	Method   string
	URL      string // plain string or regex pattern wrapped in "^pattern$"
	Match    func(req *http.Request) bool
	Response *http.Response
}

type Client struct {
	stubs []*RequestStub
	mu    sync.Mutex
}

// AssertNoRemaining asserts that current client has no unprocessed requests remaining.
func (c *Client) AssertNoRemaining() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.stubs) == 0 {
		return nil
	}

	msgParts := make([]string, 0, len(c.stubs)+1)
	msgParts = append(msgParts, "not all stub requests were processed:")
	for _, stub := range c.stubs {
		msgParts = append(msgParts, "- "+stub.Method+" "+stub.URL)
	}

	return errors.New(strings.Join(msgParts, "\n"))
}

// Do implements the [s3.HTTPClient] interface.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, stub := range c.stubs {
		if req.Method != stub.Method {
			continue
		}

		urlPattern := stub.URL
		if !strings.HasPrefix(urlPattern, "^") && !strings.HasSuffix(urlPattern, "$") {
			urlPattern = "^" + regexp.QuoteMeta(urlPattern) + "$"
		}

		urlRegex, err := regexp.Compile(urlPattern)
		if err != nil {
			return nil, err
		}

		if !urlRegex.MatchString(req.URL.String()) {
			continue
		}

		if stub.Match != nil && !stub.Match(req) {
			continue
		}

		// remove from the remaining stubs
		c.stubs = slices.Delete(c.stubs, i, i+1)

		response := stub.Response
		if response == nil {
			response = &http.Response{}
		}
		if response.Header == nil {
			response.Header = http.Header{}
		}
		if response.Body == nil {
			response.Body = http.NoBody
		}

		response.Request = req

		return response, nil
	}

	var body []byte
	if req.Body != nil {
		defer req.Body.Close()
		body, _ = io.ReadAll(req.Body)
	}

	return nil, fmt.Errorf(
		"the below request doesn't have a corresponding stub:\n%s %s\nHeaders: %v\nBody: %q",
		req.Method,
		req.URL.String(),
		req.Header,
		body,
	)
}
