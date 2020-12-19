package core

import (
	"context"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/gen0cide/laforge/ent"
	"github.com/google/uuid"
)

// CreateTagEntry ...
func CreateTagEntry(name string, description map[string]string, ctx context.Context, client *ent.Client) (*ent.Tag, error) {
	tag, err := client.Tag.
		Create().
		SetUUID(uuid.New()).
		SetName(name).
		SetDescription(description).
		Save(ctx)

	if err != nil {
		cli.Logger.Debugf("failed creating tag: %v", err)
		return nil, err
	}

	cli.Logger.Debugf("tag was created: ", tag)
	return tag, nil
}
