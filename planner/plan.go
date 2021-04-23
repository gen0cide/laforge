package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/build"
	"github.com/gen0cide/laforge/ent/environment"
	"github.com/gen0cide/laforge/ent/host"
	"github.com/gen0cide/laforge/ent/network"
	"github.com/gen0cide/laforge/ent/plan"
	"github.com/gen0cide/laforge/ent/provisionedhost"
	"github.com/gen0cide/laforge/ent/provisionednetwork"
	"github.com/gen0cide/laforge/ent/status"
	"github.com/gen0cide/laforge/ent/team"
	"github.com/gen0cide/laforge/grpc"
	"github.com/gen0cide/laforge/server/utils"
	_ "github.com/mattn/go-sqlite3"
)

var renderFiles = false

func main() {
	var wg sync.WaitGroup

	client := ent.SQLLiteOpen("file:test.sqlite?_loc=auto&cache=shared&_fk=1")
	ctx := context.Background()
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	entEnvironment, err := client.Environment.Query().Where(environment.IDEQ(1)).WithEnvironmentToBuild().Only(ctx)
	if err != nil {
		log.Fatalf("Failed to find Environment %v. Err: %v", 1, err)
	}

	entBuild, _ := createBuild(ctx, client, entEnvironment)
	if err != nil {
		log.Fatalf("Failed to create Build for Enviroment %v. Err: %v", 1, err)
	}
	fmt.Println(entBuild)

	for teamNumber := 0; teamNumber <= entEnvironment.TeamCount; teamNumber++ {
		wg.Add(1)
		go createTeam(ctx, client, entBuild, teamNumber, &wg)
	}
	wg.Wait()
}
func createPlanningStatus(ctx context.Context, client *ent.Client, statusFor status.StatusFor) (*ent.Status, error) {
	entStatus, err := client.Status.Create().SetState(status.StatePLANNING).SetStatusFor(statusFor).Save(ctx)
	if err != nil {
		log.Fatalf("Failed to create Status for %v. Err: %v", statusFor, err)
		return nil, err
	}
	return entStatus, nil
}

func createBuild(ctx context.Context, client *ent.Client, entEnvironment *ent.Environment) (*ent.Build, error) {
	entStatus, err := createPlanningStatus(ctx, client, status.StatusForBuild)
	if err != nil {
		return nil, err
	}
	entBuild, err := client.Build.Create().
		SetRevision(len(entEnvironment.Edges.EnvironmentToBuild)).
		SetBuildToEnvironment(entEnvironment).
		SetBuildToStatus(entStatus).
		Save(ctx)
	if err != nil {
		log.Fatalf("Failed to create Build %v for Enviroment %v. Err: %v", len(entEnvironment.Edges.EnvironmentToBuild), entEnvironment.HclID, err)
		return nil, err
	}
	_, err = client.Plan.Create().
		SetType(plan.TypeStartBuild).
		SetBuildID(entBuild.ID).
		SetPlanToBuild(entBuild).
		SetStepNumber(0).
		Save(ctx)
	if err != nil {
		log.Fatalf("Failed to create Plan Node for Build %v. Err: %v", entBuild.ID, err)
		return nil, err
	}
	return entBuild, nil
}

func createTeam(ctx context.Context, client *ent.Client, entBuild *ent.Build, teamNumber int, wg *sync.WaitGroup) (*ent.Team, error) {
	defer wg.Done()

	entStatus, err := createPlanningStatus(ctx, client, status.StatusForTeam)
	if err != nil {
		return nil, err
	}
	entTeam, err := client.Team.Create().
		SetTeamNumber(teamNumber).
		SetTeamToBuild(entBuild).
		SetTeamToStatus(entStatus).
		Save(ctx)
	if err != nil {
		log.Fatalf("Failed to create Team Number %v for Build %v. Err: %v", teamNumber, entBuild.ID, err)
		return nil, err
	}
	buildPlanNode, err := entBuild.QueryBuildToPlan().Only(ctx)
	if err != nil {
		log.Fatalf("Failed to Query Plan Node for Build %v. Err: %v", entBuild.ID, err)
		return nil, err
	}
	_, err = client.Plan.Create().
		AddPrevPlan(buildPlanNode).
		SetType(plan.TypeStartTeam).
		SetBuildID(entBuild.ID).
		SetPlanToTeam(entTeam).
		SetStepNumber(1).
		Save(ctx)
	if err != nil {
		log.Fatalf("Failed to create Plan Node for Team %v. Err: %v", teamNumber, err)
		return nil, err
	}
	buildNetworks, err := entBuild.QueryBuildToEnvironment().QueryEnvironmentToNetwork().All(ctx)
	if err != nil {
		log.Fatalf("Failed to Query Enviroment for Build %v. Err: %v", entBuild.ID, err)
		return nil, err
	}
	createProvisonedNetworks := []*ent.ProvisionedNetwork{}
	for _, buildNetwork := range buildNetworks {
		pNetwork, _ := createProvisionedNetworks(ctx, client, entBuild, entTeam, buildNetwork)
		createProvisonedNetworks = append(createProvisonedNetworks, pNetwork)
	}
	for _, pNetwork := range createProvisonedNetworks {
		entHosts, err := pNetwork.
			QueryProvisionedNetworkToNetwork().
			QueryNetworkToIncludedNetwork().
			QueryIncludedNetworkToHost().All(ctx)
		if err != nil {
			log.Fatalf("Failed to Query Hosts for Network %v. Err: %v", pNetwork.Name, err)
			return nil, err
		}
		networkPlan, err := pNetwork.QueryProvisionedNetworkToPlan().Only(ctx)
		if err != nil {
			log.Fatalf("Failed to Query Plan for Network %v. Err: %v", pNetwork.Name, err)
			return nil, err
		}
		for _, entHost := range entHosts {
			createProvisionedHosts(ctx, client, pNetwork, entHost, networkPlan)
		}
	}
	return entTeam, nil
}

