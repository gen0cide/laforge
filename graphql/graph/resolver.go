package graph

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/graphql/graph/generated"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver Specify all the options that are able to be resolved here
// Resolver is the resolver root.
type Resolver struct{ client *ent.Client }

// NewSchema creates a graphql executable schema.
func NewSchema(client *ent.Client) graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: &Resolver{client},
	})
}
