package spanner

// Job is a spanner job
type Job struct {
	Command    []string
	SaveOutput bool
	Silent     bool
	ExecType   string
	HostID     string
}
