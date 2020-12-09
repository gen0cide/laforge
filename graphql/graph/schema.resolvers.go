package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"time"

	"github.com/gen0cide/laforge/graphql/graph/generated"
	"github.com/gen0cide/laforge/graphql/graph/model"
)

func (r *mutationResolver) ExecutePlan(ctx context.Context, buildUUID string) (*model.Build, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Environments(ctx context.Context) ([]*model.Environment, error) {
	status := model.Status{
		State:     model.ProvisionStatusProvStatusComplete,
		StartedAt: time.Now().String(),
		EndedAt:   time.Now().String(),
		Failed:    false,
		Completed: true,
		Error:     "",
	}
	testEnviroment, _, _, _ := TestEnvironment(&status)
	r.enviroments = append(r.enviroments, testEnviroment)
	return r.enviroments, nil
}

func (r *queryResolver) Environment(ctx context.Context, envUUID string) (*model.Environment, error) {
	status := model.Status{
		State:     model.ProvisionStatusProvStatusComplete,
		StartedAt: time.Now().String(),
		EndedAt:   time.Now().String(),
		Failed:    false,
		Completed: true,
		Error:     "",
	}
	testEnviroment, _, _, _ := TestEnvironment(&status)

	if testEnviroment.ID == envUUID {
		return testEnviroment, nil
	}

	return nil, nil
}

func (r *queryResolver) ProvisionedHost(ctx context.Context, proHostUUID string) (*model.ProvisionedHost, error) {
	status := model.Status{
		State:     model.ProvisionStatusProvStatusComplete,
		StartedAt: time.Now().String(),
		EndedAt:   time.Now().String(),
		Failed:    false,
		Completed: true,
		Error:     "",
	}
	_, _, testHosts, _ := TestEnvironment(&status)
	for _, host := range testHosts {
		if host.ID == proHostUUID {
			return host, nil
		}
	}
	return nil, nil
}

func (r *queryResolver) ProvisionedNetwork(ctx context.Context, proNetUUID string) (*model.ProvisionedNetwork, error) {
	status := model.Status{
		State:     model.ProvisionStatusProvStatusComplete,
		StartedAt: time.Now().String(),
		EndedAt:   time.Now().String(),
		Failed:    false,
		Completed: true,
		Error:     "",
	}
	_, testNetworks, _, _ := TestEnvironment(&status)
	for _, network := range testNetworks {
		if network.ID == proNetUUID {
			return network, nil
		}
	}
	return nil, nil
}

func (r *queryResolver) ProvisionedStep(ctx context.Context, proStepUUID string) (*model.ProvisionedStep, error) {
	status := model.Status{
		State:     model.ProvisionStatusProvStatusComplete,
		StartedAt: time.Now().String(),
		EndedAt:   time.Now().String(),
		Failed:    false,
		Completed: true,
		Error:     "",
	}
	_, _, _, testSteps := TestEnvironment(&status)
	for _, step := range testSteps {
		if step.ID == proStepUUID {
			return step, nil
		}
	}
	return nil, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
