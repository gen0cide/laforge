package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/build"
	"github.com/gen0cide/laforge/ent/environment"
	"github.com/gen0cide/laforge/ent/plan"
	"github.com/gen0cide/laforge/ent/provisionedhost"
	"github.com/gen0cide/laforge/ent/provisionednetwork"
	"github.com/gen0cide/laforge/ent/provisioningstep"
	"github.com/gen0cide/laforge/graphql/graph/generated"
	"github.com/gen0cide/laforge/graphql/graph/model"
	"github.com/gen0cide/laforge/loader"
	"github.com/gen0cide/laforge/planner"
)

func (r *commandResolver) Vars(ctx context.Context, obj *ent.Command) ([]*model.VarsMap, error) {
	results := make([]*model.VarsMap, 0)
	for varKey, varValue := range obj.Vars {
		tempVar := &model.VarsMap{
			Key:   varKey,
			Value: varValue,
		}
		results = append(results, tempVar)
	}
	return results, nil
}

func (r *commandResolver) Tags(ctx context.Context, obj *ent.Command) ([]*model.TagMap, error) {
	results := make([]*model.TagMap, 0)
	for tagKey, tagValue := range obj.Tags {
		tempTag := &model.TagMap{
			Key:   tagKey,
			Value: tagValue,
		}
		results = append(results, tempTag)
	}
	return results, nil
}

func (r *competitionResolver) Config(ctx context.Context, obj *ent.Competition) ([]*model.ConfigMap, error) {
	results := make([]*model.ConfigMap, 0)
	for configKey, configValue := range obj.Config {
		configTag := &model.ConfigMap{
			Key:   configKey,
			Value: configValue,
		}
		results = append(results, configTag)
	}
	return results, nil
}

func (r *competitionResolver) Tags(ctx context.Context, obj *ent.Competition) ([]*model.TagMap, error) {
	results := make([]*model.TagMap, 0)
	for tagKey, tagValue := range obj.Tags {
		tempTag := &model.TagMap{
			Key:   tagKey,
			Value: tagValue,
		}
		results = append(results, tempTag)
	}
	return results, nil
}

func (r *dNSResolver) Config(ctx context.Context, obj *ent.DNS) ([]*model.ConfigMap, error) {
	results := make([]*model.ConfigMap, 0)
	for configKey, configValue := range obj.Config {
		configTag := &model.ConfigMap{
			Key:   configKey,
			Value: configValue,
		}
		results = append(results, configTag)
	}
	return results, nil
}

func (r *dNSRecordResolver) Vars(ctx context.Context, obj *ent.DNSRecord) ([]*model.VarsMap, error) {
	results := make([]*model.VarsMap, 0)
	for varKey, varValue := range obj.Vars {
		tempVar := &model.VarsMap{
			Key:   varKey,
			Value: varValue,
		}
		results = append(results, tempVar)
	}
	return results, nil
}

func (r *dNSRecordResolver) Tags(ctx context.Context, obj *ent.DNSRecord) ([]*model.TagMap, error) {
	results := make([]*model.TagMap, 0)
	for tagKey, tagValue := range obj.Tags {
		tempTag := &model.TagMap{
			Key:   tagKey,
			Value: tagValue,
		}
		results = append(results, tempTag)
	}
	return results, nil
}

func (r *environmentResolver) Config(ctx context.Context, obj *ent.Environment) ([]*model.ConfigMap, error) {
	results := make([]*model.ConfigMap, 0)
	for configKey, configValue := range obj.Config {
		configTag := &model.ConfigMap{
			Key:   configKey,
			Value: configValue,
		}
		results = append(results, configTag)
	}
	return results, nil
}

func (r *environmentResolver) Tags(ctx context.Context, obj *ent.Environment) ([]*model.TagMap, error) {
	results := make([]*model.TagMap, 0)
	for tagKey, tagValue := range obj.Tags {
		tempTag := &model.TagMap{
			Key:   tagKey,
			Value: tagValue,
		}
		results = append(results, tempTag)
	}
	return results, nil
}

func (r *fileDeleteResolver) Tags(ctx context.Context, obj *ent.FileDelete) ([]*model.TagMap, error) {
	results := make([]*model.TagMap, 0)
	for tagKey, tagValue := range obj.Tags {
		tempTag := &model.TagMap{
			Key:   tagKey,
			Value: tagValue,
		}
		results = append(results, tempTag)
	}
	return results, nil
}

