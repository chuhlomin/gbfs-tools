package gbfs

// https://github.com/NABSA/gbfs/blob/master/gbfs.md#gbfs_versionsjson-added-in-v11
// added in v1.1

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type VersionsResponse struct {
	Header
	Data VersionsData `json:"data"`
}

type VersionsData struct {
	Versions []Version `json:"versions"`
}

type Version struct {
	Version string `json:"version"`
	URL     string `json:"url"`
}

func (c *Client) LoadGBFSVersions(url string) (*VersionsResponse, error) {
	resp, err := c.sendRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	var r VersionsResponse

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}

	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errors.Wrap(err, "unmarshal JSON")
	}

	return &r, nil
}
