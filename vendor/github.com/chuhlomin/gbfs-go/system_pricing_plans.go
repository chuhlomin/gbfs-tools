package gbfs

// https://github.com/NABSA/gbfs/blob/master/gbfs.md#station_informationjson

import (
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
)

type SystemPricingPlansResponse struct {
	Header
	Data SystemPricingPlans `json:"data"`
}

type SystemPricingPlans struct {
	Plans []Plan `json:"plans"`
}

type Plan struct {
	ID            ID        `json:"plan_id"`
	URL           string    `json:"url,omitempty"`
	Name          string    `json:"name"`
	Currency      string    `json:"currency"`
	Price         Price     `json:"price"`
	IsTaxable     Bool      `json:"is_taxable"`
	Description   string    `json:"description"`
	PerKmPricing  []Pricing `json:"per_km_pricing,omitempty"`  // added in v2.1-RC2
	PerMinPricing []Pricing `json:"per_min_pricing,omitempty"` // added in v2.1-RC2
}

type Price interface{}

type Pricing struct {
	Start    uint16  `json:"start"`
	Rate     float64 `json:"rate"`
	Interval uint16  `json:"intercal"`
	End      uint16  `json:"end,omitempty"`
}

func (c *Client) LoadSystemPricingPlans(url string) (*SystemPricingPlansResponse, error) {
	resp, err := c.sendRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	var r SystemPricingPlansResponse

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read response body")
	}

	if err := json.Unmarshal(b, &r); err != nil {
		return nil, errors.Wrap(err, "unmarshal JSON")
	}

	return &r, nil
}
