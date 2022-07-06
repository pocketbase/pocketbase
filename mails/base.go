// Package mails implements various helper methods for sending user and admin
// emails like forgotten password, verification, etc.
package mails

import (
	"bytes"
	"net/url"
	"path"
	"strings"
	"text/template"
)

// normalizeUrl removes duplicated slashes from a url path.
func normalizeUrl(originalUrl string) (string, error) {
	u, err := url.Parse(originalUrl)
	if err != nil {
		return "", err
	}

	hasSlash := strings.HasSuffix(u.Path, "/")

	// clean up path by removing duplicated /
	u.Path = path.Clean(u.Path)
	u.RawPath = path.Clean(u.RawPath)

	// restore original trailing slash
	if hasSlash && !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
		u.RawPath += "/"
	}

	return u.String(), nil
}

// resolveTemplateContent resolves inline html template strings.
func resolveTemplateContent(data any, content ...string) (string, error) {
	if len(content) == 0 {
		return "", nil
	}

	t := template.New("inline_template")

	var parseErr error
	for _, v := range content {
		t, parseErr = t.Parse(v)
		if parseErr != nil {
			return "", parseErr
		}
	}

	var wr bytes.Buffer

	if executeErr := t.Execute(&wr, data); executeErr != nil {
		return "", executeErr
	}

	return wr.String(), nil
}
