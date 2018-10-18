package agent

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kardianos/service"
)

var (
	// Agent is the primary singleton engine object
	Agent *Engine

	// AsyncWorker is the primary singleton object that does work
	AsyncWorker *Worker

	// Logger is the singleton logger of the Agent
	Logger *logrus.Logger

	// LogFile is the open file handle for the singleton logger
	LogFile *os.File
)

const (
	// SvcName is the name of the service that will be installed
	SvcName = "laforge-agent"

	// SvcDisplayName is the proper name of the installed agent
	SvcDisplayName = `Laforge Provisioning Agent`

	// SvcDescription is the description of the agent
	SvcDescription = `Provides a programmatic interface for configuration of this host.`
)

// StartLogger starts the file logger
func StartLogger() error {
	os.MkdirAll(StepLogDir(), 0755)
	LogFile, err := os.OpenFile(LogFilePath(), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	Logger := logrus.New()
	Logger.SetOutput(LogFile)
	return nil
}

// CreateAgentDir attempts to create the agent directory
func CreateAgentDir() error {
	return os.MkdirAll(AssetDir(), 0700)
}

// ConfigFile returns the absolute path to the config file
func ConfigFile() string {
	return filepath.Join(AgentHomeDir, `config.json`)
}

// StepLogDir is where step stdout/stderr is logged to
func StepLogDir() string {
	return filepath.Join(AgentHomeDir, "logs")
}

// InitFile returns the absolute path to the initialized file
func InitFile() string {
	return filepath.Join(AgentHomeDir, `initialized.txt`)
}

// LogFilePath returns the absolute path of the agent's log
func LogFilePath() string {
	return filepath.Join(AgentHomeDir, "agent.log")
}

// AssetDir returns the absolute path to the asset directory
func AssetDir() string {
	return filepath.Join(AgentHomeDir, `assets`)
}

// AssetPath returns the absolute path to an asset of the provided name
func AssetPath(name string) string {
	return filepath.Join(AssetDir(), name)
}

// Initialized describes whether the local machine has been initialized
func Initialized() bool {
	if _, err := os.Stat(InitFile()); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// TouchInitFile touches the systems initialized file
func TouchInitFile() error {
	return ioutil.WriteFile(InitFile(), []byte(fmt.Sprintf("%v", time.Now().UTC())), 0600)
}

// GetServiceConfig returns the installable service config
func GetServiceConfig() *service.Config {
	return &service.Config{
		Name:             SvcName,
		DisplayName:      SvcDisplayName,
		Description:      SvcDescription,
		Arguments:        []string{"run"},
		WorkingDirectory: AgentHomeDir,
		Executable:       ExePath,
	}
}

// GetService returns the service object
func GetService() (service.Service, error) {
	if Agent == nil {
		Agent = NewEngine()
	}
	return service.New(Agent, GetServiceConfig())
}
