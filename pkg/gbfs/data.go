package gbfs

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chuhlomin/gbfs-go"
	"github.com/graphql-go/handler"
	"github.com/pkg/errors"

	"github.com/chuhlomin/gbfs-tools/pkg/database"
	"github.com/chuhlomin/gbfs-tools/pkg/structs"
)

var client *gbfs.Client
var db *database.Bolt

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

	dbPath, found := os.LookupEnv("DB_PATH")
	if !found {
		// this case should be caught in main
		panic("Missing required environment variable DB_PATH")
	}
	var err error
	db, err = database.NewBolt(dbPath)
	if err != nil {
		panic(err)
	}
}

func Handler() http.Handler {
	return handler.New(&handler.Config{
		Schema:     &Schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})
}

func AddSystem(system structs.System) error {
	return db.AddSystem(system)
}

func DisableSystem(id string) error {
	return db.DisableSystem(id)
}

func GetSystems() ([]structs.System, error) {
	var err error
	if systems == nil {
		systems, err = db.GetSystems()
		for _, s := range systems {
			urls[s.ID] = s.AutoDiscoveryURL
		}
	}

	return systems, err
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
