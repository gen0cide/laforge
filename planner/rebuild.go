package planner

import (
	"context"
	"sync"
	"time"

	"github.com/gen0cide/laforge/builder"
	"github.com/gen0cide/laforge/ent"
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

	// Mark all plans involved for rebuild
	for _, entPlan := range entPlans {
		err := markForDeleteRoutine(ctx, entPlan)
		if err != nil {
			return false, err
		}
	}

	// logrus.Debug("Marked plans for deletion, pausing for 1 minute so you can double check the statuses...")
	// time.Sleep(1 * time.Minute)

	var wg sync.WaitGroup
	for _, entPlan := range entPlans {
		wg.Add(1)
		go deleteRoutine(client, &genericBuilder, ctx, entPlan, &wg)
	}
	wg.Wait()

	logrus.Debug("waiting for deletion to propagate to all systems")
	time.Sleep(1 * time.Minute)

	for _, entPlan := range entPlans {
		err := markForRebuildRoutine(ctx, entPlan)
		if err != nil {
			return false, err
		}
	}

	// logrus.Debug("Marked plans for build, pausing for 1 minute so you can double check the statuses...")
	// time.Sleep(1 * time.Minute)

	var wg2 sync.WaitGroup
	for _, entPlan := range entPlans {
		wg2.Add(1)
		go buildRoutine(client, &genericBuilder, ctx, entPlan, &wg2)
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
