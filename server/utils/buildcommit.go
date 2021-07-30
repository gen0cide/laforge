package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/buildcommit"
	"github.com/gen0cide/laforge/ent/plandiff"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

// CreateRootCommit creates the root commit on a build
func CreateRootCommit(client *ent.Client, rdb *redis.Client, entBuild *ent.Build) (*ent.BuildCommit, error) {
	ctx := context.Background()
	defer ctx.Done()

	buildPlans, err := entBuild.QueryBuildToPlan().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("error querying plans from build: %v", err)
	}

	rootCommit, err := client.BuildCommit.Create().
		SetType(buildcommit.TypeROOT).
		SetRevision(0).
		SetBuildCommitToBuild(entBuild).
		SetState(buildcommit.StatePLANNING).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating root commit: %v", err)
	}

	var planDiffErr error = nil
	for _, buildPlan := range buildPlans {
		numExistingDiffs, err := buildPlan.QueryPlanToPlanDiffs().Count(ctx)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"buildPlan": buildPlan.ID,
			}).Errorf("error counting plan diffs from plan: %v", err)
			continue
		}
		_, planDiffErr := client.PlanDiff.Create().
			SetNewState(plandiff.NewStatePLANNING).
			SetRevision(numExistingDiffs).
			SetPlanDiffToBuildCommit(rootCommit).
			SetPlanDiffToPlan(buildPlan).
			Save(ctx)
		if planDiffErr != nil {
			logrus.WithFields(logrus.Fields{
				"buildPlan":  buildPlan.ID,
				"rootCommit": rootCommit.ID,
			}).Errorf("error creating plan diff: %v", err)
			continue
		}
	}
	if planDiffErr != nil {
		return nil, fmt.Errorf("error while generating plans (check logs): %v", planDiffErr)
	}
	rdb.Publish(ctx, "updatedBuildCommit", rootCommit.ID.String())

	return rootCommit, nil
}

// WaitForCommitReview halts program execution until a given build commit has been either approved or cancelled. Returns true if the commit was approved or false if the commit was cancelled or timeout was reached.
func WaitForCommitReview(client *ent.Client, entBuildCommit *ent.BuildCommit, timeout time.Duration) (bool, error) {
	ctx := context.Background()
	defer ctx.Done()
	startTime := time.Now()
	for {
		entBuildCommit, err := client.BuildCommit.Query().Where(buildcommit.IDEQ(entBuildCommit.ID)).Only(ctx)
		if err != nil {
			return false, err
		}

		// If the user has made a decision
		if entBuildCommit.State == buildcommit.StateCANCELLED {
			return false, nil
		} else if entBuildCommit.State == buildcommit.StateAPPROVED {
			return true, nil
		}

		// Check if we've timed out already
		if time.Since(startTime) >= timeout {
			break
		}
		// Otherwise, wait 1 second and then check again
		time.Sleep(1 * time.Second)
	}
	return false, nil
}
