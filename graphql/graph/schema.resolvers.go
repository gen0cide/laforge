package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"log"

	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/graphql/graph/generated"
	"github.com/gen0cide/laforge/graphql/graph/model"
)

func (r *buildResolver) Tags(ctx context.Context, obj *ent.Build) ([]*ent.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *buildResolver) Config(ctx context.Context, obj *ent.Build) ([]*model.ConfigMap, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *buildResolver) Maintainer(ctx context.Context, obj *ent.Build) (*ent.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *buildResolver) Teams(ctx context.Context, obj *ent.Build) ([]*ent.Team, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *commandResolver) Vars(ctx context.Context, obj *ent.Command) ([]*model.VarsMap, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *commandResolver) Tags(ctx context.Context, obj *ent.Command) ([]*ent.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *commandResolver) Maintainer(ctx context.Context, obj *ent.Command) (*ent.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *competitionResolver) Config(ctx context.Context, obj *ent.Competition) ([]*model.ConfigMap, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *competitionResolver) DNS(ctx context.Context, obj *ent.Competition) (*ent.DNS, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *dNSResolver) NTPServer(ctx context.Context, obj *ent.DNS) ([]*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *dNSResolver) Config(ctx context.Context, obj *ent.DNS) ([]*model.ConfigMap, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *dNSRecordResolver) Vars(ctx context.Context, obj *ent.DNSRecord) ([]*model.VarsMap, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *dNSRecordResolver) Tags(ctx context.Context, obj *ent.DNSRecord) ([]*ent.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *environmentResolver) Tags(ctx context.Context, obj *ent.Environment) ([]*ent.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *environmentResolver) Config(ctx context.Context, obj *ent.Environment) ([]*model.ConfigMap, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *environmentResolver) Maintainer(ctx context.Context, obj *ent.Environment) (*ent.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *environmentResolver) Networks(ctx context.Context, obj *ent.Environment) ([]*ent.Network, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *environmentResolver) Hosts(ctx context.Context, obj *ent.Environment) ([]*ent.Host, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *environmentResolver) Build(ctx context.Context, obj *ent.Environment) (*ent.Build, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *environmentResolver) Competition(ctx context.Context, obj *ent.Environment) (*ent.Competition, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *fileDownloadResolver) Templete(ctx context.Context, obj *ent.FileDownload) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *fileDownloadResolver) Tags(ctx context.Context, obj *ent.FileDownload) ([]*ent.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *findingResolver) Severity(ctx context.Context, obj *ent.Finding) (model.FindingSeverity, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *findingResolver) Difficulty(ctx context.Context, obj *ent.Finding) (model.FindingDifficulty, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *findingResolver) Maintainer(ctx context.Context, obj *ent.Finding) (*ent.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *findingResolver) Tags(ctx context.Context, obj *ent.Finding) ([]*ent.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *findingResolver) Host(ctx context.Context, obj *ent.Finding) (*ent.Host, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *hostResolver) Vars(ctx context.Context, obj *ent.Host) ([]*model.VarsMap, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *hostResolver) DependsOn(ctx context.Context, obj *ent.Host) ([]*ent.Host, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *hostResolver) Maintainer(ctx context.Context, obj *ent.Host) (*ent.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *hostResolver) Tags(ctx context.Context, obj *ent.Host) ([]*ent.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *hostResolver) DNSRecords(ctx context.Context, obj *ent.Host) ([]*ent.DNSRecord, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *hostResolver) Commands(ctx context.Context, obj *ent.Host) ([]*ent.Command, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *hostResolver) Disk(ctx context.Context, obj *ent.Host) (*ent.Disk, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *hostResolver) Scripts(ctx context.Context, obj *ent.Host) ([]*ent.Script, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *hostResolver) FileDeletes(ctx context.Context, obj *ent.Host) ([]*ent.FileDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *hostResolver) FileDownloads(ctx context.Context, obj *ent.Host) ([]*ent.FileDownload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *hostResolver) FileExtracts(ctx context.Context, obj *ent.Host) ([]*ent.FileExtract, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ExecutePlan(ctx context.Context, buildUUID string) (*ent.Build, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *networkResolver) Vars(ctx context.Context, obj *ent.Network) ([]*model.VarsMap, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *networkResolver) Tags(ctx context.Context, obj *ent.Network) ([]*ent.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisionedHostResolver) Status(ctx context.Context, obj *ent.ProvisionedHost) (*ent.Status, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisionedHostResolver) ProvisionedNetwork(ctx context.Context, obj *ent.ProvisionedHost) (*ent.ProvisionedNetwork, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisionedHostResolver) Host(ctx context.Context, obj *ent.ProvisionedHost) (*ent.Host, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisionedHostResolver) CombinedOutput(ctx context.Context, obj *ent.ProvisionedHost) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisionedHostResolver) Heartbeat(ctx context.Context, obj *ent.ProvisionedHost) (*ent.AgentStatus, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisionedNetworkResolver) Vars(ctx context.Context, obj *ent.ProvisionedNetwork) ([]*model.VarsMap, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisionedNetworkResolver) Tags(ctx context.Context, obj *ent.ProvisionedNetwork) ([]*ent.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisionedNetworkResolver) Status(ctx context.Context, obj *ent.ProvisionedNetwork) (*ent.Status, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisionedNetworkResolver) Network(ctx context.Context, obj *ent.ProvisionedNetwork) (*ent.Network, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisionedNetworkResolver) Build(ctx context.Context, obj *ent.ProvisionedNetwork) (*ent.Build, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisioningStepResolver) ProvisionType(ctx context.Context, obj *ent.ProvisioningStep) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisioningStepResolver) ProvisionedHost(ctx context.Context, obj *ent.ProvisioningStep) (*ent.ProvisionedHost, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisioningStepResolver) Status(ctx context.Context, obj *ent.ProvisioningStep) (*ent.Status, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisioningStepResolver) Script(ctx context.Context, obj *ent.ProvisioningStep) (*ent.Script, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisioningStepResolver) Command(ctx context.Context, obj *ent.ProvisioningStep) (*ent.Command, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisioningStepResolver) DNSRecord(ctx context.Context, obj *ent.ProvisioningStep) (*ent.DNSRecord, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisioningStepResolver) FileDownload(ctx context.Context, obj *ent.ProvisioningStep) (*ent.FileDownload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisioningStepResolver) FileDelete(ctx context.Context, obj *ent.ProvisioningStep) (*ent.FileDelete, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisioningStepResolver) FileExtract(ctx context.Context, obj *ent.ProvisioningStep) (*ent.FileExtract, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Environments(ctx context.Context) ([]*ent.Environment, error) {
	e, err := r.client.Environment.Query().All(ctx)
	if err != nil {
        return nil, fmt.Errorf("failed querying Environment: %v", err)
    }
    log.Println("Environment returned: ", e)
    return e, nil
}

func (r *queryResolver) Environment(ctx context.Context, envUUID string) (*ent.Environment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProvisionedHost(ctx context.Context, proHostUUID string) (*ent.ProvisionedHost, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProvisionedNetwork(ctx context.Context, proNetUUID string) (*ent.ProvisionedNetwork, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ProvisionedStep(ctx context.Context, proStepUUID string) (*ent.ProvisioningStep, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *scriptResolver) Vars(ctx context.Context, obj *ent.Script) ([]*model.VarsMap, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *scriptResolver) Tags(ctx context.Context, obj *ent.Script) ([]*ent.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *scriptResolver) Maintainer(ctx context.Context, obj *ent.Script) (*ent.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *scriptResolver) Findings(ctx context.Context, obj *ent.Script) ([]*ent.Finding, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *statusResolver) State(ctx context.Context, obj *ent.Status) (model.ProvisionStatus, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *statusResolver) StartedAt(ctx context.Context, obj *ent.Status) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *statusResolver) EndedAt(ctx context.Context, obj *ent.Status) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *tagResolver) Description(ctx context.Context, obj *ent.Tag) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *teamResolver) Config(ctx context.Context, obj *ent.Team) ([]*model.ConfigMap, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *teamResolver) Maintainer(ctx context.Context, obj *ent.Team) (*ent.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *teamResolver) Build(ctx context.Context, obj *ent.Team) (*ent.Build, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *teamResolver) Environment(ctx context.Context, obj *ent.Team) (*ent.Environment, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *teamResolver) Tags(ctx context.Context, obj *ent.Team) ([]*ent.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

// Build returns generated.BuildResolver implementation.
func (r *Resolver) Build() generated.BuildResolver { return &buildResolver{r} }

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

// FileDownload returns generated.FileDownloadResolver implementation.
func (r *Resolver) FileDownload() generated.FileDownloadResolver { return &fileDownloadResolver{r} }

// Finding returns generated.FindingResolver implementation.
func (r *Resolver) Finding() generated.FindingResolver { return &findingResolver{r} }

// Host returns generated.HostResolver implementation.
func (r *Resolver) Host() generated.HostResolver { return &hostResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Network returns generated.NetworkResolver implementation.
func (r *Resolver) Network() generated.NetworkResolver { return &networkResolver{r} }

// ProvisionedHost returns generated.ProvisionedHostResolver implementation.
func (r *Resolver) ProvisionedHost() generated.ProvisionedHostResolver {
	return &provisionedHostResolver{r}
}

// ProvisionedNetwork returns generated.ProvisionedNetworkResolver implementation.
func (r *Resolver) ProvisionedNetwork() generated.ProvisionedNetworkResolver {
	return &provisionedNetworkResolver{r}
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

// Tag returns generated.TagResolver implementation.
func (r *Resolver) Tag() generated.TagResolver { return &tagResolver{r} }

// Team returns generated.TeamResolver implementation.
func (r *Resolver) Team() generated.TeamResolver { return &teamResolver{r} }

type buildResolver struct{ *Resolver }
type commandResolver struct{ *Resolver }
type competitionResolver struct{ *Resolver }
type dNSResolver struct{ *Resolver }
type dNSRecordResolver struct{ *Resolver }
type environmentResolver struct{ *Resolver }
type fileDownloadResolver struct{ *Resolver }
type findingResolver struct{ *Resolver }
type hostResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type networkResolver struct{ *Resolver }
type provisionedHostResolver struct{ *Resolver }
type provisionedNetworkResolver struct{ *Resolver }
type provisioningStepResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type scriptResolver struct{ *Resolver }
type statusResolver struct{ *Resolver }
type tagResolver struct{ *Resolver }
type teamResolver struct{ *Resolver }
