package utils

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/servertask"
	"github.com/gen0cide/laforge/ent/status"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	UserKeyPath = path.Join("users", "%s", "%s", "keys")
	RepoPath    = path.Join("repos", "%s", "%s")
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

// CreateServerTask creates a server task under the given user with the given task type
func CreateServerTask(ctx context.Context, client *ent.Client, rdb *redis.Client, authUser *ent.AuthUser, taskType servertask.Type) (taskStatus *ent.Status, serverTask *ent.ServerTask, err error) {
	taskStatus, err = client.Status.Create().
		SetState(status.StateINPROGRESS).
		SetStartedAt(time.Now()).
		SetStatusFor(status.StatusForServerTask).
		Save(ctx)
	if err != nil {
		err = fmt.Errorf("error while creating server task status: %v", err)
		return
	}
	cwd, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("error while getting current working directory: %v", err)
		return
	}
	serverTask, err = client.ServerTask.Create().
		SetType(taskType).
		SetServerTaskToAuthUser(authUser).
		SetServerTaskToStatus(taskStatus).
		SetStartTime(time.Now()).
		SetLogFilePath(path.Join(cwd, "logs", fmt.Sprintf("%s_%s_%s.log", time.Now().Format("20060102-15-04-05"), taskType, authUser.ID))).
		Save(ctx)
	if err != nil {
		err = fmt.Errorf("error while creating server task: %v", err)
		return
	}
	// Push update to subscription
	rdb.Publish(ctx, "updatedStatus", taskStatus.ID.String())
	rdb.Publish(ctx, "updatedServerTask", serverTask.ID.String())
	return
}

// FailServerTask puts a server task into a failed state and set the errors on the task
func FailServerTask(ctx context.Context, client *ent.Client, rdb *redis.Client, taskStatus *ent.Status, serverTask *ent.ServerTask, errors ...error) (updatedTaskStatus *ent.Status, updatedServerTask *ent.ServerTask, err error) {
	updatedTaskStatus, err = client.Status.UpdateOne(taskStatus).SetState(status.StateFAILED).SetEndedAt(time.Now()).Save(ctx)
	if err != nil {
		err = fmt.Errorf("error while updating server task status: %v", err)
		return
	}
	errorStrings := make([]string, len(errors))
	for i, e := range errors {
		errorStrings[i] = e.Error()
	}
	updatedServerTask, err = client.ServerTask.UpdateOne(serverTask).SetErrors(errorStrings).SetEndTime(time.Now()).Save(ctx)
	if err != nil {
		err = fmt.Errorf("error while updating server task: %v", err)
		return
	}
	// Push update to subscription
	rdb.Publish(ctx, "updatedStatus", taskStatus.ID.String())
	rdb.Publish(ctx, "updatedServerTask", serverTask.ID.String())
	return
}

// CompleteServerTask puts a server task into the completed state
func CompleteServerTask(ctx context.Context, client *ent.Client, rdb *redis.Client, taskStatus *ent.Status, serverTask *ent.ServerTask) (updatedTaskStatus *ent.Status, updatedServerTask *ent.ServerTask, err error) {
	updatedTaskStatus, err = client.Status.UpdateOne(taskStatus).SetState(status.StateCOMPLETE).SetEndedAt(time.Now()).Save(ctx)
	if err != nil {
		err = fmt.Errorf("error while updating server task status: %v", err)
		return
	}
	updatedServerTask, err = client.ServerTask.UpdateOne(serverTask).SetEndTime(time.Now()).Save(ctx)
	if err != nil {
		err = fmt.Errorf("error while updating server task: %v", err)
		return
	}
	// Push update to subscription
	rdb.Publish(ctx, "updatedStatus", taskStatus.ID.String())
	rdb.Publish(ctx, "updatedServerTask", serverTask.ID.String())
	return
}
