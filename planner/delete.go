package planner

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gen0cide/laforge/builder"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/agenttask"
	"github.com/gen0cide/laforge/ent/plan"
	"github.com/gen0cide/laforge/ent/predicate"
	"github.com/gen0cide/laforge/ent/provisionedhost"
	"github.com/gen0cide/laforge/ent/status"
	"github.com/sirupsen/logrus"
)

func DeleteBuild(ctx context.Context, client *ent.Client, entBuild *ent.Build) (bool, error) {
	entPlans, err := entBuild.QueryBuildToPlan().All(ctx)
	if err != nil {
		return false, err
	}

	var wg sync.WaitGroup
	for _, entPlan := range entPlans {
		planStatus, err := entPlan.PlanToStatus(ctx)
		if err != nil {
			return false, err
		}

		wg.Add(1)

		go func(wg *sync.WaitGroup, planStatus *ent.Status) {
			defer wg.Done()
			planStatus.Update().SetState(status.StateTODELETE).Save(ctx)
		}(&wg, planStatus)
	}

	wg.Wait()

	rootPlans, err := entBuild.QueryBuildToPlan().Where(plan.TypeEQ(plan.TypeStartBuild)).All(ctx)
	if err != nil {
		logrus.Errorf("error querying root plans from build: %v", err)
		return false, err
	}
	logrus.Infof("ROOT PLANS: %v", rootPlans)
	environment, err := entBuild.QueryBuildToEnvironment().Only(ctx)
	if err != nil {
		logrus.Errorf("error querying environment from build: %v", err)
		return false, err
	}

	genericBuilder, err := builder.BuilderFromEnvironment(environment)
	if err != nil {
		logrus.Errorf("error generating builder: %v", err)
		return false, err
	}

	logrus.WithFields(logrus.Fields{
		"rootPlanCount": len(rootPlans),
	}).Debug("found root plans")

	for _, entPlan := range rootPlans {
		wg.Add(1)
		go deleteRoutine(client, &genericBuilder, ctx, entPlan, &wg)
	}

	wg.Wait()

	logrus.Debug("delete build done")

	// Remove all rendered files
	err = os.RemoveAll(environment.Name + "/" + fmt.Sprint(entBuild.Revision))
	if err != nil {
		return false, fmt.Errorf("error deleting rendered files: %v", err)
	}
	// err = client.Build.DeleteOne(entBuild).Exec(ctx)
	// if err != nil {
	// 	return false, err
	// }
	return true, nil
}

