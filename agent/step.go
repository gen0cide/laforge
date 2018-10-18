package agent

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Step is a unique action to be taken on the local machine
type Step struct {
	ID          int                    `json:"id,omitempty"`
	Revision    string                 `json:"revision,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	StepType    string                 `json:"step_type,omitempty"`
	Source      string                 `json:"source,omitempty"`
	Destination string                 `json:"destination,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Status      string                 `json:"status,omitempty"`
	StartedAt   time.Time              `json:"started_at,omitempty"`
	EndedAt     time.Time              `json:"ended_at,omitempty"`
	StdoutFile  string                 `json:"stdout_file,omitempty"`
	StderrFile  string                 `json:"stderr_file,omitempty"`
	Stdout      *bytes.Buffer          `json:"-"`
	Stderr      *bytes.Buffer          `json:"-"`
	ExitStatus  int                    `json:"exit_status,omitempty"`
	ExitError   error                  `json:"exit_error,omitempty"`
}

// status = [ started, errored, awaiting_reboot, finished ]

// CanProceed returns whether the step is safe to proceed onto the next one
func (s *Step) CanProceed() bool {
	if s.Status == "finished" {
		return true
	}
	return false
}

// Prepare sets some values for the Step function to be performed
func (s *Step) Prepare() error {
	s.Status = "started"
	s.StartedAt = time.Now().UTC()
	s.StdoutFile = filepath.Join(StepLogDir(), s.StdoutFilename())
	s.StderrFile = filepath.Join(StepLogDir(), s.StderrFilename())
	s.Stdout = new(bytes.Buffer)
	s.Stderr = new(bytes.Buffer)
	return nil
}

// StdoutFilename returns the templated filename for a step's stdout log
func (s *Step) StdoutFilename() string {
	return fmt.Sprintf("%d_%s.stdout.log", s.ID, s.Name)
}

// StderrFilename returns the templated filename for a step's stderr log
func (s *Step) StderrFilename() string {
	return fmt.Sprintf("%d_%s.stderr.log", s.ID, s.Name)
}

// Perform actually performs the step's functions
func (s *Step) Perform() error {
	defer func() {
		s.EndedAt = time.Now().UTC()
	}()

	switch s.StepType {
	case "dns_record":
		// noop
		return nil
	case "remote_file":
		return s.CopyFile()
	case "script":
		if AsyncWorker.Config.Host.IsWindows() {
			switch s.Metadata["language"].(string) {
			case "powershell":
				return s.ExecutePowershell()
			default:
				return s.ExecuteWindows()
			}
		}
		return s.ExecuteLinux()
	}
	return nil
}

// ExecuteLinux makes the script executable and runs it in Linux
func (s *Step) ExecuteLinux() error {
	cmd := filepath.Join(AgentHomeDir, s.Source)
	err := os.Chmod(cmd, 0700)
	if err != nil {
		return err
	}
	return s.RunCommand(cmd)
}

// ExecuteWindows simply executes the source file
func (s *Step) ExecuteWindows() error {
	cmd := filepath.Clean(filepath.Join(AgentHomeDir, s.Source))
	return s.RunCommand(cmd)
}

// ExecutePowershell attempts to execute the script source as a powershell -File argument
func (s *Step) ExecutePowershell() error {
	cmd := filepath.Clean(filepath.Join(AgentHomeDir, s.Source))
	return s.RunCommand("powershell", "-NoLogo", "-NonInteractive", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", cmd)
}

// CopyFile attempts to move the target file from one place to the next
func (s *Step) CopyFile() error {
	dst := filepath.Clean(s.Destination)
	src := filepath.Clean(filepath.Join(AgentHomeDir, s.Source))
	return copyFile(src, dst)
}

// RunCommand is a generic executor for local commands in a way that preserves the STDERR and STDOUT
// by logging to a file as well as to a buffer which can be streamed.
func (s *Step) RunCommand(base string, args ...string) error {
	cmd := exec.Command(base, args...)

	stdoutfile, err := os.Create(s.StdoutFile)
	if err != nil {
		return err
	}
	defer stdoutfile.Close()
	stderrfile, err := os.Create(s.StderrFile)
	if err != nil {
		return err
	}
	defer stderrfile.Close()

	stdoutwriter := bufio.NewWriter(stdoutfile)
	defer stdoutwriter.Flush()

	stderrwriter := bufio.NewWriter(stderrfile)
	defer stderrwriter.Flush()

	stdoutpipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrpipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go io.Copy(io.MultiWriter(stdoutwriter, s.Stdout), stdoutpipe)
	go io.Copy(io.MultiWriter(stderrwriter, s.Stderr), stderrpipe)

	err = cmd.Start()
	if err != nil {
		s.ExitError = err
		s.ExitStatus = 1
		return err
	}

	err = cmd.Wait()
	if err != nil {
		s.ExitError = err
		s.ExitStatus = 1
		return err
	}
	return nil
}

// copyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func copyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
