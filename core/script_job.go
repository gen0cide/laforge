package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"

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
func (j *ScriptJob) CanProceed(e chan error) {
	if j.Script == nil || j.Target == nil {
		e <- errors.New("cannot proceed with script job with nil targets")
		return
	}
	if j.Target.ProvisionedHost.Conn.Active {
		e <- nil
		return
		// return NewTimeoutExtension(errors.New("cannot proceed with a host with an inactive connection"))
	}

	pathToConnFile := filepath.Join(j.Base.BaseDir, j.Target.ParentLaforgeID(), "conn.laforge")

	logdir := filepath.Join(j.Base.BaseDir, j.Target.ParentLaforgeID(), "logs")
	if _, err := os.Stat(logdir); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(logdir, 0755)
		} else {
			cli.Logger.Errorf("Error creating log directory %s: %v", logdir, err)
			e <- err
			return
		}
	}

	if _, err := os.Stat(pathToConnFile); err != nil {
		if os.IsNotExist(err) {
			e <- NewTimeoutExtension(errors.New("cannot proceed with a host that has no connection definition"))
			return
		}
		e <- nil
		return
	}

	conn := &Connection{}
	err := LoadHCLFromFile(pathToConnFile, conn)
	if err != nil {
		cli.Logger.Errorf("Error loading job %s resource: %v", j.JobID, err)
		e <- err
		return
	}

	if conn.Active != true {
		e <- NewTimeoutExtension(errors.New("cannot proceed with a host with an inactive connection"))
		return
	}

	newConn, err := SmartMerge(j.Target.ProvisionedHost.Conn, conn, false)
	if err != nil {
		e <- fmt.Errorf("Error merging connections for %s", j.Target.ParentLaforgeID())
		return
	}
	j.Target.ProvisionedHost.Conn = newConn.(*Connection)

	if err != nil {
		e <- fmt.Errorf("fatal error attempting to patch connection into state tree for %s: %v", j.JobID, err)
		return
	}

	e <- nil
	return
}

// EnsureDependencies implements the Doer interface
func (j *ScriptJob) EnsureDependencies(e chan error) {
	targetAsset := filepath.Join(j.Base.BaseDir, j.Target.ParentLaforgeID(), "assets", j.Script.SourceBase())
	if _, err := os.Stat(targetAsset); err != nil {
		e <- err
		return
	}

	j.AssetPath = strings.TrimSpace(targetAsset)

	if j.Target.ProvisionedHost.Conn == nil {
		e <- fmt.Errorf("script %s has a nil connection for the parent host", j.JobID)
		return
	}

	if j.Target.ProvisionedHost.Conn.IsSSH() {
		if j.Target.ProvisionedHost.Conn.SSHAuthConfig.IdentityFile == `../../data/ssh.pem` {
			cli.Logger.Debugf("Fixing identity file for %s", j.Target.Path())
			j.Target.ProvisionedHost.Conn.SSHAuthConfig.IdentityFile = filepath.Join(j.Base.BaseDir, j.Base.CurrentBuild.Path(), "data", "ssh.pem")
		}
	}

	e <- nil
	return
}

// Do implements the Doer interface
func (j *ScriptJob) Do(e chan error) {
	cli.Logger.Warnf("Performing Script Job:\n  %s %s: %s\n  %s   %s: %s", color.HiBlueString(">>"), color.HiCyanString("SCRIPT"), color.HiGreenString("%s", j.AssetPath), color.HiBlueString(">>"), color.HiCyanString("HOST"), color.HiGreenString("%s", j.Target.ProvisionedHost.Conn.RemoteAddr))
	actualfilename := fmt.Sprintf("%d-%s", j.Target.StepNumber, filepath.Base(j.AssetPath))
	logdir := filepath.Join(j.Base.BaseDir, j.Target.ParentLaforgeID(), "logs")
	err := j.Target.ProvisionedHost.Conn.UploadExecuteAndDelete(j, j.AssetPath, actualfilename, logdir)
	if err != nil {
		cli.Logger.Errorf("Error executing %s: %v", j.JobID, err)
		e <- err
		return
	}
	e <- nil
	return
}

// CleanUp implements the Doer interface
func (j *ScriptJob) CleanUp(e chan error) {
	e <- nil
	return
}

// Finish implements the Doer interface
func (j *ScriptJob) Finish(e chan error) {
	cli.Logger.Infof("Finished %s", j.JobID)
	e <- nil
	return
}
