package ghupdate

import (
	"errors"
	"strings"
)

type releaseAsset struct {
	Name        string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
	Id          int    `json:"id"`
	Size        int    `json:"size"`
}

type release struct {
	Name      string          `json:"name"`
	Tag       string          `json:"tag_name"`
	Published string          `json:"published_at"`
	Url       string          `json:"html_url"`
	Body      string          `json:"body"`
	Assets    []*releaseAsset `json:"assets"`
	Id        int             `json:"id"`
}

// findAssetBySuffix returns the first available asset containing the specified suffix.
func (r *release) findAssetBySuffix(suffix string) (*releaseAsset, error) {
	if suffix != "" {
		for _, asset := range r.Assets {
			if strings.HasSuffix(asset.Name, suffix) {
				return asset, nil
			}
		}
	}

	return nil, errors.New("missing asset containing " + suffix)
}
