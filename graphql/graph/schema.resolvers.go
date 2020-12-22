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
	"github.com/gen0cide/laforge/ent/host"
	"github.com/gen0cide/laforge/ent/provisionedhost"
	"github.com/gen0cide/laforge/ent/provisionednetwork"
	"github.com/gen0cide/laforge/ent/provisioningstep"
	"github.com/gen0cide/laforge/graphql/graph/generated"
	"github.com/gen0cide/laforge/graphql/graph/model"
)

func (r *buildResolver) Tags(ctx context.Context, obj *ent.Build) ([]*ent.Tag, error) {
	t, err := obj.QueryTag().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Tags: %v", err)
	}

	return t, nil
}

func (r *buildResolver) Config(ctx context.Context, obj *ent.Build) ([]*model.ConfigMap, error) {
	results := make([]*model.ConfigMap, 0)

	for k, v := range obj.Config {
		results = append(results, &model.ConfigMap{k, v})
	}

	return results, nil
}

func (r *buildResolver) Maintainer(ctx context.Context, obj *ent.Build) (*ent.User, error) {
	u, err := obj.QueryMaintainer().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Maintainer: %v", err)
	}

	return u, nil
}

func (r *buildResolver) Teams(ctx context.Context, obj *ent.Build) ([]*ent.Team, error) {
	t, err := obj.QueryTeam().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Teams: %v", err)
	}

	return t, nil
}

func (r *commandResolver) Vars(ctx context.Context, obj *ent.Command) ([]*model.VarsMap, error) {
	results := make([]*model.VarsMap, 0)

	for k, v := range obj.Vars {
		results = append(results, &model.VarsMap{k, v})
	}

	return results, nil
}

func (r *commandResolver) Tags(ctx context.Context, obj *ent.Command) ([]*ent.Tag, error) {
	t, err := obj.QueryTag().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("error querying Tags: %v", err)
	}

	return t, nil
}

func (r *commandResolver) Maintainer(ctx context.Context, obj *ent.Command) (*ent.User, error) {
	u, err := obj.QueryUser().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("error querying User: %v", err)
	}

	return u, nil
}

func (r *competitionResolver) Config(ctx context.Context, obj *ent.Competition) ([]*model.ConfigMap, error) {
	results := make([]*model.ConfigMap, 0)

	for k, v := range obj.Config {
		results = append(results, &model.ConfigMap{k, v})
	}

	return results, nil
}

func (r *competitionResolver) DNS(ctx context.Context, obj *ent.Competition) (*ent.DNS, error) {
	d, err := obj.QueryDNS().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("error querying DNS: %v", err)
	}

	return d, nil
}

func (r *dNSResolver) NTPServer(ctx context.Context, obj *ent.DNS) ([]*string, error) {
	results := make([]*string, 0)

	for _, elem := range obj.DNSServers {
		results = append(results, &elem)
	}

	return results, nil
}

func (r *dNSResolver) Config(ctx context.Context, obj *ent.DNS) ([]*model.ConfigMap, error) {
	results := make([]*model.ConfigMap, 0)

	for k, v := range obj.Config {
		results = append(results, &model.ConfigMap{k, v})
	}

	return results, nil
}

func (r *dNSRecordResolver) Vars(ctx context.Context, obj *ent.DNSRecord) ([]*model.VarsMap, error) {
	results := make([]*model.VarsMap, 0)

	for k, v := range obj.Vars {
		results = append(results, &model.VarsMap{k, v})
	}

	return results, nil
}

func (r *dNSRecordResolver) Tags(ctx context.Context, obj *ent.DNSRecord) ([]*ent.Tag, error) {
	t, err := obj.QueryTag().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Tags: %v", err)
	}

	return t, nil
}

func (r *environmentResolver) Tags(ctx context.Context, obj *ent.Environment) ([]*ent.Tag, error) {
	t, err := obj.QueryTag().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Tags: %v", err)
	}

	return t, nil
}

func (r *environmentResolver) Config(ctx context.Context, obj *ent.Environment) ([]*model.ConfigMap, error) {
	results := make([]*model.ConfigMap, 0)

	for k, v := range obj.Config {
		results = append(results, &model.ConfigMap{k, v})
	}

	return results, nil
}

