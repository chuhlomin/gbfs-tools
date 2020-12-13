package structs

type System struct {
	ID               string `json:"id" db:"system_id"`
	CountryCode      string `json:"countryCode" db:"country_code"`
	Name             string `json:"name" db:"name"`
	Location         string `json:"location" db:"location"`
	URL              string `json:"url" db:"url"`
	AutoDiscoveryURL string `json:"autoDiscoveryUrl" db:"auto_discovery_url"`
	IsEnabled        bool   `db:"is_enabled"`
}
