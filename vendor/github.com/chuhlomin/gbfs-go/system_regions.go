package gbfs

// https://github.com/NABSA/gbfs/blob/master/gbfs.md#station_informationjson

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type SystemRegionsResponse struct {
	Header
	Data SystemRegions `json:"data"`
}

type SystemRegions struct {
	Regions []Region `json:"regions"`
}

type Region struct {
	ID   ID     `json:"region_id"`
	Name string `json:"name"`
}

func (c *Client) LoadSystemRegions(url string) (*SystemRegionsResponse, error) {
	resp, err := c.sendRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	var r SystemRegionsResponse

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}

	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errors.Wrap(err, "unmarshal JSON")
	}

	return &r, nil
}
