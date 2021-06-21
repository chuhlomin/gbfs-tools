package gbfs

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/liangyaopei/structmap"
	gj "github.com/paulmach/go.geojson"

	"github.com/chuhlomin/gbfs-go"
)

type stationProperties struct {
	ID          gbfs.ID `map:"id,omitempty"`
	Name        string  `map:"name,omitempty"`
	Address     string  `map:"address,omitempty"`
	CrossStreet string  `map:"crossStreet,omitempty"`
	Capacity    int     `map:"capacity,omitempty"`
	ShortName   string  `map:"shortName,omitempty"`
	RegionID    gbfs.ID `map:"regionID,omitempty"`
}

func HandlerGeoJSON() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceID := r.URL.Query().Get("systemID")
		url, err := RedisClient.GetFeedURL(serviceID, "station_information", "en")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get system %q feed URL: %v", serviceID, err), 500)
			return
		}

		si, err := Client.LoadStationInformation(url)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get station information %q: %v", url, err), 500)
			return
		}

		fc := convertStationsToGeoJSON(si.Data.Stations)
		b, err := json.MarshalIndent(fc, "", "  ")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to marshal station information: %v", err), 500)
			return
		}

		_, err = w.Write(b)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to write response: %v", err), 500)
		}
	})
}

func convertStationsToGeoJSON(stations []gbfs.StationInformation) *gj.FeatureCollection {
	fc := gj.NewFeatureCollection()
	fc.Features = []*gj.Feature{}

	for _, station := range stations {
		props := stationProperties{
			ID:          station.ID,
			Name:        station.Name,
			Address:     station.Address,
			CrossStreet: station.CrossStreet,
			Capacity:    station.Capacity,
			ShortName:   station.ShortName,
			RegionID:    station.RegionID,
		}
		m, err := structmap.StructToMap(&props, "map", "")
		if err != nil {
			log.Printf("Failed to convert struct to map: %v", err)
			continue
		}

		feature := gj.Feature{
			Geometry: &gj.Geometry{
				Type: gj.GeometryPoint,
				Point: []float64{
					station.Lon,
					station.Lat,
				},
			},
			Properties: m,
		}
		fc.Features = append(fc.Features, &feature)
	}

	return fc
}
