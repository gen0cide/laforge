package planner

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/build"
	"github.com/gen0cide/laforge/ent/buildcommit"
	"github.com/gen0cide/laforge/ent/command"
	"github.com/gen0cide/laforge/ent/competition"
	"github.com/gen0cide/laforge/ent/dnsrecord"
	"github.com/gen0cide/laforge/ent/environment"
	"github.com/gen0cide/laforge/ent/filedelete"
	"github.com/gen0cide/laforge/ent/filedownload"
	"github.com/gen0cide/laforge/ent/fileextract"
	"github.com/gen0cide/laforge/ent/host"
	"github.com/gen0cide/laforge/ent/network"
	"github.com/gen0cide/laforge/ent/plan"
	"github.com/gen0cide/laforge/ent/provisionedhost"
	"github.com/gen0cide/laforge/ent/provisionednetwork"
	"github.com/gen0cide/laforge/ent/provisioningstep"
	"github.com/gen0cide/laforge/ent/script"
	"github.com/gen0cide/laforge/ent/servertask"
	"github.com/gen0cide/laforge/ent/status"
	"github.com/gen0cide/laforge/ent/team"
	"github.com/gen0cide/laforge/grpc"
	"github.com/gen0cide/laforge/server/utils"
	"github.com/go-redis/redis/v8"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

var RenderFiles = false
var RenderFilesTask *ent.ServerTask = nil
var RenderFilesTaskStatus *ent.Status = nil

// func main() {

// 	client := ent.SQLLiteOpen("file:test.sqlite?_loc=auto&cache=shared&_fk=1")
// 	ctx := context.Background()
// 	defer client.Close()

// 	// Run the auto migration tool.
// 	if err := client.Schema.Create(ctx); err != nil {
// 		logrus.Errorf("failed creating schema resources: %v", err)
// 	}
// 	uuidString := "36579b83-cc50-4f9f-a007-6da25467dc8a"
// 	envID, err := uuid.Parse(uuidString)
// 	if err != nil {
// 		logrus.Errorf("Unable to parse UUID %v. Err: %v", uuidString, err)
// 	}

// 	entEnvironment, err := client.Environment.Query().Where(environment.ID(envID)).WithEnvironmentToBuild().Only(ctx)
// 	if err != nil {
// 		logrus.Errorf("Failed to find Environment %v. Err: %v", uuidString, err)
// 	}

// 	entBuild, _ := CreateBuild(ctx, client, entEnvironment)
// 	if err != nil {
// 		logrus.Errorf("Failed to create Build for Enviroment %v. Err: %v", 1, err)
// 	}
// 	fmt.Println(entBuild)
// }

func createPlanningStatus(ctx context.Context, client *ent.Client, statusFor status.StatusFor) (*ent.Status, error) {
	entStatus, err := client.Status.Create().SetState(status.StatePLANNING).SetStatusFor(statusFor).Save(ctx)
	if err != nil {
		logrus.Errorf("Failed to create Status for %v. Err: %v", statusFor, err)
		return nil, err
	}
	return entStatus, nil
}