func createProvisionedNetworks(ctx context.Context, client *ent.Client, entBuild *ent.Build, entTeam *ent.Team, entNetwork *ent.Network) (*ent.ProvisionedNetwork, error) {

	entStatus, err := createPlanningStatus(ctx, client, status.StatusForProvisionedNetwork)
	if err != nil {
		return nil, err
	}

	entProvisionedNetwork, err := client.ProvisionedNetwork.Create().
		SetName(entNetwork.Name).
		SetCidr(entNetwork.Cidr).
		SetProvisionedNetworkToStatus(entStatus).
		SetProvisionedNetworkToNetwork(entNetwork).
		SetProvisionedNetworkToTeam(entTeam).
		SetProvisionedNetworkToBuild(entBuild).
		Save(ctx)
	if err != nil {
		log.Fatalf("Failed to create Provisoned Network %v for Team %v. Err: %v", entNetwork.Name, entTeam.TeamNumber, err)
		return nil, err
	}
	teamPlanNode, err := entTeam.QueryTeamToPlan().Only(ctx)
	if err != nil {
		log.Fatalf("Failed to Query Plan Node for Build %v. Err: %v", entBuild.ID, err)
		return nil, err
	}
	_, err = client.Plan.Create().
		AddPrevPlan(teamPlanNode).
		SetType(plan.TypeProvisionNetwork).
		SetBuildID(entBuild.ID).
		SetPlanToProvisionedNetwork(entProvisionedNetwork).
		SetStepNumber(teamPlanNode.StepNumber + 1).
		Save(ctx)
	if err != nil {
		log.Fatalf("Failed to create Plan Node for Provisioned Network  %v. Err: %v", entProvisionedNetwork.Name, err)
		return nil, err
	}
	return entProvisionedNetwork, nil
}

