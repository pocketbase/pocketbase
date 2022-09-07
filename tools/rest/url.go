package rest

import (
	"net/url"
	"path"
	"strings"
)

// NormalizeUrl removes duplicated slashes from a url path.
func NormalizeUrl(originalUrl string) (string, error) {
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
