package core

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/juju/utils/filepath"
	"github.com/pkg/errors"
)

// RemoteFileJob attempts to upload a file to the remote machine
// easyjson:json
type RemoteFileJob struct {
	GenericJob
	Target     *ProvisioningStep `json:"-"`
	RemoteFile *RemoteFile       `json:"-"`
	AssetPath  string            `json:"asset_path,omitempty"`
}

// CreateRemoteFileJob creates a new remote file job for a Doer object within the Planner
func CreateRemoteFileJob(id string, offset int, m *Metadata, pstep *ProvisioningStep) (*RemoteFileJob, error) {
	rj := &RemoteFileJob{
		Target: pstep,
	}

	rj.Metadata = m
	rj.MetadataID = m.GetID()
	rj.Offset = offset
	rj.JobID = id
	rj.RemoteFile = rj.Target.RemoteFile
	rj.JobType = "remote_file_job"
	rj.CreatedAt = time.Now()
	return rj, nil
}

// CanProceed implements the Doer interface
func (j *RemoteFileJob) CanProceed(e chan error) {
	if j.RemoteFile == nil || j.Target == nil {
		e <- errors.New("cannot proceed with remote file job with nil target or remote file objects")
		return
	}
	if j.Target.ProvisionedHost.Conn.Active {
		e <- nil
		return
	}
	timeout := time.After(time.Duration(j.Timeout) * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	currfp, err := filepath.NewRenderer("")
	if err != nil {
		e <- err
		return
	}

	pathToConnFile := currfp.Join(j.Base.BaseDir, j.Target.ParentLaforgeID(), "conn.laforge")

	fixed := false
	for {
		if fixed {
			break
		}
		select {
		case <-timeout:
			e <- fmt.Errorf("laforge was incapable of reaching remote host %s because its conn.laforge file was not available", j.Target.ParentLaforgeID())
			return
		case <-tick:
			conn := &Connection{}
			err := LoadHCLFromFile(pathToConnFile, conn)
			if err != nil {
				cli.Logger.Errorf("Error loading job %s resource: %v", j.JobID, err)
				continue
			}
			if conn.Active {
				err = j.Target.ProvisionedHost.Conn.Swap(conn)
				if err != nil {
					e <- fmt.Errorf("fatal error attempting to patch connection into state tree for %s: %v", j.JobID, err)
					return
				}
				fixed = true
				continue
			}
		}
	}
	e <- nil
}

// EnsureDependencies implements the Doer interface
func (j *RemoteFileJob) EnsureDependencies(e chan error) {
	currfp, err := filepath.NewRenderer("")
	if err != nil {
		e <- err
		return
	}
	assetfilename, err := j.RemoteFile.AssetName()
	if err != nil {
		e <- err
		return
	}

	targetAsset := currfp.Join(j.Base.BaseDir, j.Base.CurrentBuild.Path(), "data", assetfilename)
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
		if j.Target.ProvisionedHost.Conn.SSHAuthConfig.IdentityFile == sshKeyPath {
			cli.Logger.Debugf("Fixing identity file for %s", j.Target.Path())
			j.Target.ProvisionedHost.Conn.SSHAuthConfig.IdentityFile = currfp.Join(j.Base.BaseDir, j.Base.CurrentBuild.Path(), "data", "ssh.pem")
		}
	}

	e <- nil
}

// Do implements the Doer interface
func (j *RemoteFileJob) Do(e chan error) {
	cli.Logger.Warnf("Uploading remote file %s on %s to %s", j.AssetPath, j.Target.ProvisionedHost.Path(), j.RemoteFile.Destination)
	err := j.Target.ProvisionedHost.Conn.Upload(j.AssetPath, j.RemoteFile.Destination)
	if err != nil {
		cli.Logger.Errorf("Error uploading %s: %v", j.JobID, err)
		e <- err
		return
	}
	e <- nil
}

// CleanUp implements the Doer interface
func (j *RemoteFileJob) CleanUp(e chan error) {
	e <- nil
}

// Finish implements the Doer interface
func (j *RemoteFileJob) Finish(e chan error) {
	cli.Logger.Infof("Finished %s", j.JobID)
	e <- nil
}
