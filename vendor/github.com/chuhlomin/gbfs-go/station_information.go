package gbfs

// https://github.com/NABSA/gbfs/blob/master/gbfs.md#system_informationjson

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type StationInformationResponse struct {
	Header
	Data StationInformationData `json:"data"`
}

type StationInformationData struct {
	Stations []StationInformation `json:"stations"`
}

type StationInformation struct {
	ID                  ID             `json:"station_id"`
	Name                string         `json:"name"`
	ShortName           string         `json:"short_name,omitempty"`
	Lat                 float64        `json:"lat"`
	Lon                 float64        `json:"lon"`
	Address             string         `json:"address,omitempty"`
	CrossStreet         string         `json:"cross_street,omitempty"`
	RegionID            ID             `json:"region_id,omitempty"`
	PostCode            string         `json:"post_code,omitempty"`
	RentalMethods       []RentalMethod `json:"rental_methods,omitempty"`
	IsVirtualStation    *bool          `json:"is_virtual_station,omitempty"` // added in v2.1-RC
	StationArea         interface{}    `json:"station_area,omitempty"`       // added in v2.1-RC
	Capacity            int            `json:"capacity,omitempty"`
	VehicleCapacity     Capacity       `json:"vehicle_capacity"`                // added in v2.1-RC
	IsValetStation      *bool          `json:"is_valet_station,omitempty"`      // added in v2.1-RC
	RentalURIs          RentalURIs     `json:"rental_uris,omitempty"`           // added in v1.1
	VehicleTypeCapacity Capacity       `json:"vehicle_type_capacity,omitempty"` // added in v2.1
}

type RentalMethod string

const RentalMethodKey RentalMethod = "KEY"
const RentalMethodCreditCard RentalMethod = "CREDITCARD"
const RentalMethodPayPass RentalMethod = "PAYPASS"
const RentalMethodApplePay RentalMethod = "APPLEPAY"
const RentalMethodAndroidPay RentalMethod = "ANDROIDPAY"
const RentalMethodTransitCard RentalMethod = "TRANSITCARD"
const RentalMethodAccountNumber RentalMethod = "ACCOUNTNUMBER"
const RentalMethodPhone RentalMethod = "PHONE"

type RentalURIs struct {
	Android string `json:"android,omitempty"`
	IOS     string `json:"ios,omitempty"`
	Web     string `json:"web,omitempty"`
}

type Capacity map[string]int

func (c *Client) LoadStationInformation(url string) (*StationInformationResponse, error) {
	resp, err := c.sendRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	var r StationInformationResponse

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}

	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errors.Wrap(err, "unmarshal JSON")
	}

	return &r, nil
}
