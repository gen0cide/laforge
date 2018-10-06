package spanner

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gen0cide/laforge/core"
	"github.com/masterzen/winrm"
	"golang.org/x/crypto/ssh"
)

// Worker is a single unit of work and corrasponds to a single team element within a build
type Worker struct {
	Job
	ID         int64
	TeamID     int
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

	tlf, err := core.Bootstrap()
	if err != nil {
		return err
	}

	provisionedHost, found := tlf.Team.ActiveHosts[w.HostID]
	if !found || (provisionedHost != nil && provisionedHost.Active == false) {
		return fmt.Errorf("Host %s is currently not active in team %d environment", w.HostID, w.TeamID)
	}

	w.Host = provisionedHost
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

func (w *Worker) RunRemoteCommand(wc chan *Worker) {
	defer func() {
		wc <- w
	}()
	w.BeginTime = time.Now()
	defer func() {
		w.EndTime = time.Now()
	}()
	if w.Host.IsWinRM() {
		w.RunWinRMCommand(wc)
	} else {
		w.RunSSHCommand(wc)
	}
	return
}

func (w *Worker) RunWinRMCommand(wc chan *Worker) {
	endpoint := winrm.NewEndpoint(w.Host.RemoteAddr, w.Host.WinRMAuthConfig.Port, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, w.Host.WinRMAuthConfig.User, w.Host.WinRMAuthConfig.Password)
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
	if err != nil {
		w.ExitStatus = exitcode
		w.ExitError = err
		return
	}

	core.Logger.Infof("TEAM-%d/%s (output) >>>\n%s", w.TeamID, w.HostID, stdoutbuf.String())
	if stderrbuf.Len() > 0 {
		core.Logger.Errorf("TEAM-%d/%s (error) >>>\n%s", w.TeamID, w.HostID, stdoutbuf.String())
	}
	return
}

func (w *Worker) RunSSHCommand(wc chan *Worker) {
	// privateKey could be read from a file, or retrieved from another storage
	// source, such as the Secret Service / GNOME Keyring
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
		User: w.Host.SSHAuthConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", w.Host.SSHAuthConfig.Hostname, w.Host.SSHAuthConfig.Port), config)
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}
	// Create a session. It is one session per command.
	session, err := client.NewSession()
	if err != nil {
		w.ExitStatus = 1
		w.ExitError = err
		return
	}

	defer session.Close()

	var b bytes.Buffer  // import "bytes"
	session.Stdout = &b // get output
	// you can also pass what gets input to the stdin, allowing you to pipe
	// content from client to server
	//      session.Stdin = bytes.NewBufferString("My input")

	// Finally, run the command
	// err = session.Run(cmd)
	// return b.String(), err
	return
}