func CreateBuild(ctx context.Context, client *ent.Client, rdb *redis.Client, currentUser *ent.AuthUser, entEnvironment *ent.Environment) (*ent.Build, error) {
	taskStatus, serverTask, err := utils.CreateServerTask(ctx, client, rdb, currentUser, servertask.TypeCREATEBUILD)
	if err != nil {
		return nil, fmt.Errorf("error creating server task: %v", err)
	}
	serverTask, err = client.ServerTask.UpdateOne(serverTask).SetServerTaskToEnvironment(entEnvironment).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error assigning environment to create build server task: %v", err)
	}
	rdb.Publish(ctx, "updatedServerTask", serverTask.ID.String())
	if RenderFiles {
		RenderFilesTask, err = client.ServerTask.UpdateOne(RenderFilesTask).SetServerTaskToEnvironment(entEnvironment).Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("error assigning environment to render files server task: %v", err)
		}
		rdb.Publish(ctx, "updatedServerTask", RenderFilesTask.ID.String())
	}

	var wg sync.WaitGroup
	entStatus, err := createPlanningStatus(ctx, client, status.StatusForBuild)
	if err != nil {
		_, _, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask, err)
		if err != nil {
			return nil, fmt.Errorf("error failing server task: %v", err)
		}
		return nil, err
	}
	entCompetition, err := client.Competition.Query().Where(competition.HclIDEQ(entEnvironment.CompetitionID)).Only(ctx)
	if err != nil {
		logrus.Errorf("Failed to Query Competition %v for Enviroment %v. Err: %v", len(entEnvironment.CompetitionID), entEnvironment.HclID, err)
		_, _, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask, err)
		if err != nil {
			return nil, fmt.Errorf("error failing server task: %v", err)
		}
		return nil, err
	}
	entBuild, err := client.Build.Create().
		SetRevision(len(entEnvironment.Edges.EnvironmentToBuild)).
		SetEnvironmentRevision(entEnvironment.Revision).
		SetBuildToEnvironment(entEnvironment).
		SetBuildToStatus(entStatus).
		SetBuildToCompetition(entCompetition).
		Save(ctx)
	if err != nil {
		logrus.Errorf("Failed to create Build %v for Enviroment %v. Err: %v", len(entEnvironment.Edges.EnvironmentToBuild), entEnvironment.HclID, err)
		_, _, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask, err)
		if err != nil {
			return nil, fmt.Errorf("error failing server task: %v", err)
		}
		return nil, err
	}
	serverTask, err = client.ServerTask.UpdateOne(serverTask).SetServerTaskToBuild(entBuild).Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error assigning environment to create build server task: %v", err)
	}
	rdb.Publish(ctx, "updatedServerTask", serverTask.ID.String())
	if RenderFiles {
		RenderFilesTask, err = client.ServerTask.UpdateOne(RenderFilesTask).SetServerTaskToBuild(entBuild).Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("error linking build to render files server task: %v", err)
		}
		rdb.Publish(ctx, "updatedServerTask", RenderFilesTask.ID.String())
	}
	entPlanStatus, err := createPlanningStatus(ctx, client, status.StatusForPlan)
	if err != nil {
		_, _, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask, err)
		if err != nil {
			return nil, fmt.Errorf("error failing server task: %v", err)
		}
		return nil, err
	}

	_, err = client.Plan.Create().
		SetType(plan.TypeStartBuild).
		SetBuildID(entBuild.ID.String()).
		SetPlanToBuild(entBuild).
		SetStepNumber(0).
		SetPlanToStatus(entPlanStatus).
		Save(ctx)
	if err != nil {
		logrus.Errorf("Failed to create Plan Node for Build %v. Err: %v", entBuild.ID, err)
		_, _, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask, err)
		if err != nil {
			return nil, fmt.Errorf("error failing server task: %v", err)
		}
		return nil, err
	}
	for teamNumber := 0; teamNumber < entEnvironment.TeamCount; teamNumber++ {
		wg.Add(1)
		go createTeam(client, entBuild, teamNumber, &wg)
	}

	go func(wg *sync.WaitGroup, entBuild *ent.Build) {
		wg.Wait()

		ctx := context.Background()
		defer ctx.Done()

		entCommit, err := utils.CreateRootCommit(client, rdb, entBuild)
		if err != nil {
			_, _, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask, err)
			return
		}
		err = entBuild.Update().SetBuildToLatestBuildCommit(entCommit).Exec(ctx)
		if err != nil {
			_, _, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask, err)
			return
		}
		rdb.Publish(ctx, "updatedBuild", entBuild.ID.String())

		if RenderFilesTask != nil {
			RenderFilesTaskStatus, RenderFilesTask, err = utils.CompleteServerTask(ctx, client, rdb, RenderFilesTaskStatus, RenderFilesTask)
			if err != nil {
				_, _, err = utils.FailServerTask(ctx, client, rdb, RenderFilesTaskStatus, RenderFilesTask, err)
				return
			}
		}
		_, serverTask, err = utils.CompleteServerTask(ctx, client, rdb, taskStatus, serverTask)
		if err != nil {
			_, _, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask, err)
			return
		}
		serverTask, err = client.ServerTask.UpdateOne(serverTask).SetServerTaskToBuild(entBuild).Save(ctx)
		if err != nil {
			_, _, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask, err)
			return
		}
		rdb.Publish(ctx, "updatedServerTask", serverTask.ID.String())
		rdb.Publish(ctx, "updatedBuild", entBuild.ID.String())
		// entBuild.Update().SetCompletedPlan(true).SaveX(ctx)

		logrus.Debug("-----\nWAITING FOR COMMIT REVIEW\n-----")
		isApproved, err := utils.WaitForCommitReview(client, entCommit, 20*time.Minute)
		if err != nil {
			logrus.Errorf("error while waiting for root commit to be approved: %v", err)
			entCommit.Update().SetState(buildcommit.StateCANCELLED).Exec(ctx)
			rdb.Publish(ctx, "updatedBuildCommit", entCommit.ID.String())
			return
		}
		if isApproved {
			logrus.Debug("-----\nCOMMIT APPROVED\n-----")
			go StartBuild(client, currentUser, entBuild)
		} else {
			logrus.Debug("-----\nCOMMIT CANCELLED/TIMED OUT\n-----")
			logrus.Errorf("root commit has been cancelled or 20 minute timeout has been reached")
			entCommit.Update().SetState(buildcommit.StateCANCELLED).Exec(ctx)
			rdb.Publish(ctx, "updatedBuildCommit", entCommit.ID.String())
		}
	}(&wg, entBuild)

	return entBuild, nil
}

