package graphql

import (
	"github.com/graphql-go/graphql"
)

// GetSchema returns the GraphQL schema
func GetSchema() (graphql.Schema, error) {
	// Define basic types
	streamType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Stream",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	// Define Query type
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hello": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return "Hello from StreamHub API! 🚀", nil
				},
			},
			"streams": &graphql.Field{
				Type: graphql.NewList(streamType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// Return empty array for now
					return []interface{}{}, nil
				},
			},
		},
	})

	// Create schema
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
}
