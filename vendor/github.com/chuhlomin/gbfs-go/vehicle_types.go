package gbfs

// https://github.com/NABSA/gbfs/blob/master/gbfs.md#vehicle_typesjson-added-in-v21-rc
// added in v2.1-RC

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type VehicleTypesResponse struct {
	Header
	Data VehicleTypesData `json:"data"`
}

type VehicleTypesData struct {
	VehicleTypes []VehicleType `json:"vehicle_types"`
}

type VehicleType struct {
	VehicleTypeID  ID             `json:"vehicle_type_id"`
	FormFactor     FormFactor     `json:"form_factor"`
	PropulsionType PropulsionType `json:"propulsion_type"`
	MaxRangeMeters float64        `json:"max_range_meters,omitempty"`
	Name           string         `json:"name,omitempty"`
}

type FormFactor string

const FormFactorBicycle FormFactor = "bicycle"
const FormFactorCar FormFactor = "car"
const FormFactorMoped FormFactor = "moped"
const FormFactorScooter FormFactor = "scooter"
const FormFactorOther FormFactor = "other"

type PropulsionType string

const PropulsionTypeHuman = "human"
const PropulsionTypeElectricAssist = "electric_assist"
const PropulsionTypeElectric = "electric"
const PropulsionTypeCombustion = "combustion"

func (c *Client) LoadVehicleTypes(url string) (*VehicleTypesResponse, error) {
	resp, err := c.sendRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	var r VehicleTypesResponse

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}

	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errors.Wrap(err, "unmarshal JSON")
	}

	return &r, nil
}
