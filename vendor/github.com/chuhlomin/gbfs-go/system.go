package gbfs

import (
	"encoding/csv"
	"io"

	"github.com/pkg/errors"
)

// North American Bike Share Association
// URL to systems.csv
// See https://github.com/NABSA/gbfs#systems-implementing-gbfs
const SystemsNABSA = "https://github.com/NABSA/gbfs/raw/master/systems.csv"

// System represents system publishing GBFS feeds
type System struct {
	ID               string `json:"id" db:"id"`
	CountryCode      string `json:"country_code" db:"country_code"`
	Name             string `json:"name" db:"name"`
	Location         string `json:"location" db:"location"`
	URL              string `json:"url" db:"url"`
	AutoDiscoveryURL string `json:"auto_discovery_url" db:"auto_discovery_url"`
}

// LoadSystem gets URL to systems.csv file and returns parsed systems
func (c *Client) LoadSystems(url string) ([]System, error) {
	resp, err := c.sendRequest(url)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	return ParseSystemsCSV(resp.Body)
}

func ParseSystemsCSV(reader io.Reader) ([]System, error) {
	var systems []System

	r := csv.NewReader(reader)
	var header []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return systems, errors.Wrap(err, "read CSV body")
		}

		if header == nil {
			header = record
		} else {
			system := System{}
			for i, column := range header {
				switch column {
				case "Country Code":
					system.CountryCode = record[i]
				case "Name":
					system.Name = record[i]
				case "Location":
					system.Location = record[i]
				case "System ID":
					system.ID = record[i]
				case "URL":
					system.URL = record[i]
				case "Auto-Discovery URL":
					system.AutoDiscoveryURL = record[i]
				}
			}
			systems = append(systems, system)
		}
	}

	return systems, nil
}
