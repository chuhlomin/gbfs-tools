package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/chuhlomin/gbfs-go"
	"github.com/chuhlomin/gbfs-tools/pkg/structs"
	"github.com/pkg/errors"
)

func main() {
	log.Println("Starting...")

	if err := run(); err != nil {
		log.Fatalf("ERROR: Failed to %v", err)
	}

	log.Println("Stopped.")
}

const query = `mutation(
		$id: String,
		$name: String,
		$countryCode: String,
		$location: String,
		$url: String,
		$autoDiscoveryUrl: String
) {
	addSystem(
		id: $id,
		name: $name,
		countryCode: $countryCode,
		location: $location,
		url: $url,
		autoDiscoveryUrl: $autoDiscoveryUrl
	){id name countryCode location url autoDiscoveryUrl}
}`

type request struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

func run() error {
	c := gbfs.NewClient("github.com/chuhlomin/gbfs-tools", 30*time.Second)
	systems, err := c.LoadSystems(gbfs.SystemsNABSA)
	if err != nil {
		return errors.Wrapf(err, "load systems from %q", gbfs.SystemsNABSA)
	}

	httpClient := http.Client{
		Timeout: 30 * time.Second,
	}

	for _, s := range systems {
		b, err := json.Marshal(request{
			Query: query,
			Variables: structs.System{
				ID:               s.ID,
				CountryCode:      s.CountryCode,
				Name:             s.Name,
				Location:         s.Location,
				URL:              s.URL,
				AutoDiscoveryURL: s.AutoDiscoveryURL,
			},
		})
		if err != nil {
			return errors.Wrap(err, "marshal request")
		}

		response, err := httpClient.Post(
			"http://127.0.0.1:8082/graphql",
			"application/json",
			bytes.NewReader(b),
		)
		if err != nil {
			return errors.Wrap(err, "post request")
		}
		log.Printf("Success: %v", response)
	}

	return nil
}
