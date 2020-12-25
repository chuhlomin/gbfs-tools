package structs

type System struct {
	ID               string `json:"id" db:"system_id" storm:"id"`
	CountryCode      string `json:"countryCode" db:"country_code" storm:"index"`
	Name             string `json:"name" db:"name"`
	Location         string `json:"location" db:"location" storm:"index"`
	URL              string `json:"url" db:"url"`
	AutoDiscoveryURL string `json:"autoDiscoveryUrl" db:"auto_discovery_url"`
	IsEnabled        bool   `json:"isEnabled" db:"is_enabled" storm:"index"`
}