func createProvisionedHosts(ctx context.Context, client *ent.Client, pNetwork *ent.ProvisionedNetwork, entHost *ent.Host, prevPlan *ent.Plan) (*ent.ProvisionedHost, error) {
	prevPlans := []*ent.Plan{prevPlan}
	planStepNumber := prevPlan.StepNumber + 1
	entProvisionedHost, err := client.ProvisionedHost.Query().Where(
		provisionedhost.And(
			provisionedhost.HasProvisionedHostToProvisionedNetworkWith(
				provisionednetwork.IDEQ(pNetwork.ID),
			),
			provisionedhost.HasProvisionedHostToHostWith(
				host.IDEQ(entHost.ID),
			),
		),
	).Only(ctx)
	if err != nil {
		if err != err.(*ent.NotFoundError) {
			log.Fatalf("Failed to Query Existing Host %v. Err: %v", entHost.HclID, err)
			return nil, err
		}
	} else {
		return entProvisionedHost, nil
	}

	entHostDependencies, err := entHost.QueryDependByHostToHostDependency().
		WithHostDependencyToDependOnHost().
		WithHostDependencyToNetwork().
		All(ctx)

	currentBuild := pNetwork.QueryProvisionedNetworkToBuild().WithBuildToEnvironment().OnlyX(ctx)
	currentTeam := pNetwork.QueryProvisionedNetworkToTeam().OnlyX(ctx)

	for _, entHostDependency := range entHostDependencies {
		entDependsOnHost, err := client.ProvisionedHost.Query().Where(
			provisionedhost.And(
				provisionedhost.HasProvisionedHostToProvisionedNetworkWith(
					provisionednetwork.And(
						provisionednetwork.HasProvisionedNetworkToNetworkWith(
							network.IDEQ(entHostDependency.Edges.HostDependencyToNetwork.ID),
						),
						provisionednetwork.HasProvisionedNetworkToBuildWith(
							build.IDEQ(currentBuild.ID),
						),
						provisionednetwork.HasProvisionedNetworkToTeamWith(
							team.IDEQ(currentTeam.ID),
						),
					),
				),
				provisionedhost.HasProvisionedHostToHostWith(
					host.IDEQ(entHostDependency.Edges.HostDependencyToDependOnHost.ID),
				),
			),
		).WithProvisionedHostToPlan().Only(ctx)
		if err != nil {
			if err != err.(*ent.NotFoundError) {
				log.Fatalf("Failed to Query Depended On Host %v for Host %v. Err: %v", entHostDependency.Edges.HostDependencyToDependOnHost.HclID, entHost.HclID, err)
				return nil, err
			} else {
				dependOnPnetwork, err := client.ProvisionedNetwork.Query().Where(
					provisionednetwork.And(
						provisionednetwork.HasProvisionedNetworkToNetworkWith(
							network.IDEQ(entHostDependency.Edges.HostDependencyToNetwork.ID),
						),
						provisionednetwork.HasProvisionedNetworkToBuildWith(
							build.IDEQ(currentBuild.ID),
						),
						provisionednetwork.HasProvisionedNetworkToTeamWith(
							team.IDEQ(currentTeam.ID),
						),
					),
				).Only(ctx)
				if err != nil {
					log.Fatalf("Failed to Query Provined Network %v for Depended On Host %v. Err: %v", entHostDependency.Edges.HostDependencyToNetwork.HclID, entHostDependency.Edges.HostDependencyToDependOnHost.HclID, err)
				}
				entDependsOnHost, err = createProvisionedHosts(ctx, client, dependOnPnetwork, entHostDependency.Edges.HostDependencyToDependOnHost, prevPlan)
			}
		}
		dependOnPlan, err := entDependsOnHost.QueryProvisionedHostToPlan().Only(ctx)
		if err != nil && err != err.(*ent.NotFoundError) {
			log.Fatalf("Failed to Query Depended On Host %v Plan for Host %v. Err: %v", entHostDependency.Edges.HostDependencyToDependOnHost.HclID, entHost.HclID, err)
			return nil, err
		}
		prevPlans = append(prevPlans, dependOnPlan)
		if planStepNumber <= dependOnPlan.StepNumber {
			planStepNumber = dependOnPlan.StepNumber + 1
		}

	}

	// When get Internet combine CIDR in pNetwork with LastOctet in entHost
	subnetIP := fmt.Sprint(entHost.LastOctet)

	entStatus, err := createPlanningStatus(ctx, client, status.StatusForProvisionedHost)
	if err != nil {
		return nil, err
	}

	entProvisionedHost, err = client.ProvisionedHost.Create().
		SetSubnetIP(subnetIP).
		SetProvisionedHostToStatus(entStatus).
		SetProvisionedHostToProvisionedNetwork(pNetwork).
		SetProvisionedHostToHost(entHost).
		Save(ctx)

	_, err = client.Plan.Create().
		AddPrevPlan(prevPlans...).
		SetType(plan.TypeProvisionHost).
		SetBuildID(prevPlan.BuildID).
		SetPlanToProvisionedHost(entProvisionedHost).
		SetStepNumber(planStepNumber).
		Save(ctx)

	if err != nil {
		log.Fatalf("Failed to create Plan Node for Provisioned Host  %v. Err: %v", entHost.HclID, err)
		return nil, err
	}

	serverAddress, ok := os.LookupEnv("GRPC_SERVER")
	if !ok {
		serverAddress = "localhost:50051"
	}
	isWindowsHost := false
	if strings.Contains(entHost.OS, "w2k") {
		isWindowsHost = true
	}

	binaryPath := path.Join(currentBuild.Edges.BuildToEnvironment.Name, fmt.Sprint(currentBuild.Revision), fmt.Sprint(currentTeam.TeamNumber), pNetwork.Name, entHost.Hostname)
	os.MkdirAll(binaryPath, 0755)
	binaryName := path.Join(binaryPath, "laforgeAgent")
	if isWindowsHost {
		binaryName = binaryName + ".exe"
	}
	binaryName, err = filepath.Abs(binaryName)
	if err != nil {
		log.Fatalf("Unable to Resolve Absolute File Path. Err: %v", err)
		return nil, err
	}
	if renderFiles {
		grpc.BuildAgent(fmt.Sprint(entProvisionedHost.ID), serverAddress, binaryName, isWindowsHost)
		entTmpUrl, err := utils.CreateTempURL(ctx, client, binaryName)
		if err != nil {
			return nil, err
		}
		_, err = entTmpUrl.Update().SetGinFileMiddlewareToProvisionedHost(entProvisionedHost).Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	return entProvisionedHost, nil
}