func deleteRoutine(client *ent.Client, builder *builder.Builder, ctx context.Context, entPlan *ent.Plan, wg *sync.WaitGroup) {
	defer wg.Done()

	planStatus, err := entPlan.QueryPlanToStatus().Only(ctx)
	if err != nil {
		logrus.Errorf("error while getting plan status: %v", err)
		return
	}
	if planStatus.State != status.StateTODELETE {
		return
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

	// Only allow tree spidering in a specific order (don't follow dependency links)
	var planFilter predicate.Plan
	switch entPlan.Type {
	case plan.TypeStartBuild:
		planFilter = plan.TypeEQ(plan.TypeStartTeam)
	case plan.TypeStartTeam:
		planFilter = plan.TypeEQ(plan.TypeProvisionNetwork)
	case plan.TypeProvisionNetwork:
		planFilter = plan.TypeEQ(plan.TypeProvisionHost)
	case plan.TypeProvisionHost:
		planFilter = plan.TypeEQ(plan.TypeExecuteStep)
	case plan.TypeExecuteStep:
		planFilter = plan.TypeEQ(plan.TypeExecuteStep)
	default:
		break
	}
	nextPlans, err := entPlan.QueryNextPlan().Where(planFilter).All(ctx)
	if err != nil {
		logrus.Errorf("error querying next plan from ent plan: %v", err)
		return
	}

	logrus.Debugf("start delete | %s - %s", entPlan.Type, entPlan.ID)
	// logrus.Infof("next plans   | %s - %s | %v", entPlan.Type, entPlan.ID, nextPlans)

	var nextPlanWg sync.WaitGroup
	for _, nextPlan := range nextPlans {
		nextPlanWg.Add(1)
		go deleteRoutine(client, builder, ctx, nextPlan, &nextPlanWg)
	}
	nextPlanWg.Wait()

	logrus.Debugf("wait childs  | %s - %s", entPlan.Type, entPlan.ID)
	for {
		hasTaintedNextPlans, err := entPlan.QueryNextPlan().Where(
			plan.And(
				planFilter,
				plan.HasPlanToStatusWith(status.StateEQ(status.StateTAINTED)),
			),
		).Exist(ctx)

		if err != nil {
			logrus.Errorf("error checking for nextPlans that are TAINTED: %v", err)
			return
		}

		if hasTaintedNextPlans {
			logrus.Errorf("error: children are TAINTED for entPlan %s", entPlan.ID)
			entStatus, err := entPlan.PlanToStatus(ctx)
			if err != nil {
				logrus.Errorf("error querying status from ent plan: %v", err)
				return
			}
			_, err = entStatus.Update().SetState(status.StateTAINTED).Save(ctx)
			if err != nil {
				logrus.Errorf("error updating ent plan status to TAINTED: %v", err)
				return
			}
			_, err = provisionedStatus.Update().SetState(status.StateTAINTED).Save(ctx)
			if err != nil {
				logrus.Errorf("error updating provisioned object status to TAINTED: %v", err)
				return
			}
			return
		}

		hasUnDeletedNextPlans, err := entPlan.QueryNextPlan().Where(
			plan.And(
				planFilter,
				plan.HasPlanToStatusWith(
					status.Or(
						status.StateEQ(status.StateTODELETE),
						status.StateEQ(status.StateDELETEINPROGRESS),
					),
				),
			),
		).Exist(ctx)

		if err != nil {
			logrus.Errorf("error checking for nextPlans that are not DELETE: %v", err)
			return
		}

		if !hasUnDeletedNextPlans {
			break
		}

		time.Sleep(time.Second)
	}

	logrus.Debugf("fr deleting  | %s - %s", entPlan.Type, entPlan.ID)

	entStatus, err := entPlan.PlanToStatus(ctx)
	if err != nil {
		logrus.Errorf("error querying status from ent plan: %v", err)
		return
	}

	// Just double check to make sure it already hasn't been deleted
	if entStatus.State == status.StateDELETEINPROGRESS || entStatus.State == status.StateDELETED {
		return
	}

	entStatus, err = entStatus.Update().SetState(status.StateDELETEINPROGRESS).Save(ctx)
	if err != nil {
		logrus.Errorf("error updating ent plan status: %v", err)
		return
	}

	var deleteErr error = nil
	switch entPlan.Type {
	case plan.TypeStartBuild:
		deleteErr = provisionedStatus.Update().SetState(status.StateDELETED).Exec(ctx)
	case plan.TypeStartTeam:
		deleteErr = provisionedStatus.Update().SetState(status.StateDELETED).Exec(ctx)
	case plan.TypeProvisionNetwork:
		entProNetwork, err := entPlan.PlanToProvisionedNetwork(ctx)
		if err != nil {
			logrus.Errorf("error querying provisioned network from ent plan: %v", err)
			return
		}
		logrus.Debugf("del network  | %s", entPlan.ID)
		deleteErr = deleteNetwork(client, builder, ctx, entProNetwork)
	case plan.TypeProvisionHost:
		entProHost, err := entPlan.PlanToProvisionedHost(ctx)
		if err != nil {
			logrus.Errorf("error querying provisioned host from ent plan: %v", err)
			return
		}
		logrus.Debugf("del host     | %s", entPlan.ID)
		deleteErr = deleteHost(client, builder, ctx, entProHost)
	case plan.TypeExecuteStep:
		step, deleteErr := entPlan.QueryPlanToProvisioningStep().Only(ctx)
		if deleteErr != nil {
			break
		}
		ginFileMiddleware, deleteErr := step.QueryProvisioningStepToGinFileMiddleware().Only(ctx)
		if deleteErr != nil {
			break
		}
		deleteErr = client.GinFileMiddleware.DeleteOne(ginFileMiddleware).Exec(ctx)
		if deleteErr != nil {
			break
		}
		deleteErr = provisionedStatus.Update().SetState(status.StateDELETED).Exec(ctx)
	default:
		break
	}

	if deleteErr != nil {
		entStatus.Update().SetState(status.StateTAINTED).SetFailed(true).Save(ctx)
		logrus.WithFields(logrus.Fields{
			"type":    entPlan.Type,
			"builder": (*builder).ID(),
		}).Errorf("error while deleting plan: %v", deleteErr)
	} else {
		logrus.Debugf("del ent plan | %s - %s", entPlan.Type, entPlan.ID)
		_, deleteErr = entStatus.Update().SetState(status.StateDELETED).Save(ctx)
		if deleteErr != nil {
			logrus.Errorf("error while setting entStatus to DELETED: %v", err)
			return
		}
	}
}

func deleteHost(client *ent.Client, builder *builder.Builder, ctx context.Context, entProHost *ent.ProvisionedHost) error {
	logrus.Infof("del host     | %s", entProHost.SubnetIP)
	hostStatus, err := entProHost.QueryProvisionedHostToStatus().Only(ctx)
	if err != nil {
		logrus.Errorf("Error while getting Provisioned Host status: %v", err)
		return err
	}
	err = (*builder).TeardownHost(ctx, entProHost)
	if err != nil {
		// Tainted state tells us something went wrong with deletion
		logrus.Errorf("error while deleting host: %v", err)
		_, saveErr := hostStatus.Update().SetState(status.StateTAINTED).Save(ctx)
		if saveErr != nil {
			logrus.Errorf("error while setting Provisioned Host status to TAINTED: %v", saveErr)
			return saveErr
		}
		return err
	} else {
		_, saveErr := hostStatus.Update().SetState(status.StateDELETED).Save(ctx)
		if saveErr != nil {
			logrus.Errorf("error while setting Provisioned Host status to DELETED: %v", saveErr)
			return saveErr
		}
	}
	logrus.Infof("deleted %s successfully", entProHost.SubnetIP)

	// Cleanup agent tasks
	_, deleteErr := client.AgentTask.Delete().Where(agenttask.HasAgentTaskToProvisionedHostWith(provisionedhost.IDEQ(entProHost.ID))).Exec(ctx)
	if deleteErr != nil {
		logrus.Errorf("error while deleting Agent Tasks for Provisioned Host: %v", err)
		return deleteErr
	}
	return nil
}

func deleteNetwork(client *ent.Client, builder *builder.Builder, ctx context.Context, entProNetwork *ent.ProvisionedNetwork) error {
	logrus.Infof("del network  | %s", entProNetwork.Name)
	networkStatus, err := entProNetwork.QueryProvisionedNetworkToStatus().Only(ctx)
	if err != nil {
		logrus.Errorf("Error while getting Provisioned Network status: %v", err)
		return err
	}
	err = (*builder).TeardownNetwork(ctx, entProNetwork)
	if err != nil {
		logrus.Errorf("error while deleteing network: %v", err)
		_, saveErr := networkStatus.Update().SetState(status.StateTAINTED).Save(ctx)
		if saveErr != nil {
			logrus.Errorf("error while setting Provisioned Network status to TAINTED: %v", saveErr)
			return saveErr
		}
		return err
	} else {
		_, saveErr := networkStatus.Update().SetState(status.StateDELETED).Save(ctx)
		if saveErr != nil {
			logrus.Errorf("error while setting Provisioned Network status to DELETED: %v", saveErr)
			return saveErr
		}
	}
	logrus.Infof("deleted %s successfully", entProNetwork.Name)
	return nil
}
