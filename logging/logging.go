package logging

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/gen0cide/laforge/ent"
	"github.com/sirupsen/logrus"
)

// Logger is to create an object used for logging to independent log files
type Logger struct {
	// Log is the logrus object used to write to the logs
	Log *logrus.Logger
	// LogFile is the absolute path to the log file
	LogFile string
}

// CreateNewLogger creates a Logger object set to output to the specified log file
func CreateNewLogger(logFilePath string) Logger {
	var log = logrus.New()
	log.SetLevel(logrus.GetLevel())

	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		log.Out = file
	} else {
		logrus.Error("Failed create log file")
	}

	return Logger{
		Log:     log,
		LogFile: logFilePath,
	}
}

func CreateLoggerForServerTask(serverTask *ent.ServerTask) (*Logger, error) {
	ctx := context.Background()
	defer ctx.Done()
	logFolder, ok := os.LookupEnv("LAFORGE_LOG_FOLDER")
	if !ok {
		// Default log location
		logFolder = "/var/log/laforge"
	}
	absPath, err := filepath.Abs(logFolder)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"logFolder":    logFolder,
			"serverTaskId": serverTask.ID,
		}).Errorf("error getting absolute path from log folder: %v", err)
		return nil, fmt.Errorf("error creating logger: %v", err)
	}
	err = os.MkdirAll(absPath, 0755)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"logFolder":    logFolder,
			"serverTaskId": serverTask.ID,
		}).Errorf("error creating log folder: %v", err)
		return nil, fmt.Errorf("error creating logger: %v", err)
	}
	filename := fmt.Sprintf("%s_%s.lfglog", time.Now().Format("20060102-15-04-05"), serverTask.Type)
	logPath := path.Join(absPath, filename)
	logrus.Info(logPath)
	log := CreateNewLogger(logPath)

	err = serverTask.Update().SetLogFilePath(logPath).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while updating server task with log file path: %v", err)
	}
	return &log, nil
}