func (r *environmentResolver) Maintainer(ctx context.Context, obj *ent.Environment) (*ent.User, error) {
	u, err := obj.QueryUser().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying User: %v", err)
	}

	return u, nil
}

func (r *environmentResolver) Networks(ctx context.Context, obj *ent.Environment) ([]*ent.Network, error) {
	n, err := obj.QueryNetwork().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Network: %v", err)
	}

	return n, nil
}

func (r *environmentResolver) Hosts(ctx context.Context, obj *ent.Environment) ([]*ent.Host, error) {
	h, err := obj.QueryHost().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Host: %v", err)
	}

	return h, nil
}

func (r *environmentResolver) Build(ctx context.Context, obj *ent.Environment) (*ent.Build, error) {
	b, err := obj.QueryBuild().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Build: %v", err)
	}

	return b, nil
}

func (r *environmentResolver) Competition(ctx context.Context, obj *ent.Environment) (*ent.Competition, error) {
	c, err := obj.QueryCompetition().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Competition: %v", err)
	}

	return c, nil
}

func (r *fileDownloadResolver) Templete(ctx context.Context, obj *ent.FileDownload) (bool, error) {
	return obj.Template, nil
}

func (r *fileDownloadResolver) Tags(ctx context.Context, obj *ent.FileDownload) ([]*ent.Tag, error) {
	t, err := obj.QueryTag().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Tags: %v", err)
	}

	return t, nil
}

func (r *findingResolver) Severity(ctx context.Context, obj *ent.Finding) (model.FindingSeverity, error) {
	return model.FindingSeverity(obj.Severity), nil
}

func (r *findingResolver) Difficulty(ctx context.Context, obj *ent.Finding) (model.FindingDifficulty, error) {
	return model.FindingDifficulty(obj.Difficulty), nil
}

func (r *findingResolver) Maintainer(ctx context.Context, obj *ent.Finding) (*ent.User, error) {
	u, err := obj.QueryUser().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying User: %v", err)
	}

	return u, nil
}

func (r *findingResolver) Tags(ctx context.Context, obj *ent.Finding) ([]*ent.Tag, error) {
	t, err := obj.QueryTag().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Tags: %v", err)
	}

	return t, nil
}

func (r *findingResolver) Host(ctx context.Context, obj *ent.Finding) (*ent.Host, error) {
	h, err := obj.QueryHost().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Host: %v", err)
	}

	return h, nil
}

func (r *hostResolver) Vars(ctx context.Context, obj *ent.Host) ([]*model.VarsMap, error) {
	results := make([]*model.VarsMap, 0)

	for k, v := range obj.Vars {
		results = append(results, &model.VarsMap{k, v})
	}

	return results, nil
}

func (r *hostResolver) DependsOn(ctx context.Context, obj *ent.Host) ([]*ent.Host, error) {
	results := make([]*ent.Host, 0)

	for dependsOn := range obj.DependsOn {
		h, err := r.client.Host.Query().Where(host.IDEQ(dependsOn)).Only(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed querying Host: %v", err)
		}

		results = append(results, h)
	}

	return results, nil
}

func (r *hostResolver) Maintainer(ctx context.Context, obj *ent.Host) (*ent.User, error) {
	u, err := obj.QueryMaintainer().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying User: %v", err)
	}

	return u, nil
}

func (r *hostResolver) Tags(ctx context.Context, obj *ent.Host) ([]*ent.Tag, error) {
	t, err := obj.QueryTag().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Tags: %v", err)
	}

	return t, nil
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

	for k, v := range obj.Vars {
		results = append(results, &model.VarsMap{k, v})
	}

	return results, nil
}

func (r *networkResolver) Tags(ctx context.Context, obj *ent.Network) ([]*ent.Tag, error) {
	t, err := obj.QueryTag().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Tags: %v", err)
	}

	return t, nil
}

func (r *provisionedHostResolver) Status(ctx context.Context, obj *ent.ProvisionedHost) (*ent.Status, error) {
	s, err := obj.QueryStatus().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Status: %v", err)
	}

	return s, nil
}

