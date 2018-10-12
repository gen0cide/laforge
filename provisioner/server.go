package provisioner

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kardianos/service"
)

// ServerPort is the port the agent API server runs on
var ServerPort = 9971

// Engine is the primary engine type within the provisioning agent
type Engine struct {
	Config  *State
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

	r.GET("/api/status", e.ReqGetStatus)
	r.GET("/api/state", e.ReqGetState)
	r.GET("/api/steps", e.ReqGetState)
	r.GET("/api/steps/current", e.ReqGetCurrentStep)
	r.GET("/api/steps/completed", e.ReqGetCompletedSteps)
	r.GET("/api/steps/awaiting", e.ReqGetAwaitingSteps)
	r.GET("/api/steps/find/:id", e.ReqGetSpecificStep)
	r.GET("/api/logs/all/:id", e.ReqGetStepLogAll)
	r.GET("/api/logs/stdout/:id", e.ReqGetStepLogStdout)
	r.GET("/api/logs/stderr/:id", e.ReqGetStepLogStderr)
	r.POST("/api/provision", e.ReqPushProvision)
	r.POST("/api/self-destruct", e.ReqSelfDestruct)

	e.Server = r

	err := os.Chdir(AgentHomeDir)
	if err != nil {
		fmt.Printf("error entering agent home directory: %v\n", err)
		return
	}

	err = e.LoadConfig()
	if err != nil {
		fmt.Printf("error loading config: %v\n", err)
		return
	}

	e.Server.Run(fmt.Sprintf(":%d", ServerPort))
	return
}

// LoadConfig loads the base configuration of the host
func (e *Engine) LoadConfig() error {
	fmt.Printf("CONFIG FILE: %s\n", ConfigFile())
	state, err := LoadStateFile(ConfigFile())
	if err != nil {
		return err
	}
	e.Config = state
	return nil
}
