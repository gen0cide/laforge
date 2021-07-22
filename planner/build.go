package planner

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gen0cide/laforge/builder"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/agenttask"
	"github.com/gen0cide/laforge/ent/plan"
	"github.com/gen0cide/laforge/ent/provisioningstep"
	"github.com/gen0cide/laforge/ent/servertask"
	"github.com/gen0cide/laforge/ent/status"
	"github.com/gen0cide/laforge/server/utils"
	"github.com/sirupsen/logrus"
)

func StartBuild(client *ent.Client, currentUser *ent.AuthUser, entBuild *ent.Build) error {
	ctx := context.Background()
	defer ctx.Done()

	entEnvironment, err := entBuild.QueryBuildToEnvironment().Only(ctx)
	if err != nil {
		logrus.Errorf("failed to query environment from build: %v", err)
		return err
	}

	taskStatus, serverTask, err := utils.CreateServerTask(ctx, client, rdb, currentUser, servertask.TypeEXECUTEBUILD)
	if err != nil {
		return fmt.Errorf("error creating server task: %v", err)
	}
	serverTask, err = client.ServerTask.UpdateOne(serverTask).SetServerTaskToBuild(entBuild).SetServerTaskToEnvironment(entEnvironment).Save(ctx)
	if err != nil {
		taskStatus, serverTask, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask)
		if err != nil {
			return fmt.Errorf("error failing execute build server task: %v", err)
		}
		return fmt.Errorf("error assigning environment and build to execute build server task: %v", err)
	}
	rdb.Publish(ctx, "updatedServerTask", serverTask.ID.String())

	entPlans, err := entBuild.QueryBuildToPlan().Where(plan.HasPlanToStatusWith(status.StateEQ(status.StatePLANNING))).All(ctx)

	if err != nil {
		taskStatus, serverTask, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask)
		if err != nil {
			return fmt.Errorf("error failing execute build server task: %v", err)
		}
		logrus.Errorf("Failed to Query Plan Nodes %v. Err: %v", entPlans, err)
		return err
	}

	var wg sync.WaitGroup

	for _, entPlan := range entPlans {
		entStatus, err := entPlan.PlanToStatus(ctx)

		if err != nil {
			logrus.Errorf("Failed to Query Status %v. Err: %v", entPlan, err)
			return err
		}

		wg.Add(1)

		go func(wg *sync.WaitGroup, entStatus *ent.Status) {
			defer wg.Done()
			ctx := context.Background()
			defer ctx.Done()
			entStatus.Update().SetState(status.StateAWAITING).Save(ctx)
			rdb.Publish(ctx, "updatedStatus", entStatus.ID.String())
		}(&wg, entStatus)

		wg.Add(1)
		go func(wg *sync.WaitGroup, entPlan *ent.Plan) {
			defer wg.Done()
			ctx := context.Background()
			defer ctx.Done()
			switch entPlan.Type {
			case plan.TypeProvisionNetwork:
				entProNetwork, err := entPlan.QueryPlanToProvisionedNetwork().Only(ctx)
				if err != nil {
					logrus.Errorf("Failed to Query Provisioned Network. Err: %v", err)
				}
				entStatus, err := entProNetwork.ProvisionedNetworkToStatus(ctx)
				if err != nil {
					logrus.Errorf("Failed to Query Status %v. Err: %v", entPlan, err)
				}
				entStatus.Update().SetState(status.StateAWAITING).Save(ctx)
				rdb.Publish(ctx, "updatedStatus", entStatus.ID.String())
			case plan.TypeProvisionHost:
				entProHost, err := entPlan.QueryPlanToProvisionedHost().Only(ctx)
				if err != nil {
					logrus.Errorf("Failed to Query Provisioned Host. Err: %v", err)
				}
				entStatus, err := entProHost.ProvisionedHostToStatus(ctx)
				if err != nil {
					logrus.Errorf("Failed to Query Status %v. Err: %v", entPlan, err)
				}
				entStatus.Update().SetState(status.StateAWAITING).Save(ctx)
				rdb.Publish(ctx, "updatedStatus", entStatus.ID.String())
			case plan.TypeExecuteStep:
				entProvisioningStep, err := entPlan.QueryPlanToProvisioningStep().Only(ctx)
				if err != nil {
					logrus.Errorf("Failed to Query Provisioning Step. Err: %v", err)
				}
				entStatus, err := entProvisioningStep.ProvisioningStepToStatus(ctx)
				if err != nil {
					logrus.Errorf("Failed to Query Status %v. Err: %v", entPlan, err)
				}
				entStatus.Update().SetState(status.StateAWAITING).Save(ctx)
				rdb.Publish(ctx, "updatedStatus", entStatus.ID.String())
			case plan.TypeStartTeam:
				entTeam, err := entPlan.QueryPlanToTeam().Only(ctx)
				if err != nil {
					logrus.Errorf("Failed to Query Provisioning Step. Err: %v", err)
				}
				entStatus, err := entTeam.TeamToStatus(ctx)
				if err != nil {
					logrus.Errorf("Failed to Query Status %v. Err: %v", entPlan, err)
				}
				entStatus.Update().SetState(status.StateAWAITING).Save(ctx)
				rdb.Publish(ctx, "updatedStatus", entStatus.ID.String())
			case plan.TypeStartBuild:
				entBuild, err := entPlan.QueryPlanToBuild().Only(ctx)
				if err != nil {
					logrus.Errorf("Failed to Query Provisioning Step. Err: %v", err)
				}
				entStatus, err := entBuild.BuildToStatus(ctx)
				if err != nil {
					logrus.Errorf("Failed to Query Status %v. Err: %v", entPlan, err)
				}
				entStatus.Update().SetState(status.StateAWAITING).Save(ctx)
				rdb.Publish(ctx, "updatedStatus", entStatus.ID.String())
			default:
				break
			}
		}(&wg, entPlan)
	}

	wg.Wait()

	rootPlans, err := entBuild.QueryBuildToPlan().Where(plan.TypeEQ(plan.TypeStartBuild)).All(ctx)
	if err != nil {
		taskStatus, serverTask, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask)
		if err != nil {
			return fmt.Errorf("error failing execute build server task: %v", err)
		}
		logrus.Errorf("Failed to Query Start Plan Nodes. Err: %v", err)
		return err
	}
	environment, err := entBuild.QueryBuildToEnvironment().Only(ctx)
	if err != nil {
		taskStatus, serverTask, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask)
		if err != nil {
			return fmt.Errorf("error failing execute build server task: %v", err)
		}
		logrus.Errorf("Failed to Query Environment. Err: %v", err)
		return err
	}

	// var genericBuilder builder.Builder

	// switch environment.Builder {
	// case "vsphere-nsxt":
	// 	genericBuilder, err = builder.NewVSphereNSXTBuilder(environment)
	// 	if err != nil {
	// 		logrus.Errorf("Failed to make vSphere NSX-T builder. Err: %v", err)
	// 		return err
	// 	}
	// }

	genericBuilder, err := builder.BuilderFromEnvironment(environment)
	if err != nil {
		taskStatus, serverTask, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask)
		if err != nil {
			return fmt.Errorf("error failing execute build server task: %v", err)
		}
		logrus.Errorf("error generating builder: %v", err)
		return err
	}

	for _, entPlan := range rootPlans {
		wg.Add(1)
		go buildRoutine(client, &genericBuilder, ctx, entPlan, &wg)
	}

	wg.Wait()

	taskStatus, serverTask, err = utils.CompleteServerTask(ctx, client, rdb, taskStatus, serverTask)
	if err != nil {
		return fmt.Errorf("error completing execute build server task: %v", err)
	}

	return nil
}

