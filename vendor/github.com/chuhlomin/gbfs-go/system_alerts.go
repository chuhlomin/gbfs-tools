package gbfs

// https://github.com/NABSA/gbfs/blob/master/gbfs.md#system_calendarjson

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type SystemAlertsResponse struct {
	Header
	Data SystemAlertsData `json:"data"`
}

type SystemAlertsData struct {
	Alerts []Alert `json:"alerts"`
}

type Alert struct {
	ID          ID          `json:"alert_id"`
	Type        AlertType   `json:"type"`
	Times       []AlertTime `json:"times,omitempty"`
	StationIDs  []string    `json:"station_ids,omitempty"`
	RegionIDs   []string    `json:"region_ids,omitempty"`
	URL         string      `json:"url,omitempty"`
	Summary     string      `json:"summary"`
	Description string      `json:"description,omitempty"`
	LastUpdated Timestamp   `json:"last_updated,omitempty"`
}

type AlertType string

const AlertSystemClosure AlertType = "SYSTEM_CLOSURE"
const AlertStationClosure AlertType = "STATION_CLOSURE"
const AlertStationMove AlertType = "STATION_MOVE"
const AlertOther AlertType = "OTHER"

type AlertTime struct {
	Start Timestamp `json:"start"`
	End   Timestamp `json:"end,omitempty"`
}

func (c *Client) LoadSystemAlerts(url string) (*SystemAlertsResponse, error) {
	resp, err := c.sendRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	var r SystemAlertsResponse

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}

	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errors.Wrap(err, "unmarshal JSON")
	}

	return &r, nil
}
