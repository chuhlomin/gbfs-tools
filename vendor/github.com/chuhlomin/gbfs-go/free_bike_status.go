package gbfs

// https://github.com/NABSA/gbfs/blob/master/gbfs.md#free_bike_statusjson

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type FreeBikeStatusResponse struct {
	Header
	Data FreeBikeStatusData `json:"data"`
}

type FreeBikeStatusData struct {
	Bikes []FreeBikeStatus `json:"bikes"`
}

type FreeBikeStatus struct {
	BikeID             ID         `json:"bike_id"`
	SystemID           ID         `json:"system_id"` // added in v3.0-RC
	Lat                float64    `json:"lat"`
	Lon                float64    `json:"lon"`
	IsReserved         Bool       `json:"is_reserved"`
	IsDisabled         Bool       `json:"is_disabled"`
	RentalURIs         RentalURIs `json:"rental_uris"`          // added in v1.1
	VehicleTypeID      ID         `json:"vehicle_type_id"`      // added in v2.1-RC
	LastReported       Timestamp  `json:"last_reported"`        // added in v2.1-RC
	CurrentRangeMeters float64    `json:"current_range_meters"` // added in v2.1-RC
	StationID          ID         `json:"station_id"`           // added in v2.1-RC
	PricingPlanID      ID         `json:"pricing_plan_id"`      // added in v2.1-RC
}

func (c *Client) LoadFreeBikeStatus(url string) (*FreeBikeStatusResponse, error) {
	resp, err := c.sendRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	var r FreeBikeStatusResponse

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}

	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errors.Wrap(err, "unmarshal JSON")
	}

	return &r, nil
}