func buildRoutine(client *ent.Client, builder *builder.Builder, ctx context.Context, entPlan *ent.Plan, wg *sync.WaitGroup) {
	defer wg.Done()

	entStatus, err := entPlan.PlanToStatus(ctx)

	if err != nil {
		logrus.Errorf("Failed to Query Status %v. Err: %v", entPlan, err)
	}

	// If it isn't marked for planning, don't worry about traversing to it
	if entStatus.State != status.StateAWAITING {
		return
	}

	prevNodes, err := entPlan.QueryPrevPlan().All(ctx)

	if err != nil {
		logrus.Errorf("Failed to Query Plan Start %v. Err: %v", prevNodes, err)
	}

	parentNodeFailed := false

	for _, prevNode := range prevNodes {
		for {

			if parentNodeFailed {
				break
			}

			prevCompletedStatus, err := prevNode.QueryPlanToStatus().Where(
				status.StateNEQ(
					status.StateCOMPLETE,
				),
			).Exist(ctx)

			if err != nil {
				logrus.Errorf("Failed to Query Status %v. Err: %v", prevNode, err)
			}

			prevFailedStatus, err := prevNode.QueryPlanToStatus().Where(
				status.StateEQ(
					status.StateFAILED,
				),
			).Exist(ctx)

			if err != nil {
				logrus.Errorf("Failed to Query Status %v. Err: %v", prevNode, err)
			}

			if !prevCompletedStatus {
				break
			}

			if prevFailedStatus {
				parentNodeFailed = true
				break
			}

			time.Sleep(time.Second)
		}
	}
	entStatus, err = entPlan.PlanToStatus(ctx)

	if err != nil {
		logrus.Errorf("Failed to Query Status %v. Err: %v", entPlan, err)
	}

	// If it's already in progress, don't worry about traversing to it
	if entStatus.State == status.StateINPROGRESS || entStatus.State == status.StateCOMPLETE {
		return
	}

	entStatus.Update().SetState(status.StateINPROGRESS).Save(ctx)
	rdb.Publish(ctx, "updatedStatus", entStatus.ID.String())

	var planErr error = nil
	switch entPlan.Type {
	case plan.TypeProvisionNetwork:
		entProNetwork, err := entPlan.QueryPlanToProvisionedNetwork().Only(ctx)
		if err != nil {
			logrus.Errorf("Failed to Query Provisioned Network. Err: %v", err)
		}
		if parentNodeFailed {
			networkStatus, err := entProNetwork.QueryProvisionedNetworkToStatus().Only(ctx)
			if err != nil {
				logrus.Errorf("Error while getting Provisioned Network status: %v", err)
			}
			_, saveErr := networkStatus.Update().SetFailed(true).SetState(status.StateFAILED).Save(ctx)
			if saveErr != nil {
				logrus.Errorf("Error while setting Provisioned Network status to FAILED: %v", saveErr)
			}
			rdb.Publish(ctx, "updatedStatus", networkStatus.ID.String())
			planErr = fmt.Errorf("parent node for Provionded Network has failed")
		} else {
			planErr = buildNetwork(client, builder, ctx, entProNetwork)
		}
	case plan.TypeProvisionHost:
		entProHost, err := entPlan.QueryPlanToProvisionedHost().Only(ctx)
		if err != nil {
			logrus.Errorf("Failed to Query Provisioned Host. Err: %v", err)
		}
		if parentNodeFailed {
			hostStatus, err := entProHost.QueryProvisionedHostToStatus().Only(ctx)
			if err != nil {
				logrus.Errorf("Error while getting Provisioned Network status: %v", err)
			}
			_, saveErr := hostStatus.Update().SetFailed(true).SetState(status.StateFAILED).Save(ctx)
			if saveErr != nil {
				logrus.Errorf("Error while setting Provisioned Network status to FAILED: %v", saveErr)
			}
			rdb.Publish(ctx, "updatedStatus", hostStatus.ID.String())
			planErr = fmt.Errorf("parent node for Provionded Host has failed")
		} else {
			planErr = buildHost(client, builder, ctx, entProHost)
		}
	case plan.TypeExecuteStep:
		entProvisioningStep, err := entPlan.QueryPlanToProvisioningStep().Only(ctx)
		if err != nil {
			logrus.Errorf("Failed to Query Provisioning Step. Err: %v", err)
		}
		if parentNodeFailed {
			stepStatus, err := entProvisioningStep.QueryProvisioningStepToStatus().Only(ctx)
			if err != nil {
				logrus.Errorf("Failed to Query Provisioning Step Status. Err: %v", err)
			}
			_, err = stepStatus.Update().SetFailed(true).SetState(status.StateFAILED).Save(ctx)
			if err != nil {
				logrus.Errorf("error while trying to set ent.ProvisioningStep.Status.State to status.StateFAILED: %v", err)
			}
			rdb.Publish(ctx, "updatedStatus", stepStatus.ID.String())
			planErr = fmt.Errorf("parent node for Provisioning Step has failed")
		} else {
			planErr = execStep(client, ctx, entProvisioningStep)
		}
	case plan.TypeStartTeam:
		entTeam, err := entPlan.QueryPlanToTeam().Only(ctx)
		if err != nil {
			logrus.Errorf("Failed to Query Provisioning Step. Err: %v", err)
		}
		entStatus, err := entTeam.TeamToStatus(ctx)
		if err != nil {
			logrus.Errorf("Failed to Query Status %v. Err: %v", entPlan, err)
		}
		entStatus.Update().SetState(status.StateCOMPLETE).Save(ctx)
		rdb.Publish(ctx, "updatedStatus", entStatus.ID.String())
	case plan.TypeStartBuild:
		entBuild, err := entPlan.QueryPlanToBuild().Only(ctx)
		if err != nil {
			logrus.Errorf("Failed to Query Provisioning Step. Err: %v", err)
		}
		entStatus, err := entBuild.BuildToStatus(ctx)
		if err != nil {
			logrus.Errorf("Failed to Query Status %v. Err: %v", entPlan, err)
		}
		entStatus.Update().SetState(status.StateCOMPLETE).Save(ctx)
		rdb.Publish(ctx, "updatedStatus", entStatus.ID.String())
	default:
		break
	}

	if planErr != nil {
		entStatus.Update().SetState(status.StateFAILED).SetFailed(true).Save(ctx)
		rdb.Publish(ctx, "updatedStatus", entStatus.ID.String())
		logrus.WithFields(logrus.Fields{
			"type":    entPlan.Type,
			"builder": (*builder).ID(),
		}).Errorf("error while executing plan: %v", planErr)
	} else {
		entStatus.Update().SetState(status.StateCOMPLETE).SetCompleted(true).Save(ctx)
		rdb.Publish(ctx, "updatedStatus", entStatus.ID.String())
	}

	nextPlans, err := entPlan.QueryNextPlan().All(ctx)
	for _, nextPlan := range nextPlans {
		wg.Add(1)
		go buildRoutine(client, builder, ctx, nextPlan, wg)
	}

}

