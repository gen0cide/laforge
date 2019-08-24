package core

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"

	"github.com/gen0cide/laforge/core/cli"
)

// CommandJob attempts to execute remote commands on the system
// easyjson:json
type CommandJob struct {
	GenericJob
	Command *Command          `json:"-"`
	Target  *ProvisioningStep `json:"-"`
}

// CreateCommandJob is a constructor that is used by the planner to create a CommandJob for a given provisioning step.
func CreateCommandJob(id string, offset int, m *Metadata, pstep *ProvisioningStep) (*CommandJob, error) {
	j := &CommandJob{
		Target: pstep,
	}
	j.Metadata = m
	j.MetadataID = m.GetID()
	j.Offset = offset
	j.JobID = id
	j.Command = j.Target.Command
	if j.Target.Command.Timeout != 0 {
		j.Timeout = j.Target.Command.Timeout
	}
	j.JobType = "command_job"
	j.CreatedAt = time.Now()
	return j, nil
}

// CanProceed makes sure we can proceed and all of our dependencies are met
func (j *CommandJob) CanProceed(e chan error) {
	//Let's make sure we have a command to run
	if j.Command == nil {
		e <- errors.New("Command job had no command to run")
		return
	}
	// We need to have a set of targets
	if j.Target == nil {
		e <- errors.New("Command job had no targets")
		return
	}
	// We need to make sure we have an active connection
	if j.Target.ProvisionedHost.Conn.Active {
		e <- nil
		return
	}

	// We need to get our connection file for info on this host
	pathToConnFile := filepath.Join(j.Base.BaseDir, j.Target.ParentLaforgeID(), "conn.laforge")

	// Output from this command will be in our log file, let's make sure the log dir exists before we proceed
	logdir := filepath.Join(j.Base.BaseDir, j.Target.ParentLaforgeID(), "logs")
	if _, err := os.Stat(logdir); err != nil {
		if os.IsNotExist(err) {
			//nolint:gosec
			mkdirErr := os.MkdirAll(logdir, 0755)
			if mkdirErr != nil {
				e <- mkdirErr
				return
			}
		} else {
			cli.Logger.Errorf("Error creating log directory %s: %v", logdir, err)
			e <- err
			return
		}
	}

	// We need to make sure we have a connection file to proceed
	if _, err := os.Stat(pathToConnFile); err != nil {
		if os.IsNotExist(err) {
			e <- NewTimeoutExtension(fmt.Errorf("cannot proceed with a host that has no connection definition: %s", pathToConnFile))
			return
		}
		e <- nil
		return
	}

	// Now we build a connection struct from the connection file to make sure it's valid
	conn := &Connection{}
	err := LoadHCLFromFile(pathToConnFile, conn)
	if err != nil {
		cli.Logger.Errorf("Error loading job %s resource: %v", j.JobID, err)
		e <- err
		return
	}

	// We check to make sure it's an active connection
	if !conn.Active {
		e <- NewTimeoutExtension(errors.New("cannot proceed with a host with an inactive connection"))
		return
	}

	// Let's make sure our connection information is merged
	newConn, err := SmartMerge(j.Target.ProvisionedHost.Conn, conn, false)
	if err != nil {
		e <- fmt.Errorf("fatal error attempting to patch connection into state tree for %s: %v", j.JobID, err)
		return
	}

	// And that all of our connection data is good
	j.Target.ProvisionedHost.Conn = newConn.(*Connection)

	// Finally, let's actually test our connection over WinRM/SSH on the network to the system
	if !j.Target.ProvisionedHost.Conn.Test() {
		e <- NewTimeoutExtensionWithDelay(errors.New("Unable to successfuly make a test connection to host, retrying after a delay"), 20)
		return
	}

	e <- nil
}

// EnsureDependencies makes sure all of our dependencies (such as asset files, connection, etc. are working
func (j *CommandJob) EnsureDependencies(e chan error) {
	// Make sure we have a valid connection again
	if j.Target.ProvisionedHost.Conn == nil {
		e <- fmt.Errorf("command %s has a nil connection for the parent host", j.JobID)
		return
	}

	// If our connection is over SSH, we need to validate our key exists.  For Windows, we'll use credentials instead
	if j.Target.ProvisionedHost.Conn.IsSSH() {
		if j.Target.ProvisionedHost.Conn.SSHAuthConfig.IdentityFile == sshKeyPath {
			cli.Logger.Debugf("Fixing identity file for %s", j.Target.Path())
			j.Target.ProvisionedHost.Conn.SSHAuthConfig.IdentityFile = filepath.Join(j.Base.BaseDir, j.Base.CurrentBuild.Path(), "data", "ssh.pem")
		}
	}

	e <- nil
}

// Do is where we actually run the command
func (j *CommandJob) Do(e chan error) {
	// Let the user know what we're doing
	cli.Logger.Warnf("Performing Command Job:\n  %s %s: %s\n   %s   %s: %s", color.HiBlueString(">>"), color.HiCyanString(ObjectTypeCommand.String()), color.HiGreenString("%s", j.Command.CommandString()), color.HiBlueString(">>"), color.HiCyanString("HOST"), color.HiGreenString("%s", j.Target.ProvisionedHost.Conn.RemoteAddr))

	// Let's get the path to our logs
	logdir := filepath.Join(j.Base.BaseDir, j.Target.ParentLaforgeID(), "logs")
	logname := fmt.Sprintf("%d-%s", j.Target.StepNumber, filepath.Base(j.Command.ID))

	// Here we actually run the command
	err := j.Target.ProvisionedHost.Conn.ExecuteString(j, j.Command.CommandString(), logdir, logname)
	if err != nil {
		cli.Logger.Errorf("Error executing command %s: %s", j.JobID, err.Error())
		e <- err
		return
	}

	e <- nil
}

// CleanUp performs any cleanup we need to do afterward we can do it here, but we don't have any
func (j *CommandJob) CleanUp(e chan error) {
	cli.Logger.Debugf("Starting cleanup, cooldown running.")
	// Now we'll wait for the tooldown as defined in our command
	if j.Command.Cooldown > 0 {
		cli.Logger.Infof("Letting command job %s cooldown for %d seconds.", j.Command.ID, j.Command.Cooldown)
		time.Sleep(time.Duration(j.Command.Cooldown) * time.Second)
	}
	cli.Logger.Debugf("Finishing cleanup, cooldown done.")

	e <- nil
}

// Finish is called finally to let the log know we're done!
func (j *CommandJob) Finish(e chan error) {
	cli.Logger.Infof("Finished command %s", j.JobID)
	e <- nil
}
