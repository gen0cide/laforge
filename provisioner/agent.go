package provisioner

import (
	"os"
	"path/filepath"

	"github.com/kardianos/service"
)

var (
	// Agent is the primary singleton engine object
	Agent *Engine
)

const (
	// SvcName is the name of the service that will be installed
	SvcName = "laforge-agent"

	// SvcDisplayName is the proper name of the installed agent
	SvcDisplayName = `Laforge Provisioning Agent`

	// SvcDescription is the description of the agent
	SvcDescription = `Provides a programmatic interface for configuration of this host.`
)

// CreateAgentDir attempts to create the agent directory
func CreateAgentDir() error {
	return os.MkdirAll(AssetDir(), 0700)
}

// ConfigFile returns the absolute path to the config file
func ConfigFile() string {
	return filepath.Join(AgentHomeDir, `config.json`)
}

// StateFile returns the absolute path to the state file
func StateFile() string {
	return filepath.Join(AgentHomeDir, `state.json`)
}

// AssetDir returns the absolute path to the asset directory
func AssetDir() string {
	return filepath.Join(AgentHomeDir, `assets`)
}

// AssetPath returns the absolute path to an asset of the provided name
func AssetPath(name string) string {
	return filepath.Join(AssetDir(), name)
}

// Install installs the program as a system service on the local machine
func Install() error {
	return nil
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