func buildHost(client *ent.Client, builder *builder.Builder, ctx context.Context, entProHost *ent.ProvisionedHost) error {
	logrus.Infof("deploying %s", entProHost.SubnetIP)
	hostStatus, err := entProHost.QueryProvisionedHostToStatus().Only(ctx)
	if err != nil {
		logrus.Errorf("Error while getting Provisioned Host status: %v", err)
		return err
	}
	_, saveErr := hostStatus.Update().SetState(status.StateINPROGRESS).Save(ctx)
	if saveErr != nil {
		logrus.Errorf("Error while setting Provisioned Host status to INPROGRESS: %v", saveErr)
		return saveErr
	}
	rdb.Publish(ctx, "updatedStatus", hostStatus.ID.String())
	err = (*builder).DeployHost(ctx, entProHost)
	if err != nil {
		logrus.Errorf("Error while deploying host: %v", err)
		_, saveErr := hostStatus.Update().SetFailed(true).SetState(status.StateFAILED).Save(ctx)
		if saveErr != nil {
			logrus.Errorf("Error while setting Provisioned Host status to FAILED: %v", saveErr)
			return saveErr
		}
		rdb.Publish(ctx, "updatedStatus", hostStatus.ID.String())
		return err
	}
	logrus.Infof("deployed %s successfully", entProHost.SubnetIP)
	_, saveErr = hostStatus.Update().SetCompleted(true).SetState(status.StateCOMPLETE).Save(ctx)
	if saveErr != nil {
		logrus.Errorf("Error while setting Provisioned Host status to COMPLETE: %v", saveErr)
		return saveErr
	}
	rdb.Publish(ctx, "updatedStatus", hostStatus.ID.String())
	return nil
}

