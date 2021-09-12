package planner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gen0cide/laforge/builder"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/buildcommit"
	"github.com/gen0cide/laforge/ent/plan"
	"github.com/gen0cide/laforge/ent/plandiff"
	"github.com/gen0cide/laforge/ent/predicate"
	"github.com/gen0cide/laforge/ent/status"
	"github.com/gen0cide/laforge/logging"
	"github.com/gen0cide/laforge/server/utils"
	"github.com/go-redis/redis/v8"
)

func Rebuild(client *ent.Client, rdb *redis.Client, logger *logging.Logger, currentUser *ent.AuthUser, serverTask *ent.ServerTask, taskStatus *ent.Status, entPlans []*ent.Plan, spawnedRebuildSuccessfully chan bool) (bool, error) {
	ctx := context.Background()
	defer ctx.Done()

	entBuild, err := entPlans[0].QueryPlanToBuild().Only(ctx)
	if err != nil {
		spawnedRebuildSuccessfully <- false
		logger.Log.Errorf("error getting build from plan: %v", err)
		return false, err
	}

	rebuildRevision, err := entBuild.QueryBuildToBuildCommits().Count(ctx)
	if err != nil {
		spawnedRebuildSuccessfully <- false
		logger.Log.Errorf("error counting commits on build: %v", err)
		return false, err
	}

	entRebuildCommit, err := client.BuildCommit.Create().
		SetRevision(rebuildRevision).
		SetType(buildcommit.TypeREBUILD).
		SetState(buildcommit.StatePLANNING).
		SetBuildCommitToBuild(entBuild).
		Save(ctx)
	if err != nil {
		spawnedRebuildSuccessfully <- false
		logger.Log.Errorf("error while creating rebuild commit: %v", err)
		return false, fmt.Errorf("error while creating rebuild commit: %v", err)
	}
	rdb.Publish(ctx, "updatedBuildCommit", entRebuildCommit.ID.String())
	err = entBuild.Update().SetBuildToLatestBuildCommit(entRebuildCommit).Exec(ctx)
	if err != nil {
		spawnedRebuildSuccessfully <- false
		logger.Log.Errorf("error while setting latest commit on build: %v", err)
		return false, fmt.Errorf("error while setting latest commit on build: %v", err)
	}
	rdb.Publish(ctx, "updatedBuild", entBuild.ID.String())

	for _, rootPlan := range entPlans {
		err = generateRebuildCommitPlans(client, ctx, rootPlan, entRebuildCommit)
		if err != nil {
			spawnedRebuildSuccessfully <- false
			logger.Log.Errorf("error generating plans for rebuild commit")
			return false, err
		}
		err = generateRebuildCommitPreviousPlans(client, ctx, rootPlan, entRebuildCommit)
		if err != nil {
			spawnedRebuildSuccessfully <- false
			logger.Log.Errorf("error generating previous plans for rebuild commit")
			return false, err
		}
	}

	spawnedRebuildSuccessfully <- true

	logger.Log.Debug("-----\nWAITING FOR COMMIT REVIEW\n-----")
	isApproved, err := utils.WaitForCommitReview(client, entRebuildCommit, 20*time.Minute)
	if err != nil {
		logger.Log.Errorf("error while waiting for rebuild commit to be reviewed: %v", err)
		entRebuildCommit.Update().SetState(buildcommit.StateCANCELLED).Exec(ctx)
		rdb.Publish(ctx, "updatedBuildCommit", entRebuildCommit.ID.String())
		return false, err
	}

	// Cancelled or timeout reached
	if !isApproved {
		logger.Log.Debug("-----\nCOMMIT CANCELLED/TIMED OUT\n-----")
		logger.Log.Errorf("rebuild commit has been cancelled or 20 minute timeout has been reached")
		err = entRebuildCommit.Update().SetState(buildcommit.StateCANCELLED).Exec(ctx)
		if err != nil {
			logger.Log.Errorf("error while cancelling rebuild commit: %v", err)
			return false, err
		}
		rdb.Publish(ctx, "updatedBuildCommit", entRebuildCommit.ID.String())
		taskStatus, serverTask, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask)
		if err != nil {
			return false, fmt.Errorf("error failing execute build server task: %v", err)
		}
		return false, fmt.Errorf("commit has been cancelled or 20 minute timeout has been reached")
	}
	logger.Log.Debug("-----\nCOMMIT APPROVED\n-----")

	env, err := entBuild.QueryBuildToEnvironment().Only(ctx)
	if err != nil {
		logger.Log.Errorf("error querying environment from build: %v", err)
		return false, err
	}

	err = entRebuildCommit.Update().SetState(buildcommit.StateINPROGRESS).Exec(ctx)
	if err != nil {
		taskStatus, serverTask, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask)
		if err != nil {
			return false, fmt.Errorf("error failing execute build server task: %v", err)
		}
		logger.Log.Errorf("error while cancelling rebuild commit: %v", err)
		return false, err
	}
	rdb.Publish(ctx, "updatedBuildCommit", entRebuildCommit.ID.String())

	genericBuilder, err := builder.BuilderFromEnvironment(env, logger)
	if err != nil {
		taskStatus, serverTask, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask)
		if err != nil {
			return false, fmt.Errorf("error failing execute build server task: %v", err)
		}
		logger.Log.Errorf("error generating builder: %v", err)
		return false, err
	}

	deleteContext := context.Background()
	// Mark all plans involved for rebuild
	// for _, entPlan := range entPlans {
	err = markForRoutine(deleteContext, logger, status.StateTODELETE, entRebuildCommit)
	if err != nil {
		taskStatus, serverTask, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask)
		if err != nil {
			return false, fmt.Errorf("error failing execute build server task: %v", err)
		}
		return false, err
	}
	// }

	// logger.Log.Debug("Marked plans for deletion, pausing for 1 minute so you can double check the statuses...")
	// time.Sleep(1 * time.Minute)

	var wg sync.WaitGroup
	for _, entPlan := range entPlans {
		wg.Add(1)
		go deleteRoutine(client, logger, &genericBuilder, deleteContext, entPlan, &wg)
	}
	wg.Wait()

	logger.Log.Debug("waiting for deletion to propagate to all systems")
	time.Sleep(1 * time.Minute)

	buildContext := context.Background()
	defer buildContext.Done()
	// for _, entPlan := range entPlans {
	err = markForRoutine(buildContext, logger, status.StateAWAITING, entRebuildCommit)
	if err != nil {
		taskStatus, serverTask, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask)
		if err != nil {
			return false, fmt.Errorf("error failing execute build server task: %v", err)
		}
		return false, err
	}
	// }

	// logrus.Debug("Marked plans for build, pausing for 1 minute so you can double check the statuses...")
	// time.Sleep(1 * time.Minute)

	var wg2 sync.WaitGroup
	for _, entPlan := range entPlans {
		wg2.Add(1)
		go buildRoutine(client, logger, &genericBuilder, buildContext, entPlan, &wg2)
	}
	wg2.Wait()

	err = entRebuildCommit.Update().SetState(buildcommit.StateAPPLIED).Exec(ctx)
	if err != nil {
		taskStatus, serverTask, err = utils.FailServerTask(ctx, client, rdb, taskStatus, serverTask)
		if err != nil {
			return false, fmt.Errorf("error failing execute build server task: %v", err)
		}
		logger.Log.Errorf("error while cancelling rebuild commit: %v", err)
		return false, err
	}
	rdb.Publish(ctx, "updatedBuildCommit", entRebuildCommit.ID.String())

	taskStatus, serverTask, err = utils.CompleteServerTask(deleteContext, client, rdb, taskStatus, serverTask)
	if err != nil {
		return false, fmt.Errorf("error completing execute build server task: %v", err)
	}
	return true, nil
}

