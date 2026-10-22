package update

import (
	"fmt"
	"net/url"
	"strings"
)

type release struct {
	Id             string   `json:"id"`
	CollectionName string   `json:"collectionName"`
	Channel        string   `json:"channel"`
	Tag            string   `json:"tag"`
	Content        string   `json:"content"`
	Created        string   `json:"created"`
	Updated        string   `json:"updated"`
	Files          []string `json:"files"`
}

// findFileBySuffix returns the first available asset containing the specified suffix.
func (r *release) findFileBySuffix(suffix string) (string, error) {
	if suffix != "" {
		for _, filename := range r.Files {
			if strings.HasSuffix(filename, suffix) {
				return filename, nil
			}
		}
	}

	return "", fmt.Errorf("failed to find release file ending with %q", suffix)
}

// buildFilePath constructs the URL path to the release file.
func (r *release) buildFilePath(filename string) string {
	return fmt.Sprintf(
		"/api/files/%s/%s/%s",
		url.PathEscape(r.CollectionName),
		url.PathEscape(r.Id),
		url.PathEscape(filename),
	)
}