func buildNetwork(client *ent.Client, builder *builder.Builder, ctx context.Context, entProNetwork *ent.ProvisionedNetwork) error {
	logrus.Infof("deploying %s", entProNetwork.Name)
	networkStatus, err := entProNetwork.QueryProvisionedNetworkToStatus().Only(ctx)
	if err != nil {
		logrus.Errorf("Error while getting Provisioned Network status: %v", err)
		return err
	}
	_, saveErr := networkStatus.Update().SetState(status.StateINPROGRESS).Save(ctx)
	if saveErr != nil {
		logrus.Errorf("Error while setting Provisioned Network status to INPROGRESS: %v", saveErr)
		return saveErr
	}
	rdb.Publish(ctx, "updatedStatus", networkStatus.ID.String())
	err = (*builder).DeployNetwork(ctx, entProNetwork)
	if err != nil {
		logrus.Errorf("Error while deploying network: %v", err)
		_, saveErr := networkStatus.Update().SetFailed(true).SetState(status.StateFAILED).Save(ctx)
		if saveErr != nil {
			logrus.Errorf("Error while setting Provisioned Network status to FAILED: %v", saveErr)
			return saveErr
		}
		rdb.Publish(ctx, "updatedStatus", networkStatus.ID.String())
		return err
	}
	logrus.Infof("deployed %s successfully", entProNetwork.Name)
	// Allow networks to set up
	time.Sleep(1 * time.Minute)

	_, saveErr = networkStatus.Update().SetCompleted(true).SetState(status.StateCOMPLETE).Save(ctx)
	if saveErr != nil {
		logrus.Errorf("Error while setting Provisioned Network status to COMPLETE: %v", saveErr)
		return saveErr
	}
	rdb.Publish(ctx, "updatedStatus", networkStatus.ID.String())
	return nil
}