func (r *fileDownloadResolver) Tags(ctx context.Context, obj *ent.FileDownload) ([]*model.TagMap, error) {
	results := make([]*model.TagMap, 0)
	for tagKey, tagValue := range obj.Tags {
		tempTag := &model.TagMap{
			Key:   tagKey,
			Value: tagValue,
		}
		results = append(results, tempTag)
	}
	return results, nil
}

func (r *fileExtractResolver) Tags(ctx context.Context, obj *ent.FileExtract) ([]*model.TagMap, error) {
	results := make([]*model.TagMap, 0)
	for tagKey, tagValue := range obj.Tags {
		tempTag := &model.TagMap{
			Key:   tagKey,
			Value: tagValue,
		}
		results = append(results, tempTag)
	}
	return results, nil
}

func (r *findingResolver) Severity(ctx context.Context, obj *ent.Finding) (model.FindingSeverity, error) {
	return model.FindingSeverity(obj.Severity), nil
}

func (r *findingResolver) Difficulty(ctx context.Context, obj *ent.Finding) (model.FindingDifficulty, error) {
	return model.FindingDifficulty(obj.Difficulty), nil
}

func (r *findingResolver) Tags(ctx context.Context, obj *ent.Finding) ([]*model.TagMap, error) {
	results := make([]*model.TagMap, 0)
	for tagKey, tagValue := range obj.Tags {
		tempTag := &model.TagMap{
			Key:   tagKey,
			Value: tagValue,
		}
		results = append(results, tempTag)
	}
	return results, nil
}

func (r *hostResolver) Vars(ctx context.Context, obj *ent.Host) ([]*model.VarsMap, error) {
	results := make([]*model.VarsMap, 0)
	for varKey, varValue := range obj.Vars {
		tempVar := &model.VarsMap{
			Key:   varKey,
			Value: varValue,
		}
		results = append(results, tempVar)
	}
	return results, nil
}

func (r *hostResolver) Tags(ctx context.Context, obj *ent.Host) ([]*model.TagMap, error) {
	results := make([]*model.TagMap, 0)
	for tagKey, tagValue := range obj.Tags {
		tempTag := &model.TagMap{
			Key:   tagKey,
			Value: tagValue,
		}
		results = append(results, tempTag)
	}
	return results, nil
}

func (r *identityResolver) Vars(ctx context.Context, obj *ent.Identity) ([]*model.VarsMap, error) {
	results := make([]*model.VarsMap, 0)
	for varKey, varValue := range obj.Vars {
		tempVar := &model.VarsMap{
			Key:   varKey,
			Value: varValue,
		}
		results = append(results, tempVar)
	}
	return results, nil
}

func (r *identityResolver) Tags(ctx context.Context, obj *ent.Identity) ([]*model.TagMap, error) {
	results := make([]*model.TagMap, 0)
	for tagKey, tagValue := range obj.Tags {
		tempTag := &model.TagMap{
			Key:   tagKey,
			Value: tagValue,
		}
		results = append(results, tempTag)
	}
	return results, nil
}

func (r *mutationResolver) LoadEnviroment(ctx context.Context, envFilePath string) ([]*ent.Environment, error) {
	return loader.LoadEnviroment(ctx, r.client, envFilePath)
}

func (r *mutationResolver) CreateBuild(ctx context.Context, envUUID string) (*ent.Build, error) {
	uuid, err := strconv.Atoi(envUUID)

	if err != nil {
		return nil, fmt.Errorf("failed casting UUID to int: %v", err)
	}

	entEnvironment, err := r.client.Environment.Query().Where(environment.IDEQ(uuid)).Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Environment: %v", err)
	}
	planner.RenderFiles = true

	return planner.CreateBuild(ctx, r.client, entEnvironment)
}

func (r *mutationResolver) ExecutePlan(ctx context.Context, buildUUID string) (*ent.Build, error) {
	uuid, err := strconv.Atoi(buildUUID)

	if err != nil {
		return nil, fmt.Errorf("failed casting UUID to int: %v", err)
	}

	b, err := r.client.Build.Query().Where(build.IDEQ(uuid)).Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Build: %v", err)
	}

	return b, nil
}

func (r *networkResolver) Vars(ctx context.Context, obj *ent.Network) ([]*model.VarsMap, error) {
	results := make([]*model.VarsMap, 0)
	for varKey, varValue := range obj.Vars {
		tempVar := &model.VarsMap{
			Key:   varKey,
			Value: varValue,
		}
		results = append(results, tempVar)
	}
	return results, nil
}

