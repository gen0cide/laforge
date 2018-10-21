package core

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

// ScriptJob attempts to upload and execute a script on the remote system
// easyjson:json
type ScriptJob struct {
	GenericJob
	Target    *ProvisioningStep `json:"-"`
	Script    *Script           `json:"-"`
	AssetPath string            `json:"asset_path,omitempty"`
}

// CreateScriptJob creates a new script job for a Doer object with the Planner
func CreateScriptJob(id string, offset int, m *Metadata, pstep *ProvisioningStep) (*ScriptJob, error) {
	sj := &ScriptJob{
		Target: pstep,
	}
	sj.Metadata = m
	sj.MetadataID = m.GetID()
	sj.Offset = offset
	sj.JobID = id
	sj.Script = sj.Target.Script
	sj.JobType = "script_job"
	sj.CreatedAt = time.Now()
	return sj, nil
}

// CanProceed implements the Doer interface
func (j *ScriptJob) CanProceed() error {
	if j.Script == nil || j.Target == nil {
		return errors.New("cannot proceed with script job with nil targets")
	}
	return nil
}

// EnsureDependencies implements the Doer interface
func (j *ScriptJob) EnsureDependencies(l *Laforge) error {
	targetAsset := filepath.Join(l.BaseDir, j.Target.ParentLaforgeID(), "assets", j.Script.SourceBase())
	if _, err := os.Stat(targetAsset); err != nil {
		return err
	}

	j.AssetPath = targetAsset

	if j.Target.ProvisionedHost.Conn == nil {
		return fmt.Errorf("script %s has a nil connection for the parent host", j.JobID)
	}

	return nil
}

// Do implements the Doer interface
func (j *ScriptJob) Do() error {
	return nil
}

// CleanUp implements the Doer interface
func (j *ScriptJob) CleanUp() error {
	return nil
}

// Finish implements the Doer interface
func (j *ScriptJob) Finish() error {
	return nil
}
