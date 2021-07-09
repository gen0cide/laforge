package utils

import (
	"context"

	"github.com/gen0cide/laforge/ent"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func CreateTempURL(ctx context.Context, client *ent.Client, filePath string) (*ent.GinFileMiddleware, error) {
	entGinURL, err := client.GinFileMiddleware.Create().
		SetFilePath(filePath).
		SetURLID(uuid.New().String()).
		Save(ctx)
	if err != nil {
		logrus.Errorf("Unable to generate temp url for %v. Err: %v", filePath, err)
		return nil, err
	}
	return entGinURL, nil
}
