package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gen0cide/laforge/core/cli"
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
	if j.Target.ProvisionedHost.Conn.Active {
		return nil
		// return NewTimeoutExtension(errors.New("cannot proceed with a host with an inactive connection"))
	}

	pathToConnFile := filepath.Join(j.Base.BaseDir, j.Target.ParentLaforgeID(), "conn.laforge")

	logdir := filepath.Join(j.Base.BaseDir, j.Target.ParentLaforgeID(), "logs")
	if _, err := os.Stat(logdir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(logdir, 0755)
		} else {
			cli.Logger.Errorf("Error creating log directory %s: %v", logdir, err)
			return err
		}
	}

	if _, err := os.Stat(pathToConnFile); err != nil {
		if os.IsNotExist(err) {
			return NewTimeoutExtension(errors.New("cannot proceed with a host that has no connection definition"))
		}
		return err
	}

	conn := &Connection{}
	err := LoadHCLFromFile(pathToConnFile, conn)
	if err != nil {
		cli.Logger.Errorf("Error loading job %s resource: %v", j.JobID, err)
		return err
	}

	if conn.Active != true {
		return NewTimeoutExtension(errors.New("cannot proceed with a host with an inactive connection"))
	}

	err = j.Target.ProvisionedHost.Conn.Swap(conn)
	if err != nil {
		return fmt.Errorf("fatal error attempting to patch connection into state tree for %s: %v", j.JobID, err)
	}

	return nil
}

// EnsureDependencies implements the Doer interface
func (j *ScriptJob) EnsureDependencies() error {
	targetAsset := filepath.Join(j.Base.BaseDir, j.Target.ParentLaforgeID(), "assets", j.Script.SourceBase())
	if _, err := os.Stat(targetAsset); err != nil {
		return err
	}

	j.AssetPath = strings.TrimSpace(targetAsset)

	if j.Target.ProvisionedHost.Conn == nil {
		return fmt.Errorf("script %s has a nil connection for the parent host", j.JobID)
	}

	if j.Target.ProvisionedHost.Conn.IsSSH() {
		if j.Target.ProvisionedHost.Conn.SSHAuthConfig.IdentityFile == `../../data/ssh.pem` {
			cli.Logger.Debugf("Fixing identity file for %s", j.Target.Path())
			j.Target.ProvisionedHost.Conn.SSHAuthConfig.IdentityFile = filepath.Join(j.Base.BaseDir, j.Base.CurrentBuild.Path(), "data", "ssh.pem")
		}
	}

	return nil
}

// Do implements the Doer interface
func (j *ScriptJob) Do() error {
	cli.Logger.Warnf("Uploading and executing %s on %s", j.AssetPath, j.Target.ProvisionedHost.Conn.RemoteAddr)
	actualfilename := fmt.Sprintf("%d-%s", j.Target.StepNumber, filepath.Base(j.AssetPath))
	logdir := filepath.Join(j.Base.BaseDir, j.Target.ParentLaforgeID(), "logs")
	err := j.Target.ProvisionedHost.Conn.UploadExecuteAndDelete(j, j.AssetPath, actualfilename, logdir)
	if err != nil {
		cli.Logger.Errorf("Error executing %s: %v", j.JobID, err)
		return err
	}
	return nil
}

// CleanUp implements the Doer interface
func (j *ScriptJob) CleanUp() error {
	return nil
}

// Finish implements the Doer interface
func (j *ScriptJob) Finish() error {
	cli.Logger.Infof("Finished %s", j.JobID)
	return nil
}
