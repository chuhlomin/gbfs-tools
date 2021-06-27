package gbfs

import (
	"net/http"

	"github.com/chuhlomin/gbfs-go"
	"github.com/graphql-go/handler"
	"github.com/pkg/errors"

	"github.com/chuhlomin/gbfs-tools/pkg/redis"
	"github.com/chuhlomin/gbfs-tools/pkg/structs"
)

var Client *gbfs.Client
var RedisClient *redis.Client

func HandlerGraphQL() http.Handler {
	return handler.New(&handler.Config{
		Schema:     &Schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})
}

func GetSystems() ([]*structs.System, error) {
	var systems []*structs.System
	systems, err := RedisClient.GetSystems()
	if err != nil {
		return nil, errors.Wrap(err, "get systems")
	}

	return systems, nil
}

func GetSystem(systemID string) (*structs.System, error) {
	return RedisClient.GetSystem(systemID)
}

func GetStationStatus(systemID string) ([]gbfs.StationStatus, error) {
	url, err := RedisClient.GetFeedURL(systemID, "station_status", "en")
	if err != nil {
		return nil, errors.Wrapf(err, "get station status for %q", systemID)
	}

	status, err := Client.LoadStationStatus(url)
	if err != nil {
		return nil, errors.Wrapf(err, "load station statis %q", url)
	}

	return status.Data.Stations, nil
}
