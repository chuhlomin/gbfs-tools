package structs

type System struct {
	ID               string `json:"id"`
	CountryCode      string `json:"countryCode"`
	Name             string `json:"name"`
	Location         string `json:"location"`
	URL              string `json:"url"`
	AutoDiscoveryURL string `json:"autoDiscoveryUrl"`
	IsEnabled        bool   `json:"isEnabled"`
}