func createTeam(client *ent.Client, entBuild *ent.Build, teamNumber int, wg *sync.WaitGroup) (*ent.Team, error) {
	defer wg.Done()

	ctx := context.Background()
	defer ctx.Done()

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
		logrus.Errorf("Failed to create Team Number %v for Build %v. Err: %v", teamNumber, entBuild.ID, err)
		return nil, err
	}
	buildPlanNode, err := entBuild.QueryBuildToPlan().Where(plan.StepNumberEQ(0)).Only(ctx)
	if err != nil {
		logrus.Errorf("Failed to Query Plan Node for Build %v. Err: %v", entBuild.ID, err)
		return nil, err
	}
	entPlanStatus, err := createPlanningStatus(ctx, client, status.StatusForPlan)
	if err != nil {
		return nil, err
	}

	_, err = client.Plan.Create().
		AddPrevPlan(buildPlanNode).
		SetType(plan.TypeStartTeam).
		SetBuildID(entBuild.ID.String()).
		SetPlanToTeam(entTeam).
		SetPlanToBuild(entBuild).
		SetStepNumber(1).
		SetPlanToStatus(entPlanStatus).
		Save(ctx)
	if err != nil {
		logrus.Errorf("Failed to create Plan Node for Team %v. Err: %v", teamNumber, err)
		return nil, err
	}
	buildNetworks, err := entBuild.QueryBuildToEnvironment().QueryEnvironmentToNetwork().All(ctx)
	if err != nil {
		logrus.Errorf("Failed to Query Enviroment for Build %v. Err: %v", entBuild.ID, err)
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
			QueryIncludedNetworkToHost().
			All(ctx)
		if err != nil {
			logrus.Errorf("Failed to Query Hosts for Network %v. Err: %v", pNetwork.Name, err)
			return nil, err
		}
		networkPlan, err := pNetwork.QueryProvisionedNetworkToPlan().Only(ctx)
		if err != nil {
			logrus.Errorf("Failed to Query Plan for Network %v. Err: %v", pNetwork.Name, err)
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
		logrus.Errorf("Failed to create Provisoned Network %v for Team %v. Err: %v", entNetwork.Name, entTeam.TeamNumber, err)
		return nil, err
	}
	teamPlanNode, err := entTeam.QueryTeamToPlan().Only(ctx)
	if err != nil {
		logrus.Errorf("Failed to Query Plan Node for Build %v. Err: %v", entBuild.ID, err)
		return nil, err
	}

	entPlanStatus, err := createPlanningStatus(ctx, client, status.StatusForPlan)
	if err != nil {
		return nil, err
	}
	_, err = client.Plan.Create().
		AddPrevPlan(teamPlanNode).
		SetType(plan.TypeProvisionNetwork).
		SetBuildID(entBuild.ID.String()).
		SetPlanToProvisionedNetwork(entProvisionedNetwork).
		SetPlanToBuild(entBuild).
		SetStepNumber(teamPlanNode.StepNumber + 1).
		SetPlanToStatus(entPlanStatus).
		Save(ctx)
	if err != nil {
		logrus.Errorf("Failed to create Plan Node for Provisioned Network  %v. Err: %v", entProvisionedNetwork.Name, err)
		return nil, err
	}
	return entProvisionedNetwork, nil
}

