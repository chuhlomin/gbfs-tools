package gbfs

import (
	"fmt"

	"github.com/chuhlomin/gbfs-go"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"github.com/pkg/errors"

	"github.com/chuhlomin/gbfs-tools/pkg/structs"
)

var systemInformationType *graphql.Object
var systemType *graphql.Object

var Schema graphql.Schema

func init() {
	systemInformationType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "SystemInformation",
		Description: "System information",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "Name",
			},
			"shortName": &graphql.Field{
				Type:        graphql.String,
				Description: "ShortName",
			},
			"operator": &graphql.Field{
				Type:        graphql.String,
				Description: "Operator",
			},
			"url": &graphql.Field{
				Type:        graphql.String,
				Description: "URL",
			},
			"purchaseURL": &graphql.Field{
				Type:        graphql.String,
				Description: "Purchase URL",
			},
			"startDate": &graphql.Field{
				Type:        graphql.String,
				Description: "Start date",
			},
			"phoneNumber": &graphql.Field{
				Type:        graphql.String,
				Description: "Phone number",
			},
			"email": &graphql.Field{
				Type:        graphql.String,
				Description: "Email",
			},
			"feedContactEmail": &graphql.Field{
				Type:        graphql.String,
				Description: "Feed contact email",
			},
			"timezone": &graphql.Field{
				Type:        graphql.String,
				Description: "Timezone",
			},
			"licenseID": &graphql.Field{
				Type:        graphql.String,
				Description: "License ID",
			},
			"licenseURL": &graphql.Field{
				Type:        graphql.String,
				Description: "License URL",
			},
			"attributionOrganizationName": &graphql.Field{
				Type:        graphql.String,
				Description: "Attribution organization name",
			},
			"attributionURL": &graphql.Field{
				Type:        graphql.String,
				Description: "Attribution URL",
			},
			"language": &graphql.Field{
				Type:        graphql.String,
				Description: "Language",
			},
		},
	})

	systemType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "System",
		Description: "Bikeshare system",
		Fields: graphql.Fields{
			// "id": relay.GlobalIDField("System", nil),
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
					case gbfs.System:
						system := source.(gbfs.System)
						lf, err := GetGBFS(system.AutoDiscoveryURL)
						if err != nil {
							return nil, fmt.Errorf("Failed to get GBFS from %q: %v", system.AutoDiscoveryURL, err)
						}

						var result []string
						for l := range *lf {
							result = append(result, l)
						}
						return result, nil

					default:
						return nil, fmt.Errorf("Unexpected type %T in source: %v", t, p.Source)
					}
				},
			},
			"information": &graphql.Field{
				Type:        systemInformationType,
				Description: "System information",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					source := p.Source
					switch t := source.(type) {
					case structs.System:
						system := source.(structs.System)
						lf, err := GetGBFS(system.AutoDiscoveryURL)
						if err != nil {
							return nil, errors.Wrapf(err, "Failed to get GBFS from %q", system.AutoDiscoveryURL)
						}

						feeds, ok := (*lf)["en"]
						if !ok {
							for _, feeds = range *lf {
								break
							}
						}

						feed, err := findFeed(feeds.Feeds, "system_information")
						if err != nil {
							return nil, err
						}

						result, err := GetSystemInformation(feed.URL)
						if err != nil {
							return nil, err
						}

						return result, nil

					default:
						return nil, fmt.Errorf("Unexpected type %T in source: %v", t, p.Source)
					}
				},
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

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"systems": &graphql.Field{
				Type: systemsConnectionDefinition.ConnectionType,
				Args: systemsArgs,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					args := relay.NewConnectionArguments(p.Args)

					countryCode, filterByCountryCode := p.Args["countryCode"]

					systems, err := GetSystems()
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
					return GetSystem(fmt.Sprintf("%v", p.Args["id"]))
				},
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addSystem": &graphql.Field{
				Type: systemType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"countryCode": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"location": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"url": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"autoDiscoveryUrl": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(string)
					name, _ := params.Args["name"].(string)
					countryCode, _ := params.Args["countryCode"].(string)
					location, _ := params.Args["location"].(string)
					url, _ := params.Args["url"].(string)
					autoDiscoveryURL, _ := params.Args["autoDiscoveryUrl"].(string)

					system := structs.System{
						ID:               id,
						Name:             name,
						CountryCode:      countryCode,
						Location:         location,
						URL:              url,
						AutoDiscoveryURL: autoDiscoveryURL,
					}
					err := AddSystem(system)
					return system, err
				},
			},
			"disableSystem": &graphql.Field{
				Type: systemType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id, _ := params.Args["id"].(string)
					err := DisableSystem(id)
					return nil, err
				},
			},
		},
	})
	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		panic(err)
	}
}

func findFeed(feeds []gbfs.Feed, name string) (gbfs.Feed, error) {
	for _, f := range feeds {
		if f.Name == name {
			return f, nil
		}
	}

	return gbfs.Feed{}, fmt.Errorf("feed %q not found", name)
}