func (r *networkResolver) Tags(ctx context.Context, obj *ent.Network) ([]*model.TagMap, error) {
	results := make([]*model.TagMap, 0)
	for tagKey, tagValue := range obj.Tags {
		tempTag := &model.TagMap{
			Key:   tagKey,
			Value: tagValue,
		}
		results = append(results, tempTag)
	}
	return results, nil
}

func (r *planResolver) Type(ctx context.Context, obj *ent.Plan) (model.PlanType, error) {
	return model.PlanType(obj.Type), nil
}

func (r *provisionedHostResolver) CombinedOutput(ctx context.Context, obj *ent.ProvisionedHost) (*string, error) {
	// TODO: Implement CombinedOutput
	todo := "not implemented"
	return &todo, nil
}

func (r *provisionedHostResolver) ProvisionedHostToAgentStatus(ctx context.Context, obj *ent.ProvisionedHost) (*ent.AgentStatus, error) {
	check, err := obj.QueryProvisionedHostToAgentStatus().Exist(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Agent Status: %v", err)
	}

	if check {
		a, err := obj.QueryProvisionedHostToAgentStatus().Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed querying Agent Status: %v", err)
		}
		return a, nil
	}

	return nil, nil
}

func (r *provisioningStepResolver) Type(ctx context.Context, obj *ent.ProvisioningStep) (model.ProvisioningStepType, error) {
	return model.ProvisioningStepType(obj.Type), nil
}

func (r *queryResolver) Environments(ctx context.Context) ([]*ent.Environment, error) {
	e, err := r.client.Environment.Query().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Environment: %v", err)
	}

	return e, nil
}

func (r *queryResolver) Environment(ctx context.Context, envUUID string) (*ent.Environment, error) {
	uuid, err := strconv.Atoi(envUUID)

	if err != nil {
		return nil, fmt.Errorf("failed casting UUID to int: %v", err)
	}

	e, err := r.client.Environment.Query().Where(environment.IDEQ(uuid)).Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Environment: %v", err)
	}

	return e, nil
}

func (r *queryResolver) ProvisionedHost(ctx context.Context, proHostUUID string) (*ent.ProvisionedHost, error) {
	uuid, err := strconv.Atoi(proHostUUID)

	if err != nil {
		return nil, fmt.Errorf("failed casting UUID to int: %v", err)
	}

	ph, err := r.client.ProvisionedHost.Query().Where(provisionedhost.IDEQ(uuid)).Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying ProvisionedHost: %v", err)
	}

	return ph, nil
}

func (r *queryResolver) ProvisionedNetwork(ctx context.Context, proNetUUID string) (*ent.ProvisionedNetwork, error) {
	uuid, err := strconv.Atoi(proNetUUID)

	if err != nil {
		return nil, fmt.Errorf("failed casting UUID to int: %v", err)
	}

	pn, err := r.client.ProvisionedNetwork.Query().Where(provisionednetwork.IDEQ(uuid)).Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying ProvisionedNetwork: %v", err)
	}

	return pn, nil
}

func (r *queryResolver) ProvisionedStep(ctx context.Context, proStepUUID string) (*ent.ProvisioningStep, error) {
	uuid, err := strconv.Atoi(proStepUUID)

	if err != nil {
		return nil, fmt.Errorf("failed casting UUID to int: %v", err)
	}

	ps, err := r.client.ProvisioningStep.Query().Where(provisioningstep.IDEQ(uuid)).Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying ProvisionedStep: %v", err)
	}

	return ps, nil
}

func (r *queryResolver) Plan(ctx context.Context, planUUID string) (*ent.Plan, error) {
	uuid, err := strconv.Atoi(planUUID)

	if err != nil {
		return nil, fmt.Errorf("failed casting UUID to int: %v", err)
	}

	plan, err := r.client.Plan.Query().Where(plan.IDEQ(uuid)).Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying ProvisionedNetwork: %v", err)
	}

	return plan, nil
}

func (r *queryResolver) Build(ctx context.Context, buildUUID string) (*ent.Build, error) {
	uuid, err := strconv.Atoi(buildUUID)

	if err != nil {
		return nil, fmt.Errorf("failed casting UUID to int: %v", err)
	}

	build, err := r.client.Build.Query().Where(build.IDEQ(uuid)).Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying ProvisionedNetwork: %v", err)
	}

	return build, nil
}

func (r *scriptResolver) Vars(ctx context.Context, obj *ent.Script) ([]*model.VarsMap, error) {
	results := make([]*model.VarsMap, 0)
	for varKey, varValue := range obj.Vars {
		tempVar := &model.VarsMap{
			Key:   varKey,
			Value: varValue,
		}
		results = append(results, tempVar)
	}
	return results, nil
}

