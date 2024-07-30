package gbfs

// https://github.com/NABSA/gbfs/blob/master/gbfs.md#station_informationjson

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type SystemInformationResponse struct {
	Header
	Data SystemInformation `json:"data"`
}

type SystemInformation struct {
	SystemID                    ID         `json:"system_id"`
	Language                    string     `json:"language"`
	Name                        string     `json:"name"`
	ShortName                   string     `json:"short_name,omitempty"`
	Operator                    string     `json:"operator,omitempty"`
	URL                         string     `json:"url,omitempty"`
	PurchaseURL                 string     `json:"purchase_url,omitempty"`
	StartDate                   Date       `json:"start_date,omitempty"`
	PhoneNumber                 string     `json:"phone_number,omitempty"`
	Email                       string     `json:"email,omitempty"`
	FeedContactEmail            string     `json:"feed_contact_email,omitempty"` // added in v1.1
	Timezone                    string     `json:"timezone"`
	LicenseID                   string     `json:"license_id,omitempty"`
	LicenseURL                  string     `json:"license_url,omitempty"`
	AttributionOrganizationName string     `json:"attribution_organization_name,omitempty"`
	AttributionURL              string     `json:"attribution_url,omitempty"`
	RentalApps                  RentalApps `json:"rental_apps,omitempty"` // added in v1.1
}

type RentalApps struct {
	Android RentalAppsStore `json:"android,omitempty"`
	IOS     RentalAppsStore `json:"ios,omitempty"`
}

type RentalAppsStore struct {
	StoreURI     string `json:"store_uri"`
	DiscoveryURI string `json:"discovery_uri"`
}

func (c *Client) LoadSystemInformation(url string) (*SystemInformationResponse, error) {
	resp, err := c.sendRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	var r SystemInformationResponse

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}

	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errors.Wrap(err, "unmarshal JSON")
	}

	return &r, nil
}
