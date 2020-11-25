package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/layer5io/meshery-adapter-library/adapter"
)

// Release is used to save the release informations
type Release struct {
	ID      int             `json:"id,omitempty"`
	TagName string          `json:"tag_name,omitempty"`
	Name    adapter.Version `json:"name,omitempty"`
	Draft   bool            `json:"draft,omitempty"`
	Assets  []*Asset        `json:"assets,omitempty"`
}

// Asset describes the github release asset object
type Asset struct {
	Name        string `json:"name,omitempty"`
	State       string `json:"state,omitempty"`
	DownloadURL string `json:"browser_download_url,omitempty"`
}

// getLatestReleaseNames returns the names of the latest releases
// limited by the "limit" parameter. The first version in the list
// is always is the latest "stable" version.
func getLatestReleaseNames(limit int) ([]adapter.Version, error) {
	releases, err := GetLatestReleases(uint(limit))
	if err != nil {
		return []adapter.Version{}, ErrGetLatestReleaseNames(err)
	}

	var releaseNames []adapter.Version

	for _, r := range releases {
		releaseNames = append(releaseNames, r.Name)
	}

	return releaseNames, nil
}

// GetLatestReleases fetches the latest releases from the kuma repository
func GetLatestReleases(releases uint) ([]*Release, error) {
	releaseAPIURL := "https://api.github.com/repos/kumahq/kuma/releases?per_page=" + fmt.Sprint(releases)
	// We need a variable url here hence using nosec
	// #nosec
	resp, err := http.Get(releaseAPIURL)
	if err != nil {
		return []*Release{}, ErrGetLatestReleases(err)
	}

	if resp.StatusCode != http.StatusOK {
		return []*Release{}, ErrGetLatestReleases(fmt.Errorf("unexpected status code: %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []*Release{}, ErrGetLatestReleases(err)
	}

	var releaseList []*Release

	if err = json.Unmarshal(body, &releaseList); err != nil {
		return []*Release{}, ErrGetLatestReleases(err)
	}

	if err = resp.Body.Close(); err != nil {
		return []*Release{}, ErrGetLatestReleases(err)
	}

	return releaseList, nil
}
