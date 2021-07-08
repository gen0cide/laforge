package utils

import (
	"context"
	"log"

	"github.com/gen0cide/laforge/ent"
	"github.com/google/uuid"
)

func CreateTempURL(ctx context.Context, client *ent.Client, filePath string) (*ent.GinFileMiddleware, error) {
	entGinURL, err := client.GinFileMiddleware.Create().
		SetFilePath(filePath).
		SetURLID(uuid.New().String()).
		Save(ctx)
	if err != nil {
		log.Fatalf("Unable to generate temp url for %v. Err: %v", filePath, err)
		return nil, err
	}
	return entGinURL, nil
}
