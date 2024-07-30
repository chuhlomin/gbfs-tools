package gbfs

// https://github.com/NABSA/gbfs/blob/master/gbfs.md#system_calendarjson

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type SystemCalendarResponse struct {
	Header
	Data SystemCalendarData `json:"data"`
}

type SystemCalendarData struct {
	Calendars []Calendar `json:"calendars"`
}

type Calendar struct {
	StartMonth uint16 `json:"start_month"`
	StartDay   uint16 `json:"start_day"`
	StartYear  uint16 `json:"start_year,omitempty"`
	EndMonth   uint16 `json:"end_month"`
	EndDay     uint16 `json:"end_day"`
	EndYear    uint16 `json:"end_year,omitempty"`
}

func (c *Client) LoadSystemCalendar(url string) (*SystemCalendarResponse, error) {
	resp, err := c.sendRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	var r SystemCalendarResponse

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}

	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errors.Wrap(err, "unmarshal JSON")
	}

	return &r, nil
}
