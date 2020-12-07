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
			"_id": relay.GlobalIDField("System", nil),
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

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"systems": &graphql.Field{
				Type: systemsConnectionDefinition.ConnectionType,
				Args: relay.ConnectionArgs,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					args := relay.NewConnectionArguments(p.Args)

					systems, err := GetSystems()
					if err != nil {
						return nil, err
					}
					result := make([]interface{}, len(systems))
					for i := range systems {
						result[i] = systems[i]
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
