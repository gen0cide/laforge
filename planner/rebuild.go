package planner

import (
	"context"
	"sync"
	"time"

	"github.com/gen0cide/laforge/builder"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/plan"
	"github.com/gen0cide/laforge/ent/status"
	"github.com/sirupsen/logrus"
)

func Rebuild(ctx context.Context, client *ent.Client, entPlans []*ent.Plan) (bool, error) {
	env, err := entPlans[0].QueryPlanToBuild().QueryBuildToEnvironment().Only(ctx)
	// environment, err := entBuild.QueryBuildToEnvironment().Only(ctx)
	if err != nil {
		logrus.Errorf("error querying environment from build: %v", err)
		return false, err
	}

	genericBuilder, err := builder.BuilderFromEnvironment(env)
	if err != nil {
		logrus.Errorf("error generating builder: %v", err)
		return false, err
	}

	deleteContext := context.Background()
	// Mark all plans involved for rebuild
	for _, entPlan := range entPlans {
		err := markForDeleteRoutine(deleteContext, entPlan)
		if err != nil {
			return false, err
		}
	}

	// logrus.Debug("Marked plans for deletion, pausing for 1 minute so you can double check the statuses...")
	// time.Sleep(1 * time.Minute)

	var wg sync.WaitGroup
	for _, entPlan := range entPlans {
		wg.Add(1)
		go deleteRoutine(client, &genericBuilder, deleteContext, entPlan, &wg)
	}
	wg.Wait()

	logrus.Debug("waiting for deletion to propagate to all systems")
	time.Sleep(1 * time.Minute)

	buildContext := context.Background()
	for _, entPlan := range entPlans {
		err := markForRebuildRoutine(buildContext, entPlan)
		if err != nil {
			return false, err
		}
	}

	// logrus.Debug("Marked plans for build, pausing for 1 minute so you can double check the statuses...")
	// time.Sleep(1 * time.Minute)

	var wg2 sync.WaitGroup
	for _, entPlan := range entPlans {
		wg2.Add(1)
		go buildRoutine(client, &genericBuilder, buildContext, entPlan, &wg2)
	}
	wg2.Wait()

	return true, nil
}

func markForDeleteRoutine(ctx context.Context, entPlan *ent.Plan) error {
	entStatus, err := entPlan.QueryPlanToStatus().Only(ctx)
	if err != nil {
		return err
	}
	if entStatus.State != status.StateTODELETE {
		err = entStatus.Update().SetState(status.StateTODELETE).Exec(ctx)
		if err != nil {
			return err
		}
	}
	var provisionedStatus *ent.Status
	var getStatusError error = nil
	switch entPlan.Type {
	case plan.TypeStartBuild:
		build, getStatusError := entPlan.QueryPlanToBuild().Only(ctx)
		if getStatusError != nil {
			break
		}
		provisionedStatus, getStatusError = build.QueryBuildToStatus().Only(ctx)
	case plan.TypeStartTeam:
		team, getStatusError := entPlan.QueryPlanToTeam().Only(ctx)
		if getStatusError != nil {
			break
		}
		provisionedStatus, getStatusError = team.QueryTeamToStatus().Only(ctx)
	case plan.TypeProvisionNetwork:
		pnet, getStatusError := entPlan.QueryPlanToProvisionedNetwork().Only(ctx)
		if getStatusError != nil {
			break
		}
		provisionedStatus, getStatusError = pnet.QueryProvisionedNetworkToStatus().Only(ctx)
	case plan.TypeProvisionHost:
		phost, getStatusError := entPlan.QueryPlanToProvisionedHost().Only(ctx)
		if getStatusError != nil {
			break
		}
		provisionedStatus, getStatusError = phost.QueryProvisionedHostToStatus().Only(ctx)
	case plan.TypeExecuteStep:
		step, getStatusError := entPlan.QueryPlanToProvisioningStep().Only(ctx)
		if getStatusError != nil {
			break
		}
		provisionedStatus, getStatusError = step.QueryProvisioningStepToStatus().Only(ctx)
	default:
		break
	}
	if getStatusError != nil {
		logrus.Errorf("error getting status of provisioned object: %v", getStatusError)
	}
	if provisionedStatus.State != status.StateTODELETE {
		err = provisionedStatus.Update().SetState(status.StateTODELETE).Exec(ctx)
		if err != nil {
			return err
		}
	}
	nextPlans, err := entPlan.QueryNextPlan().All(ctx)
	if err != nil {
		return err
	}
	for _, nextPlan := range nextPlans {
		err = markForDeleteRoutine(ctx, nextPlan)
		if err != nil {
			return err
		}
	}
	return nil
}

func markForRebuildRoutine(ctx context.Context, entPlan *ent.Plan) error {
	entStatus, err := entPlan.QueryPlanToStatus().Only(ctx)
	if err != nil {
		return err
	}
	if entStatus.State != status.StateAWAITING {
		err = entStatus.Update().SetState(status.StateAWAITING).Exec(ctx)
		if err != nil {
			return err
		}
	}
	var provisionedStatus *ent.Status
	var getStatusError error = nil
	switch entPlan.Type {
	case plan.TypeStartBuild:
		build, getStatusError := entPlan.QueryPlanToBuild().Only(ctx)
		if getStatusError != nil {
			break
		}
		provisionedStatus, getStatusError = build.QueryBuildToStatus().Only(ctx)
	case plan.TypeStartTeam:
		team, getStatusError := entPlan.QueryPlanToTeam().Only(ctx)
		if getStatusError != nil {
			break
		}
		provisionedStatus, getStatusError = team.QueryTeamToStatus().Only(ctx)
	case plan.TypeProvisionNetwork:
		pnet, getStatusError := entPlan.QueryPlanToProvisionedNetwork().Only(ctx)
		if getStatusError != nil {
			break
		}
		provisionedStatus, getStatusError = pnet.QueryProvisionedNetworkToStatus().Only(ctx)
	case plan.TypeProvisionHost:
		phost, getStatusError := entPlan.QueryPlanToProvisionedHost().Only(ctx)
		if getStatusError != nil {
			break
		}
		provisionedStatus, getStatusError = phost.QueryProvisionedHostToStatus().Only(ctx)
	case plan.TypeExecuteStep:
		step, getStatusError := entPlan.QueryPlanToProvisioningStep().Only(ctx)
		if getStatusError != nil {
			break
		}
		provisionedStatus, getStatusError = step.QueryProvisioningStepToStatus().Only(ctx)
	default:
		break
	}
	if getStatusError != nil {
		logrus.Errorf("error getting status of provisioned object: %v", getStatusError)
	}
	if provisionedStatus.State != status.StateAWAITING {
		err = provisionedStatus.Update().SetState(status.StateAWAITING).Exec(ctx)
		if err != nil {
			return err
		}
	}
	nextPlans, err := entPlan.QueryNextPlan().All(ctx)
	if err != nil {
		return err
	}
	for _, nextPlan := range nextPlans {
		err = markForRebuildRoutine(ctx, nextPlan)
		if err != nil {
			return err
		}
	}
	return nil
}
