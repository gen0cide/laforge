package provisioner

import (
	"time"
)

const (
	// StatusBootingUp is returned when the host is online, but not ready to receive API calls
	StatusBootingUp StatusCode = "BOOTING_UP"

	// StatusIdle is returned when the host is online, and the Engine's API is currently idle
	StatusIdle StatusCode = "IDLE"

	// StatusAwaitingReboot is returned when the host has been scheduled for a reboot, but it has not happened yet
	StatusAwaitingReboot StatusCode = "AWAITING_REBOOT"

	// StatusRunningStep is returned when a step is currently being run by the machine
	StatusRunningStep StatusCode = "RUNNING_STEP"

	// StatusRefreshing is returned in between steps being run to modify the state on disk
	StatusRefreshing StatusCode = "REFRESHING"

	// StatusDestroying is returned when the agent is in the process of destroying sensitive information
	StatusDestroying StatusCode = "DESTROYING"
)

// StatusCode is returned as part of an API response
type StatusCode string

// Status is a response object to API calls
type Status struct {
	Code           StatusCode
	StartedAt      time.Time
	ElapsedTime    time.Duration
	CurrentStep    *Step
	TotalSteps     int
	CompletedSteps int
}
