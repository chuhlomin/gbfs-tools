package gbfs

import (
	"fmt"
	"net/http"

	"github.com/chuhlomin/gbfs-go"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/graphql-go/relay"
	"github.com/pkg/errors"

	"github.com/chuhlomin/gbfs-tools/pkg/redis"
	"github.com/chuhlomin/gbfs-tools/pkg/structs"
)

var Client *gbfs.Client
var RedisClient *redis.Client

var Schema graphql.Schema

func HandlerGraphQL() http.Handler {
	return handler.New(&handler.Config{
		Schema:     &Schema,
		Pretty:     true,
		GraphiQL:   true,
		Playground: true,
	})
}

func init() {
	// systemInformationType := graphql.NewObject(graphql.ObjectConfig{
	// 	Name:        "SystemInformation",
	// 	Description: "System information",
	// 	Fields: graphql.Fields{
	// 		"name": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "Name",
	// 		},
	// 		"shortName": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "ShortName",
	// 		},
	// 		"operator": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "Operator",
	// 		},
	// 		"url": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "URL",
	// 		},
	// 		"purchaseURL": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "Purchase URL",
	// 		},
	// 		"startDate": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "Start date",
	// 		},
	// 		"phoneNumber": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "Phone number",
	// 		},
	// 		"email": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "Email",
	// 		},
	// 		"feedContactEmail": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "Feed contact email",
	// 		},
	// 		"timezone": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "Timezone",
	// 		},
	// 		"licenseID": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "License ID",
	// 		},
	// 		"licenseURL": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "License URL",
	// 		},
	// 		"attributionOrganizationName": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "Attribution organization name",
	// 		},
	// 		"attributionURL": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "Attribution URL",
	// 		},
	// 		"language": &graphql.Field{
	// 			Type:        graphql.String,
	// 			Description: "Language",
	// 		},
	// 	},
	// })

	feedType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "Feed",
		Description: "GBFS Feed",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "Name",
			},
			"url": &graphql.Field{
				Type:        graphql.String,
				Description: "URL",
			},
			"language": &graphql.Field{
				Type:        graphql.String,
				Description: "Language",
			},
		},
	})

	systemType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "System",
		Description: "Bikeshare system",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.String,
				Description: "System ID",
			},
			"countryCode": &graphql.Field{
				Type:        graphql.String,
				Description: "Country Code",
			},
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "Name",
			},
			"location": &graphql.Field{
				Type:        graphql.String,
				Description: "Location",
			},
			"url": &graphql.Field{
				Type:        graphql.String,
				Description: "URL",
			},
			"autoDiscoveryUrl": &graphql.Field{
				Type:        graphql.String,
				Description: "Auto-discovery URL",
			},
			"languages": &graphql.Field{
				Type:        &graphql.List{OfType: graphql.String},
				Description: "Available GTFS languages",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source := p.Source
					switch t := source.(type) {
					case structs.System:
						system := source.(structs.System)
						return RedisClient.GetFeedsLanguages(system.ID)

					default:
						return nil, fmt.Errorf("Unexpected type %T in source: %v", t, p.Source)
					}
				},
			},
			"feeds": &graphql.Field{
				Type:        &graphql.List{OfType: feedType},
				Description: "SystemFeeds",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source := p.Source
					switch t := source.(type) {

					case *structs.System:
						system := source.(*structs.System)
						return RedisClient.GetFeeds(system.ID)
					default:
						return nil, fmt.Errorf("Unexpected type %T in source: %v", t, p.Source)
					}
				},
			},
		},
	})

	stationStatusType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "StationStatus",
		Description: "Station status",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.String,
				Description: "Identifier of a station",
			},
			"numBikesAvailable": &graphql.Field{
				Type:        graphql.Int,
				Description: "Number of functional vehicles physically at the station that may be offered for rental",
			},
			"numBikesDisabled": &graphql.Field{
				Type:        graphql.Int,
				Description: "Number of disabled vehicles of any type at the station",
			},
			"numDocksAvailable": &graphql.Field{
				Type:        graphql.Int,
				Description: "Number of functional docks physically at the station that are able to accept vehicles for return",
			},
			"isInstalled": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Is the station currently on the street?",
			},
			"isRenting": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Is the station currently renting vehicles?",
			},
			"isReturning": &graphql.Field{
				Type:        graphql.Boolean,
				Description: "Is the station accepting vehicle returns?",
			},
			"lastReported": &graphql.Field{
				Type:        graphql.DateTime,
				Description: "The last time this station reported its status to the operator's backend",
			},
		},
	})

	systemsConnectionDefinition := relay.ConnectionDefinitions(relay.ConnectionConfig{
		Name:     "System",
		NodeType: systemType,
	})

	systemsArgs := relay.ConnectionArgs
	systemsArgs["countryCode"] = &graphql.ArgumentConfig{
		Type: graphql.String,
	}

	stationStatusConnectionDefinition := relay.ConnectionDefinitions(relay.ConnectionConfig{
		Name:     "StationStatus",
		NodeType: stationStatusType,
	})

	stationStatusArgs := relay.ConnectionArgs
	stationStatusArgs["systemID"] = &graphql.ArgumentConfig{
		Type:        graphql.String,
		Description: "System ID",
	}

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"systems": &graphql.Field{
				Type: systemsConnectionDefinition.ConnectionType,
				Args: systemsArgs,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					args := relay.NewConnectionArguments(p.Args)

					countryCode, filterByCountryCode := p.Args["countryCode"]

					systems, err := RedisClient.GetSystems()
					if err != nil {
						return nil, err
					}
					var result []interface{}
					for i := range systems {
						if filterByCountryCode && systems[i].CountryCode != countryCode {
							// log.Printf("Filter out %s (%s)", systems[i].ID, systems[i].CountryCode)
							continue
						}

						result = append(result, systems[i])
					}

					return relay.ConnectionFromArray(result, args), nil
				},
			},
			"system": &graphql.Field{
				Type: systemType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "System ID",
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return RedisClient.GetSystem(fmt.Sprintf("%v", p.Args["id"]))
				},
			},
			"stationStatus": &graphql.Field{
				Type: stationStatusConnectionDefinition.ConnectionType,
				Args: stationStatusArgs,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					args := relay.NewConnectionArguments(p.Args)

					if _, ok := p.Args["systemID"]; !ok {
						return nil, fmt.Errorf("Missing systemID argument")
					}
					systemID := fmt.Sprintf("%v", p.Args["systemID"])

					stations, err := getStationStatus(systemID)
					if err != nil {
						return nil, err
					}

					var result []interface{}
					for i := range stations {
						result = append(result, stations[i])
					}

					return relay.ConnectionFromArray(result, args), nil
				},
			},
			// "system_information": &graphql.Field{
			// 	Type: systemInformationType,
			// 	Args: graphql.FieldConfigArgument{
			// 		"id": &graphql.ArgumentConfig{
			// 			Type:        graphql.String,
			// 			Description: "System ID",
			// 		},
			// 		"lang": &graphql.ArgumentConfig{
			// 			Type:        graphql.String,
			// 			Description: "Language",
			// 		},
			// 	},
			// 	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			// 		system, err := RedisClient.GetSystem(p.Args["id"])
			// 		systems, err := GetSystemInformation()
			// 		if err != nil {
			// 			return nil, err
			// 		}
			// 		for _, s := range systems {
			// 			if s.ID == p.Args["id"] {
			// 				return s, nil
			// 			}
			// 		}

			// 		return nil, fmt.Errorf("System %q not found", p.Args["id"])
			// 	},
			// },
		},
	})

	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	if err != nil {
		panic(err)
	}
}

func getStationStatus(systemID string) ([]gbfs.StationStatus, error) {
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
