package gbfs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chuhlomin/gbfs-go"
	"github.com/graphql-go/handler"
	"github.com/pkg/errors"

	"github.com/chuhlomin/gbfs-tools/pkg/structs"
)

var client *gbfs.Client

// In-memory cache:
var systems []structs.System
var urls map[string]string                                  // id → url
var feeds map[string]*gbfs.LanguageFeeds                    // url → gbfs
var systemInformation map[string]*gbfs.SystemInformation    // url → system_information
var stationInformation map[string][]gbfs.StationInformation // url → []StationInformation

func init() {
	urls = map[string]string{}
	feeds = map[string]*gbfs.LanguageFeeds{}
	systemInformation = map[string]*gbfs.SystemInformation{}
	stationInformation = map[string][]gbfs.StationInformation{}

	client = gbfs.NewClient("github.com/chuhlomin/gbfs-tools", 30*time.Second)

	systemCSVPath, found := os.LookupEnv("SYSTEMS_CSV_PATH")
	if !found {
		// this case should be caught in main
		panic("Missing required environment variable SYSTEMS_CSV_PATH")
	}

	f, err := ioutil.ReadFile(systemCSVPath)
	if err != nil {
		panic("Failed to read file at SYSTEMS_CSV_PATH")
	}
	systemsRaw, err := gbfs.ParseSystemsCSV(bytes.NewReader(f))
	if err != nil {
		panic("Failed to parse systems CSV at SYSTEMS_CSV_PATH")
	}
	populateSystems(systemsRaw)
	go fetchSystemsFeeds()
}

func Handler() http.Handler {
	return handler.New(&handler.Config{
		Schema:     &Schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})
}

func GetSystems() ([]structs.System, error) {
	return systems, nil
}

func GetSystem(systemID string) (*structs.System, error) {
	systems, err := GetSystems()
	if err != nil {
		return nil, errors.Wrap(err, "get systems")
	}

	for _, s := range systems {
		if s.ID == systemID {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("system %q not found", systemID)
}

func GetStationStatus(systemID string) ([]gbfs.StationStatus, error) {
	system, err := GetSystem(systemID)
	if err != nil {
		return nil, errors.Wrapf(err, "get system %s", systemID)
	}

	gbfsInfo, err := GetGBFS(system.AutoDiscoveryURL)
	if err != nil {
		return nil, errors.Wrapf(err, "get GBFS %s", system.AutoDiscoveryURL)
	}

	feeds, err := gbfsInfo.GetDataFeeds("en")
	if err != nil {
		return nil, errors.Wrap(err, "get data feeds")
	}

	feed, err := feeds.GetFeed("station_status")
	if err != nil {
		return nil, errors.Wrap(err, "get feed")
	}

	status, err := client.LoadStationStatus(feed.URL)
	if err != nil {
		return nil, errors.Wrapf(err, "load station statis %s", feed.URL)
	}

	return status.Data.Stations, nil
}

func GetGBFS(url string) (*gbfs.LanguageFeeds, error) {
	var err error
	if f, ok := feeds[url]; ok {
		return f, nil
	}

	log.Printf("GET GBFS %q", url)
	resp, err := client.LoadGBFS(url)
	if err != nil {
		return nil, err
	}

	if resp.Data == nil {
		return nil, errors.New("empty response data")
	}

	feeds[url] = &resp.Data
	return &resp.Data, nil
}

func GetSystemInformation(url string) (*gbfs.SystemInformation, error) {
	var err error
	if s, ok := systemInformation[url]; ok {
		return s, nil
	}

	log.Printf("GET System information %q", url)
	resp, err := client.LoadSystemInformation(url)
	if err != nil {
		return nil, err
	}

	systemInformation[url] = &resp.Data
	return &resp.Data, nil
}

func populateSystems(raw []gbfs.System) {
	for _, s := range raw {
		systems = append(systems, structs.System{
			ID:               s.ID,
			CountryCode:      s.CountryCode,
			Name:             s.Name,
			Location:         s.Location,
			URL:              s.URL,
			AutoDiscoveryURL: s.AutoDiscoveryURL,
			IsEnabled:        true,
		})
		urls[s.ID] = s.AutoDiscoveryURL
	}
}

func fetchSystemsFeeds() {
	for systemID, gbfsURL := range urls {
		_, err := GetGBFS(gbfsURL)
		if err != nil {
			log.Printf("[ERROR] Failed to pre-fetch system feed for %q, URL %s: %v", systemID, gbfsURL, err)
		}
		time.Sleep(2 * time.Second)
	}
}
