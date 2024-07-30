package gbfs

// https://github.com/NABSA/gbfs/blob/master/gbfs.md#system_hoursjson

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type SystemHoursResponse struct {
	Header
	Data SystemHoursData `json:"data"`
}

type SystemHoursData struct {
	RentalHours []RentalHours `json:"rental_hours"`
}

type RentalHours struct {
	UserTypes []UserType `json:"user_types"`
	Days      []Weekday  `json:"days"`
	StartTime Clock      `json:"start_time"`
	EndTime   Clock      `json:"end_time"`
}

type UserType string

const UserTypeMember UserType = "member"
const UserTypeNonMember UserType = "nonmember"

func (c *Client) LoadSystemHours(url string) (*SystemHoursResponse, error) {
	resp, err := c.sendRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	var r SystemHoursResponse

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}

	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errors.Wrap(err, "unmarshal JSON")
	}

	return &r, nil
}