func execStep(client *ent.Client, ctx context.Context, entStep *ent.ProvisioningStep) error {
	stepStatus, err := entStep.QueryProvisioningStepToStatus().Only(ctx)
	if err != nil {
		logrus.Errorf("Failed to Query Provisioning Step Status. Err: %v", err)
	}
	_, err = stepStatus.Update().SetState(status.StateINPROGRESS).Save(ctx)
	if err != nil {
		logrus.Errorf("error while trying to set ent.ProvisioningStep.Status.State to status.StateCOMPLETED: %v", err)
	}
	rdb.Publish(ctx, "updatedStatus", stepStatus.ID.String())
	GQLHostName, ok := os.LookupEnv("GRAPHQL_HOSTNAME")
	downloadURL := ""

	if !ok {
		downloadURL = "http://localhost:8080/api/download/"
	} else {
		downloadURL = "http://" + GQLHostName + "/api/download/"
	}

	entPorovisionedHost, err := entStep.QueryProvisioningStepToProvisionedHost().Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying Provisioned Host for Provioning Step: %v", err)
	}

	taskCount, err := entPorovisionedHost.QueryProvisionedHostToAgentTask().Count(ctx)
	if err != nil {
		return fmt.Errorf("failed querying Number of Tasks: %v", err)
	}

	switch entStep.Type {
	case provisioningstep.TypeScript:
		entScript, err := entStep.QueryProvisioningStepToScript().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying Script for Provioning Step: %v", err)
		}
		entGinMiddleware, err := entStep.QueryProvisioningStepToGinFileMiddleware().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying Gin File Middleware for Script: %v", err)
		}
		_, err = client.AgentTask.Create().
			SetCommand(agenttask.CommandDOWNLOAD).
			SetArgs(entScript.Source + "💔" + downloadURL + entGinMiddleware.URLID).
			SetNumber(taskCount).
			SetState(agenttask.StateAWAITING).
			SetAgentTaskToProvisionedHost(entPorovisionedHost).
			SetAgentTaskToProvisioningStep(entStep).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed Creating Agent Task for Script Download: %v", err)
		}
		// TODO: Add the Ability to change permissions of a file into the agent
		_, err = client.AgentTask.Create().
			SetCommand(agenttask.CommandEXECUTE).
			SetArgs(entScript.Source + "💔" + strings.Join(entScript.Args, " ")).
			SetNumber(taskCount + 1).
			SetState(agenttask.StateAWAITING).
			SetAgentTaskToProvisionedHost(entPorovisionedHost).
			SetAgentTaskToProvisioningStep(entStep).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed Creating Agent Task for Script Execute: %v", err)
		}
		_, err = client.AgentTask.Create().
			SetCommand(agenttask.CommandDELETE).
			SetArgs(entScript.Source).
			SetNumber(taskCount + 2).
			SetState(agenttask.StateAWAITING).
			SetAgentTaskToProvisionedHost(entPorovisionedHost).
			SetAgentTaskToProvisioningStep(entStep).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed Creating Agent Task for Script Delete: %v", err)
		}
	case provisioningstep.TypeCommand:
		entCommand, err := entStep.QueryProvisioningStepToCommand().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying Command for Provioning Step: %v", err)
		}
		// Check if reboot command
		if entCommand.Program == "REBOOT" {
			_, err = client.AgentTask.Create().
				SetCommand(agenttask.CommandREBOOT).
				SetArgs("").
				SetNumber(taskCount).
				SetState(agenttask.StateAWAITING).
				SetAgentTaskToProvisionedHost(entPorovisionedHost).
				SetAgentTaskToProvisioningStep(entStep).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed Creating Agent Task for Reboot Command: %v", err)
			}
		} else {
			_, err = client.AgentTask.Create().
				SetCommand(agenttask.CommandEXECUTE).
				SetArgs(entCommand.Program + "💔" + strings.Join(entCommand.Args, " ")).
				SetNumber(taskCount).
				SetState(agenttask.StateAWAITING).
				SetAgentTaskToProvisionedHost(entPorovisionedHost).
				SetAgentTaskToProvisioningStep(entStep).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed Creating Agent Task for Command: %v", err)
			}
		}
	case provisioningstep.TypeFileDelete:
		entFileDelete, err := entStep.QueryProvisioningStepToFileDelete().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying File Delete for Provioning Step: %v", err)
		}
		_, err = client.AgentTask.Create().
			SetCommand(agenttask.CommandDELETE).
			SetArgs(entFileDelete.Path).
			SetNumber(taskCount).
			SetState(agenttask.StateAWAITING).
			SetAgentTaskToProvisionedHost(entPorovisionedHost).
			SetAgentTaskToProvisioningStep(entStep).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed Creating Agent Task for File Delete: %v", err)
		}
	case provisioningstep.TypeFileDownload:
		entFileDownload, err := entStep.QueryProvisioningStepToFileDownload().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying File Download for Provioning Step: %v", err)
		}
		entGinMiddleware, err := entStep.QueryProvisioningStepToGinFileMiddleware().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying Gin File Middleware for File Download: %v", err)
		}
		_, err = client.AgentTask.Create().
			SetCommand(agenttask.CommandDOWNLOAD).
			SetArgs(entFileDownload.Destination + "💔" + downloadURL + entGinMiddleware.URLID).
			SetNumber(taskCount).
			SetState(agenttask.StateAWAITING).
			SetAgentTaskToProvisionedHost(entPorovisionedHost).
			SetAgentTaskToProvisioningStep(entStep).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed Creating Agent Task for File Download: %v", err)
		}
	case provisioningstep.TypeFileExtract:
		entFileExtract, err := entStep.QueryProvisioningStepToFileExtract().Only(ctx)
		if err != nil {
			return fmt.Errorf("failed querying File Extract for Provioning Step: %v", err)
		}
		_, err = client.AgentTask.Create().
			SetCommand(agenttask.CommandEXTRACT).
			SetArgs(entFileExtract.Source + "💔" + entFileExtract.Destination).
			SetNumber(taskCount).
			SetState(agenttask.StateAWAITING).
			SetAgentTaskToProvisionedHost(entPorovisionedHost).
			SetAgentTaskToProvisioningStep(entStep).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed Creating Agent Task for File Extract: %v", err)
		}
	case provisioningstep.TypeDNSRecord:
		break
	default:
		break
	}

	for {
		taskFailed, err := entStep.QueryProvisioningStepToAgentTask().Where(
			agenttask.StateEQ(
				agenttask.StateFAILED,
			),
		).Exist(ctx)

		if err != nil {
			logrus.Errorf("Failed to Query Agent Task State. Err: %v", err)
		}

		if taskFailed {
			_, err = stepStatus.Update().SetFailed(true).SetState(status.StateFAILED).Save(ctx)
			if err != nil {
				logrus.Errorf("error while trying to set ent.ProvisioningStep.Status.State to status.StateFAILED: %v", err)
			}
			rdb.Publish(ctx, "updatedStatus", stepStatus.ID.String())
			return fmt.Errorf("one or more agent tasks failed")
		}

		taskRunning, err := entStep.QueryProvisioningStepToAgentTask().Where(
			agenttask.StateNEQ(
				agenttask.StateCOMPLETE,
			),
		).Exist(ctx)

		if err != nil {
			logrus.Errorf("Failed to Query Agent Task State. Err: %v", err)
		}

		if !taskRunning {
			break
		}

		time.Sleep(time.Second)
	}
	_, err = stepStatus.Update().SetCompleted(true).SetState(status.StateCOMPLETE).Save(ctx)
	if err != nil {
		logrus.Errorf("error while trying to set ent.ProvisioningStep.Status.State to status.StateCOMPLETED: %v", err)
	}
	rdb.Publish(ctx, "updatedStatus", stepStatus.ID.String())

	return nil
}
