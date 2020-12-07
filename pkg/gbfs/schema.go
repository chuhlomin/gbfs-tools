package gbfs

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
)

var systemType *graphql.Object

var Schema graphql.Schema

func init() {
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
