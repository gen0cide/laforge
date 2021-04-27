package planner

import (
	"context"

	"github.com/gen0cide/laforge/ent"
)

func DeleteBuild(ctx context.Context, client *ent.Client, entBuild *ent.Build) (bool, error) {
	err := client.Build.DeleteOne(entBuild).Exec(ctx)
	if err != nil {
		return false, nil
	}
	return true, nil

}
