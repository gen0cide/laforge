package utils

import (
	"context"
	"fmt"

	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/buildcommit"
	"github.com/gen0cide/laforge/ent/plandiff"
	"github.com/sirupsen/logrus"
)

// CreateRootCommit creates the root commit on a build
func CreateRootCommit(client *ent.Client, entBuild *ent.Build) (*ent.BuildCommit, error) {
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

	return rootCommit, nil
}