func (r *scriptResolver) Tags(ctx context.Context, obj *ent.Script) ([]*model.TagMap, error) {
	results := make([]*model.TagMap, 0)
	for tagKey, tagValue := range obj.Tags {
		tempTag := &model.TagMap{
			Key:   tagKey,
			Value: tagValue,
		}
		results = append(results, tempTag)
	}
	return results, nil
}

func (r *statusResolver) State(ctx context.Context, obj *ent.Status) (model.ProvisionStatus, error) {
	return model.ProvisionStatus(obj.State), nil
}

func (r *statusResolver) StatusFor(ctx context.Context, obj *ent.Status) (model.ProvisionStatusFor, error) {
	return model.ProvisionStatusFor(obj.StatusFor), nil
}

func (r *statusResolver) StartedAt(ctx context.Context, obj *ent.Status) (string, error) {
	return obj.StartedAt.String(), nil
}

func (r *statusResolver) EndedAt(ctx context.Context, obj *ent.Status) (string, error) {
	return obj.EndedAt.String(), nil
}

// Command returns generated.CommandResolver implementation.
func (r *Resolver) Command() generated.CommandResolver { return &commandResolver{r} }

// Competition returns generated.CompetitionResolver implementation.
func (r *Resolver) Competition() generated.CompetitionResolver { return &competitionResolver{r} }

// DNS returns generated.DNSResolver implementation.
func (r *Resolver) DNS() generated.DNSResolver { return &dNSResolver{r} }

// DNSRecord returns generated.DNSRecordResolver implementation.
func (r *Resolver) DNSRecord() generated.DNSRecordResolver { return &dNSRecordResolver{r} }

// Environment returns generated.EnvironmentResolver implementation.
func (r *Resolver) Environment() generated.EnvironmentResolver { return &environmentResolver{r} }

// FileDelete returns generated.FileDeleteResolver implementation.
func (r *Resolver) FileDelete() generated.FileDeleteResolver { return &fileDeleteResolver{r} }

// FileDownload returns generated.FileDownloadResolver implementation.
func (r *Resolver) FileDownload() generated.FileDownloadResolver { return &fileDownloadResolver{r} }

// FileExtract returns generated.FileExtractResolver implementation.
func (r *Resolver) FileExtract() generated.FileExtractResolver { return &fileExtractResolver{r} }

// Finding returns generated.FindingResolver implementation.
func (r *Resolver) Finding() generated.FindingResolver { return &findingResolver{r} }

// Host returns generated.HostResolver implementation.
func (r *Resolver) Host() generated.HostResolver { return &hostResolver{r} }

// Identity returns generated.IdentityResolver implementation.
func (r *Resolver) Identity() generated.IdentityResolver { return &identityResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Network returns generated.NetworkResolver implementation.
func (r *Resolver) Network() generated.NetworkResolver { return &networkResolver{r} }

// Plan returns generated.PlanResolver implementation.
func (r *Resolver) Plan() generated.PlanResolver { return &planResolver{r} }

// ProvisionedHost returns generated.ProvisionedHostResolver implementation.
func (r *Resolver) ProvisionedHost() generated.ProvisionedHostResolver {
	return &provisionedHostResolver{r}
}

// ProvisioningStep returns generated.ProvisioningStepResolver implementation.
func (r *Resolver) ProvisioningStep() generated.ProvisioningStepResolver {
	return &provisioningStepResolver{r}
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Script returns generated.ScriptResolver implementation.
func (r *Resolver) Script() generated.ScriptResolver { return &scriptResolver{r} }

// Status returns generated.StatusResolver implementation.
func (r *Resolver) Status() generated.StatusResolver { return &statusResolver{r} }

type commandResolver struct{ *Resolver }
type competitionResolver struct{ *Resolver }
type dNSResolver struct{ *Resolver }
type dNSRecordResolver struct{ *Resolver }
type environmentResolver struct{ *Resolver }
type fileDeleteResolver struct{ *Resolver }
type fileDownloadResolver struct{ *Resolver }
type fileExtractResolver struct{ *Resolver }
type findingResolver struct{ *Resolver }
type hostResolver struct{ *Resolver }
type identityResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type networkResolver struct{ *Resolver }
type planResolver struct{ *Resolver }
type provisionedHostResolver struct{ *Resolver }
type provisioningStepResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type scriptResolver struct{ *Resolver }
type statusResolver struct{ *Resolver }
