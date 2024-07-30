package gbfs

// https://github.com/NABSA/gbfs/blob/master/gbfs.md#gbfsjson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/pkg/errors"
)

// GBFSResponse represent a single system or geographic area
// in which vehicles are operated.
type GBFSResponse struct {
	Header
	Data LanguageFeeds `json:"data"`
}

// Every JSON file presented in the specification
// contains the same common header information
// at the top level of the JSON response object.
type Header struct {
	LastUpdated Timestamp `json:"last_updated"`
	TTL         int       `json:"ttl"`
	Version     string    `json:"version"` // added in v1.1
}

type LanguageFeeds map[string]DataFeeds

type DataFeeds struct {
	Feeds []Feed `json:"feeds"`
}

func (df *DataFeeds) GetFeed(name string) (*Feed, error) {
	for _, f := range df.Feeds {
		if f.Name == name {
			return &f, nil
		}
	}
	return nil, fmt.Errorf("no language feeds found")
}

type Feed struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (lf LanguageFeeds) GetDataFeeds(language string) (*DataFeeds, error) {
	if f, found := lf[language]; found {
		return &f, nil
	}
	for _, f := range lf {
		return &f, nil
	}
	return nil, fmt.Errorf("no language feeds found")
}

func (lf *LanguageFeeds) UnmarshalJSON(data []byte) error {
	var v map[string]json.RawMessage
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if message, ok := v["feeds"]; ok {
		// Some systems don't follow specification directly
		// and exposing "feeds" right inside "data", omitting "LanguageFeeds" map
		var feeds []Feed
		if err := json.Unmarshal(message, &feeds); err != nil {
			return err
		}

		*(*LanguageFeeds)(lf) = map[string]DataFeeds{
			"default": {
				Feeds: feeds,
			},
		}
		return nil
	}

	result := LanguageFeeds{}

	for language, message := range v {
		var df DataFeeds
		if err := json.Unmarshal(message, &df); err != nil {
			return err
		}
		result[language] = df
	}

	*lf = result
	return nil
}

func (c *Client) LoadGBFS(url string) (*GBFSResponse, error) {
	resp, err := c.sendRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	var r GBFSResponse

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}

	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errors.Wrap(err, "unmarshal JSON")
	}

	return &r, nil
}
