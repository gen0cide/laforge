package provisioner

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
)

// ServerPort is the port the agent API server runs on
var ServerPort = 9971

// Engine is the primary engine type within the provisioning agent
type Engine struct {
	Config  *State
	Current *State
	Server  *gin.Engine
	Service service.Service
}

// NewEngine creates a bare engine
func NewEngine() *Engine {
	e := &Engine{}
	return e
}

// Start implements the Interface type for service functionality
func (e *Engine) Start(s service.Service) error {
	e.Service = s
	go e.Serve()
	return nil
}

// Stop implements the Interface type for service functionality
func (e *Engine) Stop(s service.Service) error {
	return nil
}

// Serve is the long running loop for the Engine's API
func (e *Engine) Serve() {
	if !service.Interactive() {
		gin.DisableConsoleColor()
	}

	r := gin.Default()

	r.GET("/api/initialize", e.ReqInitialize)
	r.GET("/api/status", e.ReqGetStatus)
	r.GET("/api/state", e.ReqGetState)
	// r.GET("/api/steps", e.ReqGetState)
	// r.GET("/api/steps/current", e.ReqGetCurrentStep)
	// r.GET("/api/steps/completed", e.ReqGetCompletedSteps)
	// r.GET("/api/steps/awaiting", e.ReqGetAwaitingSteps)
	// r.GET("/api/steps/find/:id", e.ReqGetSpecificStep)
	// r.GET("/api/logs/all/:id", e.ReqGetStepLogAll)
	r.GET("/api/logs/stdout/:id", e.ReqGetStepLogStdout)
	r.GET("/api/logs/stderr/:id", e.ReqGetStepLogStderr)
	// r.POST("/api/provision", e.ReqPushProvision)
	r.POST("/api/self-destruct", e.ReqSelfDestruct)

	e.Server = r

	err := os.Chdir(AgentHomeDir)
	if err != nil {
		fmt.Printf("error entering agent home directory: %v\n", err)
		return
	}

	e.LoadConfig()

	e.Server.Run(fmt.Sprintf(":%d", ServerPort))
	return
}

// LoadConfig loads the base configuration of the host
func (e *Engine) LoadConfig() error {
	state, err := LoadStateFile(ConfigFile())
	if err != nil {
		return err
	}
	e.Config = state
	return nil
}

// GetStatus returns the current status of the engine
func (e *Engine) GetStatus() *Status {
	status := NewEmptyStatus()
	if e.Config == nil {
		status.Code = StatusBootingUp
		return status
	}

	status.TotalSteps = len(e.Config.Steps)
	status.CompletedSteps = len(e.Config.Completed)
	status.StartedAt = e.Config.InitializedAt

	switch e.Config.CurrentState {
	case "finished":
		status.Code = StatusIdle
		status.ElapsedTime = e.Config.CompletedAt.Sub(e.Config.InitializedAt)
		status.CompletedAt = e.Config.CompletedAt
		return status
	case "provisioning":
		status.Code = StatusRunningStep
		status.ElapsedTime = time.Since(e.Config.InitializedAt)
		status.CurrentStep = e.Config.CurrentStep
		return status
	}
	return status
}