func createProvisionedHosts(ctx context.Context, client *ent.Client, pNetwork *ent.ProvisionedNetwork, entHost *ent.Host, prevPlan *ent.Plan) (*ent.ProvisionedHost, error) {
	prevPlans := []*ent.Plan{prevPlan}
	// logrus.Infof("START  %s | %s | %v", pNetwork.Name, entHost.Hostname, prevPlans)
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
			logrus.Errorf("Failed to Query Existing Host %v. Err: %v", entHost.HclID, err)
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
				logrus.Errorf("Failed to Query Depended On Host %v for Host %v. Err: %v", entHostDependency.Edges.HostDependencyToDependOnHost.HclID, entHost.HclID, err)
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
					logrus.Errorf("Failed to Query Provined Network %v for Depended On Host %v. Err: %v", entHostDependency.Edges.HostDependencyToNetwork.HclID, entHostDependency.Edges.HostDependencyToDependOnHost.HclID, err)
				}
				dependOnPnetworkPlan, err := dependOnPnetwork.QueryProvisionedNetworkToPlan().Only(ctx)
				if err != nil {
					logrus.Errorf("error while retrieving plan from provisioned network: %v", err)
					return nil, err
				}
				entDependsOnHost, err = createProvisionedHosts(ctx, client, dependOnPnetwork, entHostDependency.Edges.HostDependencyToDependOnHost, dependOnPnetworkPlan)
			}
		}
		dependOnPlan, err := entDependsOnHost.QueryProvisionedHostToEndStepPlan().Only(ctx)
		if err != nil && err != err.(*ent.NotFoundError) {
			logrus.Errorf("Failed to Query Depended On Host %v Plan for Host %v. Err: %v", entHostDependency.Edges.HostDependencyToDependOnHost.HclID, entHost.HclID, err)
			return nil, err
		}
		prevPlans = append(prevPlans, dependOnPlan)
		if planStepNumber <= dependOnPlan.StepNumber {
			planStepNumber = dependOnPlan.StepNumber + 1
		}

	}

	subnetIP, err := CalcIP(pNetwork.Cidr, entHost.LastOctet)
	if err != nil {
		return nil, err
	}

	entStatus, err := createPlanningStatus(ctx, client, status.StatusForProvisionedHost)
	if err != nil {
		return nil, err
	}

	entProvisionedHost, err = client.ProvisionedHost.Create().
		SetSubnetIP(subnetIP).
		SetProvisionedHostToStatus(entStatus).
		SetProvisionedHostToProvisionedNetwork(pNetwork).
		SetProvisionedHostToHost(entHost).
		// SetProvisionedHostToBuild(entBuild).
		Save(ctx)

	if entHost.Tags["root-dns"] == "true" {
		entProvisionedHost, err = entProvisionedHost.Update().SetAddonType(provisionedhost.AddonTypeDNS).Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	entPlanStatus, err := createPlanningStatus(ctx, client, status.StatusForPlan)
	if err != nil {
		return nil, err
	}

	// logrus.Infof("CREATE %s | %s | %v", pNetwork.Name, entHost.Hostname, prevPlans)
	endPlanNode, err := client.Plan.Create().
		AddPrevPlan(prevPlans...).
		SetType(plan.TypeProvisionHost).
		SetBuildID(prevPlan.BuildID).
		SetPlanToProvisionedHost(entProvisionedHost).
		SetStepNumber(planStepNumber).
		SetPlanToBuild(currentBuild).
		SetPlanToStatus(entPlanStatus).
		Save(ctx)

	if err != nil {
		logrus.Errorf("Failed to create Plan Node for Provisioned Host  %v. Err: %v", entHost.HclID, err)
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

	binaryPath := path.Join("builds", currentBuild.Edges.BuildToEnvironment.Name, fmt.Sprint(currentBuild.Revision), fmt.Sprint(currentTeam.TeamNumber), pNetwork.Name, entHost.Hostname)
	os.MkdirAll(binaryPath, 0755)
	binaryName := path.Join(binaryPath, "laforgeAgent")
	if isWindowsHost {
		binaryName = binaryName + ".exe"
	}
	binaryName, err = filepath.Abs(binaryName)
	if err != nil {
		logrus.Errorf("Unable to Resolve Absolute File Path. Err: %v", err)
		return nil, err
	}
	if RenderFiles {
		grpc.BuildAgent(fmt.Sprint(entProvisionedHost.ID), serverAddress, binaryName, isWindowsHost)
		entTmpUrl, err := utils.CreateTempURL(ctx, client, binaryName)
		if err != nil {
			return nil, err
		}
		_, err = entTmpUrl.Update().SetGinFileMiddlewareToProvisionedHost(entProvisionedHost).Save(ctx)
		if err != nil {
			return nil, err
		}
		if RenderFilesTask != nil {
			RenderFilesTask, err = RenderFilesTask.Update().AddServerTaskToGinFileMiddleware(entTmpUrl).Save(ctx)
			if err != nil {
				return nil, err
			}
		}
	}
	userDataScriptID, ok := entHost.Vars["user_data_script_id"]
	if ok {
		userDataScript, err := client.Script.Query().Where(script.HclIDEQ(userDataScriptID)).Only(ctx)
		if err != nil {
			logrus.Errorf("Failed to Query Script %v. Err: %v", userDataScriptID, err)
			return nil, err
		}
		entUserDataStatus, err := client.Status.Create().SetState(status.StateCOMPLETE).SetStatusFor(status.StatusForProvisioningStep).Save(ctx)
		if err != nil {
			logrus.Errorf("Failed to Create Provisioning Step Status for Script %v. Err: %v", userDataScriptID, err)
			return nil, err
		}
		entUserDataProvisioningStep, err := client.ProvisioningStep.Create().
			SetStepNumber(0).
			SetType(provisioningstep.TypeScript).
			SetProvisioningStepToScript(userDataScript).
			SetProvisioningStepToProvisionedHost(entProvisionedHost).
			SetProvisioningStepToStatus(entUserDataStatus).
			Save(ctx)
		if err != nil {
			logrus.Errorf("Failed to Create Provisioning Step for Script %v. Err: %v", userDataScriptID, err)
			return nil, err
		}
		if RenderFiles {
			filePath, err := renderScript(ctx, client, entUserDataProvisioningStep)
			if err != nil {
				return nil, err
			}
			entTmpUrl, err := utils.CreateTempURL(ctx, client, filePath)
			if err != nil {
				return nil, err
			}
			_, err = entTmpUrl.Update().SetGinFileMiddlewareToProvisioningStep(entUserDataProvisioningStep).Save(ctx)
			if err != nil {
				return nil, err
			}
			if RenderFilesTask != nil {
				RenderFilesTask, err = RenderFilesTask.Update().AddServerTaskToGinFileMiddleware(entTmpUrl).Save(ctx)
				if err != nil {
					return nil, err
				}
			}
		}

	}

	for stepNumber, pStep := range entHost.ProvisionSteps {
		stepNumber = stepNumber + 1
		entProvisioningStep, err := createProvisioningStep(ctx, client, pStep, stepNumber, entProvisionedHost, endPlanNode)
		if err != nil {
			return nil, err
		}
		endPlanNode, err = entProvisioningStep.QueryProvisioningStepToPlan().Only(ctx)
		if err != nil {
			return nil, err
		}
	}
	_, err = entProvisionedHost.Update().SetProvisionedHostToEndStepPlan(endPlanNode).Save(ctx)
	if err != nil {
		logrus.Errorf("Unable to Update The End Step. Err: %v", err)
		return nil, err
	}

	return entProvisionedHost, nil
}

func createProvisioningStep(ctx context.Context, client *ent.Client, hclID string, stepNumber int, pHost *ent.ProvisionedHost, prevPlan *ent.Plan) (*ent.ProvisioningStep, error) {
	var entProvisioningStep *ent.ProvisioningStep
	currentEnviroment, err := pHost.QueryProvisionedHostToHost().QueryHostToEnvironment().Only(ctx)
	currentBuild := pHost.QueryProvisionedHostToProvisionedNetwork().QueryProvisionedNetworkToBuild().WithBuildToEnvironment().OnlyX(ctx)
	if err != nil {
		logrus.Errorf("Failed to Query Current Enviroment for Provisoned Host %v. Err: %v", pHost.ID, err)
		return nil, err
	}
	entStatus, err := createPlanningStatus(ctx, client, status.StatusForProvisioningStep)
	if err != nil {
		return nil, err
	}
	entScript, err := client.Script.Query().Where(
		script.And(
			script.HasScriptToEnvironmentWith(
				environment.IDEQ(currentEnviroment.ID),
			),
			script.HclIDEQ(hclID),
		),
	).Only(ctx)
	if err != nil {
		if err != err.(*ent.NotFoundError) {
			logrus.Errorf("Failed to Query Script %v. Err: %v", hclID, err)
			return nil, err
		} else {
			entCommand, err := client.Command.Query().Where(
				command.And(
					command.HasCommandToEnvironmentWith(
						environment.IDEQ(currentEnviroment.ID),
					),
					command.HclIDEQ(hclID),
				)).Only(ctx)
			if err != nil {
				if err != err.(*ent.NotFoundError) {
					logrus.Errorf("Failed to Query Command %v. Err: %v", hclID, err)
					return nil, err
				} else {
					entFileDownload, err := client.FileDownload.Query().Where(
						filedownload.And(
							filedownload.HasFileDownloadToEnvironmentWith(
								environment.IDEQ(currentEnviroment.ID),
							),
							filedownload.HclIDEQ(hclID),
						)).Only(ctx)
					if err != nil {
						if err != err.(*ent.NotFoundError) {
							logrus.Errorf("Failed to Query FileDownload %v. Err: %v", hclID, err)
							return nil, err
						} else {
							entFileExtract, err := client.FileExtract.Query().Where(
								fileextract.And(
									fileextract.HasFileExtractToEnvironmentWith(
										environment.IDEQ(currentEnviroment.ID),
									),
									fileextract.HclIDEQ(hclID),
								)).Only(ctx)
							if err != nil {
								if err != err.(*ent.NotFoundError) {
									logrus.Errorf("Failed to Query FileExtract %v. Err: %v", hclID, err)
									return nil, err
								} else {
									entFileDelete, err := client.FileDelete.Query().Where(
										filedelete.And(
											filedelete.HasFileDeleteToEnvironmentWith(
												environment.IDEQ(currentEnviroment.ID),
											),
											filedelete.HclIDEQ(hclID),
										)).Only(ctx)
									if err != nil {
										if err != err.(*ent.NotFoundError) {
											logrus.Errorf("Failed to Query FileDelete %v. Err: %v", hclID, err)
											return nil, err
										} else {
											entDNSRecord, err := client.DNSRecord.Query().Where(
												dnsrecord.And(
													dnsrecord.HasDNSRecordToEnvironmentWith(
														environment.IDEQ(currentEnviroment.ID),
													),
													dnsrecord.HclIDEQ(hclID),
												)).Only(ctx)
											if err != nil {
												if err != err.(*ent.NotFoundError) {
													logrus.Errorf("Failed to Query FileDelete %v. Err: %v", hclID, err)
													return nil, err
												} else {
													logrus.Errorf("No Provisioning Steps found for %v. Err: %v", hclID, err)
													return nil, err
												}
											} else {
												entProvisioningStep, err = client.ProvisioningStep.Create().
													SetStepNumber(stepNumber).
													SetType(provisioningstep.TypeDNSRecord).SetProvisioningStepToDNSRecord(entDNSRecord).
													SetProvisioningStepToStatus(entStatus).
													SetProvisioningStepToProvisionedHost(pHost).
													Save(ctx)
												if err != nil {
													logrus.Errorf("Failed to Create Provisioning Step for FileDelete %v. Err: %v", hclID, err)
													return nil, err
												}
											}
										}
									} else {
										entProvisioningStep, err = client.ProvisioningStep.Create().
											SetStepNumber(stepNumber).
											SetType(provisioningstep.TypeFileDelete).SetProvisioningStepToFileDelete(entFileDelete).
											SetProvisioningStepToStatus(entStatus).
											SetProvisioningStepToProvisionedHost(pHost).
											Save(ctx)
										if err != nil {
											logrus.Errorf("Failed to Create Provisioning Step for FileDelete %v. Err: %v", hclID, err)
											return nil, err
										}
									}
								}
							} else {
								entProvisioningStep, err = client.ProvisioningStep.Create().
									SetStepNumber(stepNumber).
									SetType(provisioningstep.TypeFileExtract).
									SetProvisioningStepToFileExtract(entFileExtract).
									SetProvisioningStepToStatus(entStatus).
									SetProvisioningStepToProvisionedHost(pHost).
									Save(ctx)
								if err != nil {
									logrus.Errorf("Failed to Create Provisioning Step for FileExtract %v. Err: %v", hclID, err)
									return nil, err
								}
							}
						}
					} else {
						entProvisioningStep, err = client.ProvisioningStep.Create().
							SetStepNumber(stepNumber).
							SetType(provisioningstep.TypeFileDownload).
							SetProvisioningStepToFileDownload(entFileDownload).
							SetProvisioningStepToStatus(entStatus).
							SetProvisioningStepToProvisionedHost(pHost).
							Save(ctx)
						if err != nil {
							logrus.Errorf("Failed to Create Provisioning Step for FileDownload %v. Err: %v", hclID, err)
							return nil, err
						}
						if RenderFiles {
							filePath, err := renderFileDownload(ctx, entProvisioningStep)
							if err != nil {
								return nil, err
							}
							entTmpUrl, err := utils.CreateTempURL(ctx, client, filePath)
							if err != nil {
								return nil, err
							}
							_, err = entTmpUrl.Update().SetGinFileMiddlewareToProvisioningStep(entProvisioningStep).Save(ctx)
							if err != nil {
								return nil, err
							}
							if RenderFilesTask != nil {
								RenderFilesTask, err = RenderFilesTask.Update().AddServerTaskToGinFileMiddleware(entTmpUrl).Save(ctx)
								if err != nil {
									return nil, err
								}
							}
						}
					}
				}
			} else {
				entProvisioningStep, err = client.ProvisioningStep.Create().
					SetStepNumber(stepNumber).
					SetType(provisioningstep.TypeCommand).
					SetProvisioningStepToCommand(entCommand).
					SetProvisioningStepToStatus(entStatus).
					SetProvisioningStepToProvisionedHost(pHost).
					Save(ctx)
				if err != nil {
					logrus.Errorf("Failed to Create Provisioning Step for Command %v. Err: %v", hclID, err)
					return nil, err
				}
			}
		}
	} else {
		entProvisioningStep, err = client.ProvisioningStep.Create().
			SetStepNumber(stepNumber).
			SetType(provisioningstep.TypeScript).
			SetProvisioningStepToScript(entScript).
			SetProvisioningStepToStatus(entStatus).
			SetProvisioningStepToProvisionedHost(pHost).
			Save(ctx)
		if err != nil {
			logrus.Errorf("Failed to Create Provisioning Step for Script %v. Err: %v", hclID, err)
			return nil, err
		}
		if RenderFiles {
			filePath, err := renderScript(ctx, client, entProvisioningStep)
			if err != nil {
				return nil, err
			}
			entTmpUrl, err := utils.CreateTempURL(ctx, client, filePath)
			if err != nil {
				return nil, err
			}
			_, err = entTmpUrl.Update().SetGinFileMiddlewareToProvisioningStep(entProvisioningStep).Save(ctx)
			if err != nil {
				return nil, err
			}
			if RenderFilesTask != nil {
				RenderFilesTask, err = RenderFilesTask.Update().AddServerTaskToGinFileMiddleware(entTmpUrl).Save(ctx)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	entPlanStatus, err := createPlanningStatus(ctx, client, status.StatusForPlan)
	if err != nil {
		return nil, err
	}

	_, err = client.Plan.Create().
		AddPrevPlan(prevPlan).
		SetType(plan.TypeExecuteStep).
		SetBuildID(prevPlan.BuildID).
		SetPlanToProvisioningStep(entProvisioningStep).
		SetStepNumber(prevPlan.StepNumber + 1).
		SetPlanToBuild(currentBuild).
		SetPlanToStatus(entPlanStatus).
		Save(ctx)

	if err != nil {
		logrus.Errorf("Failed to Create Plan Node for Provisioning Step %v. Err: %v", entProvisioningStep.ID, err)
		return nil, err
	}

	return entProvisioningStep, nil

}

func renderScript(ctx context.Context, client *ent.Client, pStep *ent.ProvisioningStep) (string, error) {
	currentProvisionedHost := pStep.QueryProvisioningStepToProvisionedHost().OnlyX(ctx)
	currentScript := pStep.QueryProvisioningStepToScript().OnlyX(ctx)
	currentProvisionedNetwork := currentProvisionedHost.QueryProvisionedHostToProvisionedNetwork().OnlyX(ctx)
	currentTeam := currentProvisionedNetwork.QueryProvisionedNetworkToTeam().OnlyX(ctx)
	currentBuild := currentTeam.QueryTeamToBuild().OnlyX(ctx)
	currentEnvironment := currentBuild.QueryBuildToEnvironment().OnlyX(ctx)
	currentIncludedNetwork := currentEnvironment.QueryEnvironmentToIncludedNetwork().WithIncludedNetworkToHost().WithIncludedNetworkToNetwork().AllX(ctx)
	currentCompetition := currentBuild.QueryBuildToCompetition().OnlyX(ctx)
	currentNetwork := currentProvisionedNetwork.QueryProvisionedNetworkToNetwork().OnlyX(ctx)
	currentHost := currentProvisionedHost.QueryProvisionedHostToHost().OnlyX(ctx)
	currentIdentities := currentEnvironment.QueryEnvironmentToIdentity().AllX(ctx)
	agentScriptFile := currentProvisionedHost.QueryProvisionedHostToGinFileMiddleware().OnlyX(ctx)
	// Need to Make Unique and change how it's loaded in
	currentDNS := currentCompetition.QueryCompetitionToDNS().FirstX(ctx)
	templateData := TempleteContext{
		Build:              currentBuild,
		Competition:        currentCompetition,
		Environment:        currentEnvironment,
		Host:               currentHost,
		DNS:                currentDNS,
		IncludedNetworks:   currentIncludedNetwork,
		Network:            currentNetwork,
		Script:             currentScript,
		Team:               currentTeam,
		Identities:         currentIdentities,
		ProvisionedNetwork: currentProvisionedNetwork,
		ProvisionedHost:    currentProvisionedHost,
		ProvisioningStep:   pStep,
		AgentSlug:          agentScriptFile.URLID,
	}
	t, err := template.New(strings.Replace(currentScript.Source, "./", "", -1)).Funcs(TemplateFuncLib).ParseFiles(currentScript.AbsPath)
	if err != nil {
		logrus.Errorf("Failed to Parse template for script %v. Err: %v", currentScript.Name, err)
		return "", err
	}
	fileRelativePath := path.Join("builds", currentEnvironment.Name, fmt.Sprint(currentBuild.Revision), fmt.Sprint(currentTeam.TeamNumber), currentProvisionedNetwork.Name, currentHost.Hostname)
	os.MkdirAll(fileRelativePath, 0755)
	fileName := filepath.Base(currentScript.Source)
	fileName = path.Join(fileRelativePath, fileName)
	fileName, err = filepath.Abs(fileName)
	if err != nil {
		return "", err
	}
	f, err := os.Create(fileName)
	if err != nil {
		logrus.Errorf("Error Generating Script %v. Err: %v", currentScript.Name, err)
		return "", err
	}
	err = t.Execute(f, templateData)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"scriptName": currentScript.Name,
			"path":       currentScript.AbsPath,
		}).Errorf("error while parsing template for script: %v", err)
	}
	f.Close()
	return fileName, nil
}

func renderFileDownload(ctx context.Context, pStep *ent.ProvisioningStep) (string, error) {
	currentFileDownload := pStep.QueryProvisioningStepToFileDownload().OnlyX(ctx)
	currentProvisionedHost := pStep.QueryProvisioningStepToProvisionedHost().OnlyX(ctx)
	currentProvisionedNetwork := currentProvisionedHost.QueryProvisionedHostToProvisionedNetwork().OnlyX(ctx)
	currentHost := currentProvisionedHost.QueryProvisionedHostToHost().OnlyX(ctx)
	currentTeam := currentProvisionedNetwork.QueryProvisionedNetworkToTeam().OnlyX(ctx)
	currentBuild := currentTeam.QueryTeamToBuild().OnlyX(ctx)
	currentEnvironment := currentBuild.QueryBuildToEnvironment().OnlyX(ctx)

	fileRelativePath := path.Join("builds", currentEnvironment.Name, fmt.Sprint(currentBuild.Revision), fmt.Sprint(currentTeam.TeamNumber), currentProvisionedNetwork.Name, currentHost.Hostname)
	os.MkdirAll(fileRelativePath, 0755)
	fileName := filepath.Base(currentFileDownload.Source)
	fileName = path.Join(fileRelativePath, fileName)
	fileName, err := filepath.Abs(fileName)
	if err != nil {
		return "", err
	}
	destFile, err := os.Create(fileName)
	if err != nil {
		err = fmt.Errorf("error creating file download: %v", err)
		return "", err
	}
	defer destFile.Close()

	// TODO: SOMETHING
	if currentFileDownload.SourceType == "remote" {
		// http.Get(currentFileDownload.Source)
	} else {
		srcFile, err := os.Open(currentFileDownload.AbsPath)
		if err != nil {
			return "", err
		}
		defer srcFile.Close()
		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return "", err
		}
	}
	return fileName, nil
}

// IPv42Int converts net.IP address objects to their uint32 representation
func IPv42Int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

// Int2IPv4 converts uint32s to their net.IP object
func Int2IPv4(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}
