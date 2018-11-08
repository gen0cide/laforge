package spanner

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gen0cide/laforge/core"
	"github.com/gen0cide/laforge/core/cli"
	"github.com/masterzen/winrm"
	"golang.org/x/crypto/ssh"
)

// Worker is a single unit of work and corrasponds to a single team element within a build
type Worker struct {
	Job
	ID         int64
	TeamNumber int
	TeamID     string
	Team       *core.Team
	State      *core.State
	Laforge    *core.Laforge
	Parent     *Spanner
	Host       *core.ProvisionedHost
	BeginTime  time.Time
	EndTime    time.Time
	LogFile    string
	TeamDir    string
	ExitStatus int
	ExitError  error
}

// ResolveProvisionedHost gets the provisioned host config from the team
func (w *Worker) ResolveProvisionedHost() error {
	err := os.Chdir(w.TeamDir)
	if err != nil {
		return err
	}

	state, err := core.BootstrapWithState(true)
	if err != nil {
		return err
	}
	if state == nil {
		return errors.New("cannot proceed with a nil state")
	}

	return nil
}

// RunLocalCommand executes the worker's task, piping output to the logfile
func (w *Worker) RunLocalCommand(wc chan *Worker) {
	defer func() {
		wc <- w
	}()
	w.BeginTime = time.Now()
	defer func() {
		w.EndTime = time.Now()
	}()
	cmd := exec.Command(w.Job.Command[0], w.Job.Command[1:]...)
	cmd.Dir = w.TeamDir
	cmd.Env = append(os.Environ())

	stdoutFile := fmt.Sprintf("%s.stdout", w.LogFile)
	stderrFile := fmt.Sprintf("%s.stderr", w.LogFile)

	stdoutfile, err := os.Create(stdoutFile)
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}
	defer stdoutfile.Close()

	stderrfile, err := os.Create(stderrFile)
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}
	defer stderrfile.Close()

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}

	stdoutwriter := bufio.NewWriter(stdoutfile)
	defer stdoutwriter.Flush()

	stderrwriter := bufio.NewWriter(stderrfile)
	defer stderrwriter.Flush()

	go io.Copy(stdoutwriter, stdoutPipe)
	go io.Copy(stderrwriter, stderrPipe)

	err = cmd.Start()
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}

	err = cmd.Wait()
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}
	return
}

// RunRemoteCommand begins the process for executing commands across the Spanner's worker.
func (w *Worker) RunRemoteCommand(wc chan *Worker) {
	defer func() {
		wc <- w
	}()
	w.BeginTime = time.Now()
	defer func() {
		w.EndTime = time.Now()
	}()
	if w.Host.Conn.IsWinRM() {
		w.RunWinRMCommand(wc)
	} else {
		w.RunSSHCommand(wc)
	}
	return
}

// RunWinRMCommand executes the remote command over the WinRM protocol on remote Windows hosts
func (w *Worker) RunWinRMCommand(wc chan *Worker) {
	endpoint := winrm.NewEndpoint(w.Host.Conn.RemoteAddr, w.Host.Conn.WinRMAuthConfig.Port, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, w.Host.Conn.WinRMAuthConfig.User, w.Host.Conn.WinRMAuthConfig.Password)
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}

	cmd := strings.Join(w.Command, " ")

	stdoutFile := fmt.Sprintf("%s.stdout", w.LogFile)
	stderrFile := fmt.Sprintf("%s.stderr", w.LogFile)

	stdoutfile, err := os.Create(stdoutFile)
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}
	defer stdoutfile.Close()

	stderrfile, err := os.Create(stderrFile)
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}
	defer stderrfile.Close()

	var stdoutbuf bytes.Buffer
	var stderrbuf bytes.Buffer

	stdoutwriter := bufio.NewWriter(stdoutfile)
	defer stdoutwriter.Flush()

	stderrwriter := bufio.NewWriter(stderrfile)
	defer stderrwriter.Flush()

	exitcode, err := client.Run(cmd, io.MultiWriter(&stdoutbuf, stdoutwriter), io.MultiWriter(&stderrbuf, stderrwriter))

	defer func() {
		cli.Logger.Infof("TEAM-%d/%s (output) >>>\n%s", w.TeamID, w.HostID, stdoutbuf.String())
		if stderrbuf.Len() > 0 {
			cli.Logger.Errorf("TEAM-%d/%s (error) >>>\n%s", w.TeamID, w.HostID, stdoutbuf.String())
		}
	}()

	if err != nil {
		w.ExitStatus = exitcode
		w.ExitError = err
		return
	}
	return
}

// RunSSHCommand executes the remote command over the SSH remote protocol
func (w *Worker) RunSSHCommand(wc chan *Worker) {
	keydata, err := ioutil.ReadFile(filepath.Join(w.Parent.BuildDir, "data", "ssh.pem"))
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}
	key, err := ssh.ParsePrivateKey(keydata)
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}

	config := &ssh.ClientConfig{
		User: w.Host.Conn.SSHAuthConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", w.Host.Conn.SSHAuthConfig.RemoteAddr, w.Host.Conn.SSHAuthConfig.Port), config)
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}

	session, err := client.NewSession()
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}

	defer session.Close()

	cmd := strings.Join(w.Command, " ")

	stdoutFile := fmt.Sprintf("%s.stdout", w.LogFile)
	stderrFile := fmt.Sprintf("%s.stderr", w.LogFile)

	stdoutfile, err := os.Create(stdoutFile)
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}
	defer stdoutfile.Close()

	stderrfile, err := os.Create(stderrFile)
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}
	defer stderrfile.Close()

	var stdoutbuf bytes.Buffer
	var stderrbuf bytes.Buffer

	stdoutwriter := bufio.NewWriter(stdoutfile)
	defer stdoutwriter.Flush()

	stderrwriter := bufio.NewWriter(stderrfile)
	defer stderrwriter.Flush()

	session.Stdout = io.MultiWriter(&stdoutbuf, stdoutwriter)
	session.Stderr = io.MultiWriter(&stderrbuf, stderrwriter)

	err = session.Run(cmd)

	defer func() {
		cli.Logger.Infof("TEAM-%d/%s (output) >>>\n%s", w.TeamID, w.HostID, stdoutbuf.String())
		if stderrbuf.Len() > 0 {
			cli.Logger.Errorf("TEAM-%d/%s (error) >>>\n%s", w.TeamID, w.HostID, stdoutbuf.String())
		}
	}()

	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}

	return
}

// Verify attempts to validate the constructs of the spanner
func (w *Worker) Verify() error {
	if w.ExecType == "remote-exec" {
		provisionedHostFile := filepath.Join(w.TeamDir, w.HostID, fmt.Sprintf("%s.laforge", "conn"))
		if _, err := os.Stat(provisionedHostFile); os.IsNotExist(err) {
			return fmt.Errorf("team %s does not have an active host %s", w.TeamID, w.HostID)
		}
		err := w.ResolveProvisionedHost()
		if err != nil {
			return err
		}
	}

	return nil
}
