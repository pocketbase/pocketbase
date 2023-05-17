package ghupdate

import (
	"errors"
	"strings"
)

type releaseAsset struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Size        int    `json:"size"`
	DownloadUrl string `json:"browser_download_url"`
}

type release struct {
	Id        int             `json:"id"`
	Name      string          `json:"name"`
	Tag       string          `json:"tag_name"`
	Published string          `json:"published_at"`
	Url       string          `json:"html_url"`
	Assets    []*releaseAsset `json:"assets"`
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
