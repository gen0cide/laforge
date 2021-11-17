package graph

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/graphql/auth"
	"github.com/gen0cide/laforge/graphql/graph/generated"
	"github.com/gen0cide/laforge/graphql/graph/model"
	"github.com/go-redis/redis/v8"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver Specify all the options that are able to be resolved here
// Resolver is the resolver root.
type Resolver struct {
	client *ent.Client
	rdb    *redis.Client
}

// NewSchema creates a graphql executable schema.
func NewSchema(client *ent.Client, rdb *redis.Client) graphql.ExecutableSchema {
	GQLConfig := generated.Config{
		Resolvers: &Resolver{
			client: client,
			rdb:    rdb,
		},
	}
	GQLConfig.Directives.HasRole = func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []model.RoleLevel) (res interface{}, err error) {
		currentUser, err := auth.ForContext(ctx)

		if err != nil {
			return nil, err
		}

		for _, role := range roles {
			if role.String() == string(currentUser.Role) {
				return next(ctx)
			}
		}
		return nil, &gqlerror.Error{
			Message: "not authorized",
			Extensions: map[string]interface{}{
				"code": "401",
			},
		}

	}
	return generated.NewExecutableSchema(GQLConfig)
}