func generateRebuildCommitPlans(client *ent.Client, ctx context.Context, rootPlan *ent.Plan, entBuildCommit *ent.BuildCommit) error {
	diffRevision, err := rootPlan.QueryPlanToPlanDiffs().Count(ctx)
	if err != nil {
		return err
	}

	planDiffExists, err := entBuildCommit.QueryBuildCommitToPlanDiffs().Where(plandiff.HasPlanDiffToPlanWith(plan.IDEQ(rootPlan.ID))).Exist(ctx)
	if err != nil {
		return err
	} else if !planDiffExists {
		_, err = client.PlanDiff.Create().
			SetNewState(plandiff.NewStateTOREBUILD).
			SetPlanDiffToBuildCommit(entBuildCommit).
			SetPlanDiffToPlan(rootPlan).
			SetRevision(diffRevision).
			Save(ctx)
	}

	nextPlans, err := rootPlan.QueryNextPlan().All(ctx)
	if err != nil {
		return err
	}

	for _, nextPlan := range nextPlans {
		err := generateRebuildCommitPlans(client, ctx, nextPlan, entBuildCommit)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateRebuildCommitPreviousPlans(client *ent.Client, ctx context.Context, rootPlan *ent.Plan, entBuildCommit *ent.BuildCommit) error {
	diffRevision, err := rootPlan.QueryPlanToPlanDiffs().Count(ctx)
	if err != nil {
		return err
	}

	planDiffExists, err := entBuildCommit.QueryBuildCommitToPlanDiffs().Where(plandiff.HasPlanDiffToPlanWith(plan.IDEQ(rootPlan.ID))).Exist(ctx)
	if err != nil {
		return err
	} else if !planDiffExists {
		_, err = client.PlanDiff.Create().
			SetNewState(plandiff.NewStateCOMPLETE).
			SetPlanDiffToBuildCommit(entBuildCommit).
			SetPlanDiffToPlan(rootPlan).
			SetRevision(diffRevision).
			Save(ctx)
	}

	var prevPlanPredicate predicate.Plan = nil
	switch rootPlan.Type {
	case plan.TypeProvisionHost:
		prevPlanPredicate = plan.TypeEQ(plan.TypeProvisionNetwork)
	case plan.TypeProvisionNetwork:
		prevPlanPredicate = plan.TypeEQ(plan.TypeStartTeam)
	default:
		return nil
	}

	var prevPlans []*ent.Plan
	if prevPlanPredicate != nil {
		prevPlans, err = rootPlan.QueryPrevPlan().Where(prevPlanPredicate).All(ctx)
	} else {
		prevPlans, err = rootPlan.QueryPrevPlan().All(ctx)
	}
	if err != nil {
		return err
	}

	for _, prevPlan := range prevPlans {
		err := generateRebuildCommitPreviousPlans(client, ctx, prevPlan, entBuildCommit)
		if err != nil {
			return err
		}
	}
	return nil
}

func markForRoutine(ctx context.Context, logger *logging.Logger, targetStatus status.State, entRebuildCommit *ent.BuildCommit) error {
	entPlanDiffs, err := entRebuildCommit.QueryBuildCommitToPlanDiffs().All(ctx)
	if err != nil {
		return err
	}
	for _, entPlanDiff := range entPlanDiffs {
		// Skip anything that isn't to be rebuilt
		if entPlanDiff.NewState != plandiff.NewStateTOREBUILD {
			continue
		}
		entPlan, err := entPlanDiff.QueryPlanDiffToPlan().Only(ctx)
		if err != nil {
			return err
		}
		entStatus, err := entPlan.QueryPlanToStatus().Only(ctx)
		if err != nil {
			return err
		}
		if entStatus.State != targetStatus {
			err = entStatus.Update().SetState(targetStatus).Exec(ctx)
			if err != nil {
				return err
			}
			rdb.Publish(ctx, "updatedStatus", entStatus.ID.String())
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
		}
		if getStatusError != nil {
			logger.Log.Errorf("error getting status of provisioned object: %v", getStatusError)
		}
		if provisionedStatus.State != targetStatus {
			err = provisionedStatus.Update().SetState(targetStatus).Exec(ctx)
			if err != nil {
				return err
			}
			rdb.Publish(ctx, "updatedStatus", provisionedStatus.ID.String())
		}
		// nextPlans, err := entPlan.QueryNextPlan().All(ctx)
		// if err != nil {
		// 	return err
		// }
		// for _, nextPlan := range nextPlans {
		// 	err = markForDeleteRoutine(ctx, nextPlan)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
	}
	return nil
}

// func markForRebuildRoutine(ctx context.Context, entPlan *ent.Plan) error {
// 	entStatus, err := entPlan.QueryPlanToStatus().Only(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	if entStatus.State != status.StateAWAITING {
// 		err = entStatus.Update().SetState(status.StateAWAITING).Exec(ctx)
// 		if err != nil {
// 			return err
// 		}
// 		rdb.Publish(ctx, "updatedStatus", entStatus.ID.String())
// 	}
// 	var provisionedStatus *ent.Status
// 	var getStatusError error = nil
// 	switch entPlan.Type {
// 	case plan.TypeStartBuild:
// 		build, getStatusError := entPlan.QueryPlanToBuild().Only(ctx)
// 		if getStatusError != nil {
// 			break
// 		}
// 		provisionedStatus, getStatusError = build.QueryBuildToStatus().Only(ctx)
// 	case plan.TypeStartTeam:
// 		team, getStatusError := entPlan.QueryPlanToTeam().Only(ctx)
// 		if getStatusError != nil {
// 			break
// 		}
// 		provisionedStatus, getStatusError = team.QueryTeamToStatus().Only(ctx)
// 	case plan.TypeProvisionNetwork:
// 		pnet, getStatusError := entPlan.QueryPlanToProvisionedNetwork().Only(ctx)
// 		if getStatusError != nil {
// 			break
// 		}
// 		provisionedStatus, getStatusError = pnet.QueryProvisionedNetworkToStatus().Only(ctx)
// 	case plan.TypeProvisionHost:
// 		phost, getStatusError := entPlan.QueryPlanToProvisionedHost().Only(ctx)
// 		if getStatusError != nil {
// 			break
// 		}
// 		provisionedStatus, getStatusError = phost.QueryProvisionedHostToStatus().Only(ctx)
// 	case plan.TypeExecuteStep:
// 		step, getStatusError := entPlan.QueryPlanToProvisioningStep().Only(ctx)
// 		if getStatusError != nil {
// 			break
// 		}
// 		provisionedStatus, getStatusError = step.QueryProvisioningStepToStatus().Only(ctx)
// 	default:
// 		break
// 	}
// 	if getStatusError != nil {
// 		logrus.Errorf("error getting status of provisioned object: %v", getStatusError)
// 	}
// 	if provisionedStatus.State != status.StateAWAITING {
// 		err = provisionedStatus.Update().SetState(status.StateAWAITING).Exec(ctx)
// 		if err != nil {
// 			return err
// 		}
// 		rdb.Publish(ctx, "updatedStatus", provisionedStatus.ID.String())
// 	}
// 	nextPlans, err := entPlan.QueryNextPlan().All(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	for _, nextPlan := range nextPlans {
// 		err = markForRebuildRoutine(ctx, nextPlan)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
