package provisioner

import (
	"github.com/gin-gonic/gin"
)

// ReqInitialize is used to initialize the configuration state
func (e *Engine) ReqInitialize(c *gin.Context) {
	err := e.LoadConfig()
	if err != nil {
		c.JSON(500, map[string]string{
			"status": "error",
			"path":   c.Request.RequestURI,
			"source": c.ClientIP(),
			"error":  err.Error(),
		})
		return
	}
	c.JSON(200, map[string]string{
		"status":          "ok",
		"path":            c.Request.RequestURI,
		"source":          c.ClientIP(),
		"config_hostname": e.Config.Host.Hostname,
	})
}

// ReqGetStatus returns the current agent status
func (e *Engine) ReqGetStatus(c *gin.Context) {
	c.JSON(200, e.GetStatus())
}

// ReqGetState returns a full dump of the agent's state
func (e *Engine) ReqGetState(c *gin.Context) {
	c.JSON(200, e.Config)
}

// ReqGetSteps returns a full dump of the agent's steps
func (e *Engine) ReqGetSteps(c *gin.Context) {
	c.JSON(200, map[string]string{
		"status": "ok",
		"path":   c.Request.RequestURI,
		"source": c.ClientIP(),
	})
}

// ReqGetCurrentStep returns the currently working step (if there is one)
func (e *Engine) ReqGetCurrentStep(c *gin.Context) {
	c.JSON(200, map[string]string{
		"status": "ok",
		"path":   c.Request.RequestURI,
		"source": c.ClientIP(),
	})
}

// ReqGetCompletedSteps returns the list of completed steps
func (e *Engine) ReqGetCompletedSteps(c *gin.Context) {
	c.JSON(200, map[string]string{
		"status": "ok",
		"path":   c.Request.RequestURI,
		"source": c.ClientIP(),
	})
}

// ReqGetAwaitingSteps returns the list of steps awaiting run
func (e *Engine) ReqGetAwaitingSteps(c *gin.Context) {
	c.JSON(200, map[string]string{
		"status": "ok",
		"path":   c.Request.RequestURI,
		"source": c.ClientIP(),
	})
}

// ReqGetSpecificStep returns details about a specified step
func (e *Engine) ReqGetSpecificStep(c *gin.Context) {
	c.JSON(200, map[string]string{
		"status":  "ok",
		"path":    c.Request.RequestURI,
		"source":  c.ClientIP(),
		"step_id": c.Param("id"),
	})
}

// ReqGetStepLogStdout returns the step's stdout, streaming if it's in progress
func (e *Engine) ReqGetStepLogStdout(c *gin.Context) {
	c.JSON(200, map[string]string{
		"status":  "ok",
		"path":    c.Request.RequestURI,
		"source":  c.ClientIP(),
		"step_id": c.Param("id"),
	})
}

// ReqGetStepLogStderr returns the step's stderr, streaming if it's in progress
func (e *Engine) ReqGetStepLogStderr(c *gin.Context) {
	c.JSON(200, map[string]string{
		"status":  "ok",
		"path":    c.Request.RequestURI,
		"source":  c.ClientIP(),
		"step_id": c.Param("id"),
	})
}

// ReqGetStepLogAll returns the step's logs, combined into one streaming chunk
func (e *Engine) ReqGetStepLogAll(c *gin.Context) {
	c.JSON(200, map[string]string{
		"status":  "ok",
		"path":    c.Request.RequestURI,
		"source":  c.ClientIP(),
		"step_id": c.Param("id"),
	})
}

// ReqPushProvision is what kicks off the provisioning process
func (e *Engine) ReqPushProvision(c *gin.Context) {
	c.JSON(200, map[string]string{
		"status":  "ok",
		"path":    c.Request.RequestURI,
		"source":  c.ClientIP(),
		"step_id": c.Param("id"),
	})
}

// ReqSelfDestruct removes all sensitive data and deletes the server
func (e *Engine) ReqSelfDestruct(c *gin.Context) {
	c.JSON(200, map[string]string{
		"status":  "ok",
		"path":    c.Request.RequestURI,
		"source":  c.ClientIP(),
		"step_id": c.Param("id"),
	})
}