func (r *provisionedHostResolver) ProvisionedNetwork(ctx context.Context, obj *ent.ProvisionedHost) (*ent.ProvisionedNetwork, error) {
	pn, err := obj.QueryProvisionedNetwork().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying ProvisionedNetwork: %v", err)
	}

	return pn, nil
}

func (r *provisionedHostResolver) Host(ctx context.Context, obj *ent.ProvisionedHost) (*ent.Host, error) {
	h, err := obj.QueryHost().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Host: %v", err)
	}

	return h, nil
}

func (r *provisionedHostResolver) CombinedOutput(ctx context.Context, obj *ent.ProvisionedHost) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *provisionedHostResolver) Heartbeat(ctx context.Context, obj *ent.ProvisionedHost) (*ent.AgentStatus, error) {
	a, err := obj.QueryAgentStatus().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Agent Status for Host: %v", err)
	}

	return a, nil
}

func (r *provisionedNetworkResolver) Status(ctx context.Context, obj *ent.ProvisionedNetwork) (*ent.Status, error) {
	s, err := obj.QueryStatus().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Status: %v", err)
	}

	return s, nil
}

func (r *provisionedNetworkResolver) Network(ctx context.Context, obj *ent.ProvisionedNetwork) (*ent.Network, error) {
	n, err := obj.QueryNetwork().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Network: %v", err)
	}

	return n, nil
}

func (r *provisionedNetworkResolver) Build(ctx context.Context, obj *ent.ProvisionedNetwork) (*ent.Build, error) {
	b, err := obj.QueryBuild().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Build: %v", err)
	}

	return b, nil
}

func (r *provisioningStepResolver) ProvisionType(ctx context.Context, obj *ent.ProvisioningStep) (string, error) {
	return obj.ProvisionerType, nil
}

func (r *provisioningStepResolver) ProvisionedHost(ctx context.Context, obj *ent.ProvisioningStep) (*ent.ProvisionedHost, error) {
	ph, err := obj.QueryProvisionedHost().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying ProvisionedHost: %v", err)
	}

	return ph, nil
}

func (r *provisioningStepResolver) Status(ctx context.Context, obj *ent.ProvisioningStep) (*ent.Status, error) {
	s, err := obj.QueryStatus().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Status: %v", err)
	}

	return s, nil
}

func (r *provisioningStepResolver) Script(ctx context.Context, obj *ent.ProvisioningStep) (*ent.Script, error) {
	check, err := obj.QueryScript().Exist(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Script: %v", err)
	}

	if check {
		s, err := obj.QueryScript().Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed querying Script: %v", err)
		}
		return s, nil
	}

	return nil, nil
}

func (r *provisioningStepResolver) Command(ctx context.Context, obj *ent.ProvisioningStep) (*ent.Command, error) {
	check, err := obj.QueryCommand().Exist(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Command: %v", err)
	}

	if check {
		c, err := obj.QueryCommand().Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed querying Command: %v", err)
		}
		return c, nil
	}

	return nil, nil
}

func (r *provisioningStepResolver) DNSRecord(ctx context.Context, obj *ent.ProvisioningStep) (*ent.DNSRecord, error) {
	check, err := obj.QueryDNSRecord().Exist(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying DNSRecord: %v", err)
	}

	if check {
		d, err := obj.QueryDNSRecord().Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed querying DNSRecord: %v", err)
		}
		return d, nil
	}

	return nil, nil
}

func (r *provisioningStepResolver) RemoteFile(ctx context.Context, obj *ent.ProvisioningStep) (*ent.RemoteFile, error) {
	check, err := obj.QueryRemoteFile().Exist(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Remote File: %v", err)
	}

	if check {
		f, err := obj.QueryRemoteFile().Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed querying Remote File: %v", err)
		}
		return f, nil
	}

	return nil, nil
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

func (r *remoteFileResolver) Vars(ctx context.Context, obj *ent.RemoteFile) ([]*model.VarsMap, error) {
	results := make([]*model.VarsMap, 0)

	for k, v := range obj.Vars {
		results = append(results, &model.VarsMap{k, v})
	}

	return results, nil
}

