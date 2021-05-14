package planner

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/agenttask"
	"github.com/gen0cide/laforge/ent/plan"
	"github.com/gen0cide/laforge/ent/provisioningstep"
	"github.com/gen0cide/laforge/ent/status"
)

func StartBuild(client *ent.Client, entBuild *ent.Build) error {
	ctx := context.Background()
	defer ctx.Done()

	entPlans, err := entBuild.QueryBuildToPlan().All(ctx)

	if err != nil {
		log.Fatalf("Failed to Query Plan Nodes %v. Err: %v", entPlans, err)
		return err
	}

	var wg sync.WaitGroup

	for _, entPlan := range entPlans {
		status, err := entPlan.PlanToStatus(ctx)

		if err != nil {
			log.Fatalf("Failed to Query Status %v. Err: %v", entPlan, err)
			return err
		}

		wg.Add(1)

		go func(wg *sync.WaitGroup, status *ent.Status) {
			defer wg.Done()
			ctx := context.Background()
			defer ctx.Done()
			status.Update().SetState("AWAITING").Save(ctx)
		}(&wg, status)
	}

	wg.Wait()

	rootPlans, err := entBuild.QueryBuildToPlan().Where(plan.TypeEQ(plan.TypeStartBuild)).All(ctx)
	if err != nil {
		log.Fatalf("Failed to Query Start Plan Nodes. Err: %v", err)
		return err
	}

	for _, entPlan := range rootPlans {
		wg.Add(1)
		go buildRoutine(client, ctx, entPlan, &wg)
	}

	wg.Wait()

	return nil
}

func buildRoutine(client *ent.Client, ctx context.Context, entPlan *ent.Plan, wg *sync.WaitGroup) {
	defer wg.Done()
	prevNodes, err := entPlan.QueryPrevPlan().All(ctx)

	if err != nil {
		log.Fatalf("Failed to Query Plan Start %v. Err: %v", prevNodes, err)
	}

	for _, prevNode := range prevNodes {
		for {
			prevStatus, err := prevNode.QueryPlanToStatus().Where(
				status.StateNEQ(
					status.StateCOMPLETE,
				),
			).Exist(ctx)

			if err != nil {
				log.Fatalf("Failed to Query Status %v. Err: %v", prevNode, err)
			}

			if !prevStatus {
				break
			}

			time.Sleep(time.Second)
		}
	}
	entStatus, err := entPlan.PlanToStatus(ctx)

	if err != nil {
		log.Fatalf("Failed to Query Status %v. Err: %v", entPlan, err)
	}

	entStatus.Update().SetState(status.StateINPROGRESS).Save(ctx)

	switch entPlan.Type {
	case plan.TypeProvisionNetwork:
		entProNetwork, err := entPlan.QueryPlanToProvisionedNetwork().Only(ctx)
		if err != nil {
			log.Fatalf("Failed to Query Provisioned Network. Err: %v", err)
		}
		buildNetwork(client, ctx, entProNetwork)
	case plan.TypeProvisionHost:
		entProHost, err := entPlan.QueryPlanToProvisionedHost().Only(ctx)
		if err != nil {
			log.Fatalf("Failed to Query Provisioned Host. Err: %v", err)
		}
		buildHost(client, ctx, entProHost)
	case plan.TypeExecuteStep:
		entProvisioningStep, err := entPlan.QueryPlanToProvisioningStep().Only(ctx)
		if err != nil {
			log.Fatalf("Failed to Query Provisioning Step. Err: %v", err)
		}
		execStep(client, ctx, entProvisioningStep)
	default:
		break
	}

	entStatus.Update().SetState(status.StateCOMPLETE).Save(ctx)
	entStatus.Update().SetCompleted(true).Save(ctx)

	nextPlans, err := entPlan.QueryNextPlan().All(ctx)
	for _, nextPlan := range nextPlans {
		wg.Add(1)
		go buildRoutine(client, ctx, nextPlan, wg)
	}
}

func buildHost(client *ent.Client, ctx context.Context, entProHost *ent.ProvisionedHost) error {

	entEnvironment, err := entProHost.
		QueryProvisionedHostToProvisionedNetwork().
		QueryProvisionedNetworkToBuild().
		QueryBuildToEnvironment().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying Enviroment for Provioning Network: %v", err)
	}

	switch entEnvironment.Builder {
	default:
		break
	}
	return nil
}

func buildNetwork(client *ent.Client, ctx context.Context, entProNetwork *ent.ProvisionedNetwork) error {

	entEnvironment, err := entProNetwork.
		QueryProvisionedNetworkToBuild().
		QueryBuildToEnvironment().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed querying Enviroment for Provioning Network: %v", err)
	}

	switch entEnvironment.Builder {
	default:
		break
	}
	return nil
}

func execStep(client *ent.Client, ctx context.Context, entStep *ent.ProvisioningStep) error {
	GQLHostName, ok := os.LookupEnv("GRAPQL_HOSTNAME")
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
			SetArgs(entScript.Source + "," + downloadURL + entGinMiddleware.URLID).
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
			SetArgs(entScript.Source + " " + strings.Join(entScript.Args, " ")).
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
		_, err = client.AgentTask.Create().
			SetCommand(agenttask.CommandEXECUTE).
			SetArgs(entCommand.Program + " " + strings.Join(entCommand.Args, " ")).
			SetNumber(taskCount).
			SetState(agenttask.StateAWAITING).
			SetAgentTaskToProvisionedHost(entPorovisionedHost).
			SetAgentTaskToProvisioningStep(entStep).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed Creating Agent Task for Command: %v", err)
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
			SetArgs(entFileDownload.Destination + "," + downloadURL + entGinMiddleware.URLID).
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
			SetArgs(entFileExtract.Source + "," + entFileExtract.Destination).
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
		taskRunning, err := entStep.QueryProvisioningStepToAgentTask().Where(
			agenttask.StateNEQ(
				agenttask.StateCOMPLETE,
			),
		).Exist(ctx)

		if err != nil {
			log.Fatalf("Failed to Query Agent Task State. Err: %v", err)
		}

		if !taskRunning {
			break
		}

		time.Sleep(time.Second)
	}

	return nil
}
