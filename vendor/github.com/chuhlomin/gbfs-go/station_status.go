package gbfs

// https://github.com/NABSA/gbfs/blob/master/gbfs.md#station_statusjson

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type StationStatusResponse struct {
	Header
	Data StationStatusData `json:"data"`
}

type StationStatusData struct {
	Stations []StationStatus `json:"stations"`
}

type StationStatus struct {
	ID                    ID                    `json:"station_id"`
	NumBikesAvailable     uint                  `json:"num_bikes_available"`
	NumBikesDisabled      uint                  `json:"num_bikes_disabled,omitempty"`
	NumDocksAvailable     uint                  `json:"num_docks_available,omitempty"`
	IsInstalled           Bool                  `json:"is_installed"`
	IsRenting             Bool                  `json:"is_renting"`
	IsReturning           Bool                  `json:"is_returning"`
	LastReported          Timestamp             `json:"last_reported"`
	VehicleTypesAvailable []VehicleAvailability `json:"vehicle_types_available,omitempty"` // added in v2.1-RC
	VehicleDocksAvailable []VehicleAvailability `json:"vehicle_docks_available"`           // added in v2.1-RC
}

type VehicleAvailability struct {
	VehicleTypeID ID   `json:"vehicle_type_id"`
	Count         uint `json:"count"`
}

func (c *Client) LoadStationStatus(url string) (*StationStatusResponse, error) {
	resp, err := c.sendRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	var r StationStatusResponse

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}

	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errors.Wrap(err, "unmarshal JSON")
	}

	return &r, nil
}