func (r *remoteFileResolver) Templete(ctx context.Context, obj *ent.RemoteFile) (bool, error) {
	return obj.Template, nil
}

func (r *remoteFileResolver) Tags(ctx context.Context, obj *ent.RemoteFile) ([]*ent.Tag, error) {
	t, err := obj.QueryTag().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Tag: %v", err)
	}

	return t, nil
}

func (r *scriptResolver) Vars(ctx context.Context, obj *ent.Script) ([]*model.VarsMap, error) {
	results := make([]*model.VarsMap, 0)

	for k, v := range obj.Vars {
		results = append(results, &model.VarsMap{k, v})
	}

	return results, nil
}

func (r *scriptResolver) Tags(ctx context.Context, obj *ent.Script) ([]*ent.Tag, error) {
	t, err := obj.QueryTag().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Tag: %v", err)
	}

	return t, nil
}

func (r *scriptResolver) Maintainer(ctx context.Context, obj *ent.Script) (*ent.User, error) {
	u, err := obj.QueryMaintainer().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying User: %v", err)
	}

	return u, nil
}

func (r *scriptResolver) Findings(ctx context.Context, obj *ent.Script) ([]*ent.Finding, error) {
	f, err := obj.QueryFinding().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Finding: %v", err)
	}

	return f, nil
}

func (r *statusResolver) State(ctx context.Context, obj *ent.Status) (model.ProvisionStatus, error) {
	switch state := obj.State; state {
	case "AWAITING":
		return model.ProvisionStatusProvStatusAwaiting, nil
	case "INPROGRESS":
		return model.ProvisionStatusProvStatusInProgress, nil
	case "FAILED":
		return model.ProvisionStatusProvStatusFailed, nil
	case "COMPLETE":
		return model.ProvisionStatusProvStatusComplete, nil
	case "TAINTED":
		return model.ProvisionStatusProvStatusTainted, nil
	default:
		return model.ProvisionStatusProvStatusUndefined, nil
	}
}

func (r *statusResolver) StartedAt(ctx context.Context, obj *ent.Status) (string, error) {
	return obj.StartedAt.String(), nil
}

func (r *statusResolver) EndedAt(ctx context.Context, obj *ent.Status) (string, error) {
	return obj.EndedAt.String(), nil
}

func (r *tagResolver) Description(ctx context.Context, obj *ent.Tag) (*string, error) {
	desc := fmt.Sprint(obj.Description)
	return &desc, nil
}

func (r *teamResolver) Config(ctx context.Context, obj *ent.Team) ([]*model.ConfigMap, error) {
	results := make([]*model.ConfigMap, 0)

	for k, v := range obj.Config {
		results = append(results, &model.ConfigMap{k, v})
	}

	return results, nil
}

func (r *teamResolver) Maintainer(ctx context.Context, obj *ent.Team) (*ent.User, error) {
	m, err := obj.QueryMaintainer().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying User: %v", err)
	}

	return m, nil
}

func (r *teamResolver) Build(ctx context.Context, obj *ent.Team) (*ent.Build, error) {
	b, err := obj.QueryBuild().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Build: %v", err)
	}

	return b, nil
}

func (r *teamResolver) Environment(ctx context.Context, obj *ent.Team) (*ent.Environment, error) {
	e, err := obj.QueryTeamToEnvironment().Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Environment: %v", err)
	}

	return e, nil
}

func (r *teamResolver) Tags(ctx context.Context, obj *ent.Team) ([]*ent.Tag, error) {
	t, err := obj.QueryTag().All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Tag: %v", err)
	}

	return t, nil
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

// RemoteFile returns generated.RemoteFileResolver implementation.
func (r *Resolver) RemoteFile() generated.RemoteFileResolver { return &remoteFileResolver{r} }

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
type remoteFileResolver struct{ *Resolver }
type scriptResolver struct{ *Resolver }
type statusResolver struct{ *Resolver }
type tagResolver struct{ *Resolver }
type teamResolver struct{ *Resolver }

