package core

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/sftp"

	"github.com/docker/docker/pkg/term"
	"github.com/pkg/errors"
	"github.com/shiena/ansicolor"
	"golang.org/x/crypto/ssh"
)

const (
	// DefaultShebang is added at the top of a SSH script file
	DefaultShebang = "#!/bin/bash\n"

	sshKeyPath = `../../data/ssh.pem`
)

var randLock sync.Mutex
var randShared *rand.Rand

// SSHClient represents the client that connects to a remote server via SSH
type SSHClient struct {
	connInfo *SSHAuthConfig
	client   *ssh.Client
	config   *sshConfig
	conn     net.Conn

	lock sync.Mutex
}

type sshConfig struct {
	config     *ssh.ClientConfig
	connection func() (net.Conn, error)
	Pty        bool
}

type fatalError struct {
	error
}

func (e fatalError) FatalError() error {
	return e.error
}

// SSHClientConfig attempts to create an x/ssh client configuration for connection authentication
func SSHClientConfig(sshconf *SSHAuthConfig, overrideKey string) (*ssh.ClientConfig, error) {
	realKeyPath := ""
	if _, err := os.Stat(sshconf.IdentityFile); err != nil && os.IsNotExist(err) {
		if sshconf.IdentityFileRef == nil && overrideKey == "" {
			return nil, errors.New("could not locate SSH private key for authentication")
		}
		if sshconf.IdentityFileRef == nil {
			realKeyPath = "/dev/null"
		}
		if sshconf.IdentityFileRef != nil {
			if _, err := os.Stat(sshconf.IdentityFileRef.AbsPath); err == nil {
				realKeyPath = sshconf.IdentityFileRef.AbsPath
			}
		}
	} else {
		realKeyPath = sshconf.IdentityFile
	}

	keys := []ssh.Signer{}
	if realKeyPath != "" && realKeyPath != "/dev/null" {
		buf, err := ioutil.ReadFile(realKeyPath)
		if err != nil {
			return nil, err
		}
		aKey, err := ssh.ParsePrivateKey(buf)
		if err != nil {
			return nil, err
		}
		keys = append(keys, aKey)
	}

	if overrideKey != "" {
		aKey, err := ssh.ParsePrivateKey([]byte(overrideKey))
		if err != nil {
			return nil, err
		}
		keys = append(keys, aKey)
	}

	if len(keys) == 0 {
		return nil, errors.New("no public keys were available")
	}
	pubkeys := ssh.PublicKeys(keys...)

	config := &ssh.ClientConfig{
		User: sshconf.User,
		Auth: []ssh.AuthMethod{
			pubkeys,
		},
		//nolint:gosec
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config, nil
}

// NewSSHClient creates a new communicator implementation over SSH.
func NewSSHClient(sshconf *SSHAuthConfig, overrideKey string) (*SSHClient, error) {

	config, err := SSHClientConfig(sshconf, overrideKey)
	if err != nil {
		return nil, err
	}

	randLock.Lock()
	defer randLock.Unlock()
	if randShared == nil {
		randShared = rand.New(rand.NewSource(
			time.Now().UnixNano() * int64(os.Getpid())))
	}

	comm := &SSHClient{
		connInfo: sshconf,
		config: &sshConfig{
			config:     config,
			connection: ConnectFunc("tcp", fmt.Sprintf("%s:%d", sshconf.RemoteAddr, sshconf.Port)),
		},
	}

	return comm, nil
}

// Connect implementation of communicator.Communicator interface
func (s *SSHClient) Connect() (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.conn != nil {
		err := s.conn.Close()
		if err != nil {
			return err
		}
	}

	s.conn = nil
	s.client = nil

	s.conn, err = s.config.connection()
	if err != nil {
		s.conn = nil
		return err
	}

	host := fmt.Sprintf("%s:%d", s.connInfo.RemoteAddr, s.connInfo.Port)
	sshConn, sshChan, req, err := ssh.NewClientConn(s.conn, host, s.config.config)
	if err != nil {
		return err
	}

	s.client = ssh.NewClient(sshConn, sshChan, req)
	return err
}

// Disconnect implementation of communicator.Communicator interface
func (s *SSHClient) Disconnect() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.conn != nil {
		conn := s.conn
		s.conn = nil
		return conn.Close()
	}

	return nil
}

// LaunchInteractiveShell launches an interactive SSH session through the terminal
func (s *SSHClient) LaunchInteractiveShell() error {
	termWidth, termHeight := 80, 24
	session, err := s.newSession()
	if err != nil {
		return err
	}

	//nolint:gosec,errcheck
	defer session.Close()

	session.Stdout = ansicolor.NewAnsiColorWriter(os.Stdout)
	session.Stderr = ansicolor.NewAnsiColorWriter(os.Stderr)
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.ECHOCTL:       1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	fd := os.Stdin.Fd()

	if term.IsTerminal(fd) {
		oldState, err := term.MakeRaw(fd)
		if err != nil {
			return err
		}

		//nolint:gosec,errcheck
		defer term.RestoreTerminal(fd, oldState)

		winsize, err := term.GetWinsize(fd)
		if err == nil {
			termWidth = int(winsize.Width)
			termHeight = int(winsize.Height)
		}
	}

	if err := session.RequestPty("xterm", termHeight, termWidth, modes); err != nil {

		return err
	}

	if err := session.Shell(); err != nil {
		return err
	}

	go monWinCh(session, os.Stdout.Fd())

	return session.Wait()
}

// Start executes a remote command on the host
func (s *SSHClient) Start(cmd *RemoteCommand) error {
	session, err := s.newSession()
	if err != nil {
		return err
	}

	session.Stdin = cmd.Stdin
	session.Stdout = cmd.Stdout
	session.Stderr = cmd.Stderr

	if s.config.Pty {
		termModes := ssh.TerminalModes{
			ssh.ECHO:          0,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}

		if err := session.RequestPty("xterm", 80, 40, termModes); err != nil {
			return err
		}
	}

	err = session.Start(strings.TrimSpace(cmd.Command) + "\n")
	if err != nil {
		return err
	}

	go func() {
		//nolint:gosec,errcheck
		defer session.Close()

		err := session.Wait()
		exitStatus := 0
		if err != nil {
			exitErr, ok := err.(*ssh.ExitError)
			if ok {
				exitStatus = exitErr.ExitStatus()
			}
		}

		cmd.SetExitStatus(exitStatus, err)
	}()

	return nil
}

// UploadScriptV2 uses the 3rd party pkg/sftp Go package to upload instead of native x/ssh with scp modes.
func (s *SSHClient) UploadScriptV2(src, dst string) error {
	sftp, err := sftp.NewClient(s.client)
	if err != nil {
		return err
	}

	//nolint:gosec,errcheck
	defer sftp.Close()

	f, err := sftp.Create(dst)
	if err != nil {
		return err
	}

	//nolint:gosec
	fileInput, err := os.Open(src)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, fileInput)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	err = sftp.Chmod(dst, 0777)
	if err != nil {
		return err
	}

	return nil
}

// UploadFileV2 uses the 3rd party pkg/sftp Go package to upload instead of native x/ssh with scp modes.
func (s *SSHClient) UploadFileV2(src, dst string) error {
	sftp, err := sftp.NewClient(s.client)
	if err != nil {
		return err
	}

	//nolint:gosec,errcheck
	defer sftp.Close()

	f, err := sftp.Create(dst)
	if err != nil {
		return err
	}

	//nolint:gosec
	fileInput, err := os.Open(src)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, fileInput)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

// DeleteScriptV2 uses the 3rd party pkg/sftp Go package to securely erase a file
func (s *SSHClient) DeleteScriptV2(remotefile string) error {
	sftp, err := sftp.NewClient(s.client)
	if err != nil {
		return err
	}

	//nolint:gosec,errcheck
	defer sftp.Close()

	fi, err := sftp.Lstat(remotefile)
	if err != nil {
		return err
	}

	file, err := sftp.OpenFile(remotefile, os.O_RDWR)
	if err != nil {
		return err
	}

	size := fi.Size()
	zerobytes := make([]byte, size)

	copy(zerobytes, "0")

	_, err = file.Write(zerobytes)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	err = sftp.Remove(remotefile)
	if err != nil {
		return err
	}

	return nil
}

// Upload implementation of communicator.Communicator interface
func (s *SSHClient) Upload(path string, input io.Reader) error {
	targetDir := filepath.Dir(path)
	targetFile := filepath.Base(path)
	targetDir = filepath.ToSlash(targetDir)

	size := int64(0)

	switch src := input.(type) {
	case *os.File:
		fi, err := src.Stat()
		if err != nil {
			size = fi.Size()
		}
	case *bytes.Buffer:
		size = int64(src.Len())
	case *bytes.Reader:
		size = int64(src.Len())
	case *strings.Reader:
		size = int64(src.Len())
	}

	scpFunc := func(w io.Writer, stdoutR *bufio.Reader) error {
		return scpUploadFile(targetFile, input, w, stdoutR, size)
	}

	return s.scpSession("scp -vt "+targetDir, scpFunc)
}

// UploadScript uploads a script for execution to the remote host
func (s *SSHClient) UploadScript(path string, input io.Reader) error {
	reader := bufio.NewReader(input)
	prefix, err := reader.Peek(2)
	if err != nil {
		return fmt.Errorf("Error reading script: %s", err)
	}

	var script bytes.Buffer
	if string(prefix) != "#!" {
		//nolint:gosec,errcheck
		script.WriteString(DefaultShebang)
	}

	//nolint:gosec,errcheck
	script.ReadFrom(reader)
	if err := s.Upload(path, &script); err != nil {
		return err
	}

	var stdout, stderr bytes.Buffer
	cmd := &RemoteCommand{
		Command: fmt.Sprintf("chmod 0777 %s", path),
		Stdout:  &stdout,
		Stderr:  &stderr,
	}
	if err := s.Start(cmd); err != nil {
		return fmt.Errorf(
			"Error chmodding script file to 0777 in remote "+
				"machine: %s", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf(
			"Error chmodding script file to 0777 in remote "+
				"machine %v: %s %s", err, stdout.String(), stderr.String())
	}

	return nil
}

// UploadDir uploads a directory to the remote host
func (s *SSHClient) UploadDir(dst string, src string) error {
	scpFunc := func(w io.Writer, r *bufio.Reader) error {
		uploadEntries := func() error {
			//nolint:gosec
			f, err := os.Open(src)
			if err != nil {
				return err
			}

			//nolint:gosec,errcheck
			defer f.Close()

			entries, err := f.Readdir(-1)
			if err != nil {
				return err
			}

			return scpUploadDir(src, entries, w, r)
		}

		if src[len(src)-1] != '/' {
			return scpUploadDirProtocol(filepath.Base(src), w, r, uploadEntries)
		}
		return uploadEntries()
	}

	return s.scpSession("scp -rvt "+dst, scpFunc)
}

func (s *SSHClient) newSession() (session *ssh.Session, err error) {
	if s.client == nil {
		err = errors.New("ssh client is not connected")
	} else {
		session, err = s.client.NewSession()
	}

	if err != nil {
		if err := s.Connect(); err != nil {
			return nil, err
		}

		return s.client.NewSession()
	}

	return session, nil
}

func (s *SSHClient) scpSession(scpCommand string, f func(io.Writer, *bufio.Reader) error) error {
	session, err := s.newSession()
	if err != nil {
		return err
	}
	//nolint:gosec,errcheck
	defer session.Close()

	stdinW, err := session.StdinPipe()
	if err != nil {
		return err
	}

	defer func() {
		if stdinW != nil {
			//nolint:gosec,errcheck
			stdinW.Close()
		}
	}()

	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return err
	}
	stdoutR := bufio.NewReader(stdoutPipe)

	stderr := new(bytes.Buffer)
	session.Stderr = stderr

	if err := session.Start(scpCommand); err != nil {
		return err
	}

	if err := f(stdinW, stdoutR); err != nil && err != io.EOF {
		return err
	}

	//nolint:gosec,errcheck
	stdinW.Close()
	stdinW = nil
	err = session.Wait()
	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			if exitErr.ExitStatus() == 127 {
				return errors.New(
					"SCP failed to start. This usually means that SCP is not\n" +
						"properly installed on the remote system.",
				)
			}
		}

		return err
	}

	scpErr := stderr.String()
	if len(scpErr) > 0 {
		return fmt.Errorf("scp stderr: %q", stderr)
	}

	return nil
}

// checkSCPStatus checks that a prior command sent to SCP completed successfully
func checkSCPStatus(r *bufio.Reader) error {
	code, err := r.ReadByte()
	if err != nil {
		return err
	}

	if code != 0 {
		message, _, err := r.ReadLine()
		if err != nil {
			return fmt.Errorf("Error reading error message: %s", err)
		}
		return errors.New(string(message))
	}

	return nil
}

func scpUploadFile(dst string, src io.Reader, w io.Writer, r *bufio.Reader, size int64) error {
	if size == 0 {
		tf, err := ioutil.TempFile("", "temporary-file")
		if err != nil {
			return fmt.Errorf("Error creating temporary file for upload: %s", err)
		}
		//nolint:gosec,errcheck
		defer os.Remove(tf.Name())
		//nolint:gosec,errcheck
		defer tf.Close()

		if _, err := io.Copy(tf, src); err != nil {
			return err
		}

		if err := tf.Sync(); err != nil {
			return fmt.Errorf("Error creating temporary file for upload: %s", err)
		}

		if _, err := tf.Seek(0, 0); err != nil {
			return fmt.Errorf("Error creating temporary file for upload: %s", err)
		}

		fi, err := tf.Stat()
		if err != nil {
			return fmt.Errorf("Error creating temporary file for upload: %s", err)
		}

		src = tf
		size = fi.Size()
	}

	fmt.Fprintln(w, "C0644", size, dst)
	if err := checkSCPStatus(r); err != nil {
		return err
	}

	if _, err := io.CopyN(w, src, size); err != nil {
		return err
	}

	fmt.Fprint(w, "\x00")
	if err := checkSCPStatus(r); err != nil {
		return err
	}

	return nil
}

func scpUploadDirProtocol(name string, w io.Writer, r *bufio.Reader, f func() error) error {
	fmt.Fprintln(w, "D0755 0", name)
	err := checkSCPStatus(r)
	if err != nil {
		return err
	}

	if err := f(); err != nil {
		return err
	}

	_, err = fmt.Fprintln(w, "E")
	if err != nil {
		return err
	}

	return nil
}

func scpUploadDir(root string, fs []os.FileInfo, w io.Writer, r *bufio.Reader) error {
	for _, fi := range fs {
		realPath := filepath.Join(root, fi.Name())

		isSymlinkToDir := false
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			symPath, err := filepath.EvalSymlinks(realPath)
			if err != nil {
				return err
			}

			symFi, err := os.Lstat(symPath)
			if err != nil {
				return err
			}
			isSymlinkToDir = symFi.IsDir()
		}

		if !fi.IsDir() && !isSymlinkToDir {
			//nolint:gosec
			f, err := os.Open(realPath)
			if err != nil {
				return err
			}

			err = func() error {
				//nolint:gosec,errcheck
				defer f.Close()
				return scpUploadFile(fi.Name(), f, w, r, fi.Size())
			}()

			if err != nil {
				return err
			}

			continue
		}

		err := scpUploadDirProtocol(fi.Name(), w, r, func() error {
			//nolint:gosec
			f, err := os.Open(realPath)
			if err != nil {
				return err
			}

			//nolint:gosec,errcheck
			defer f.Close()

			entries, err := f.Readdir(-1)
			if err != nil {
				return err
			}
			return scpUploadDir(realPath, entries, w, r)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// ConnectFunc is a convenience method for returning a function
// that just uses net.Dial to communicate with the remote end that
// is suitable for use with the SSH communicator configuration.
func ConnectFunc(network, addr string) func() (net.Conn, error) {
	return func() (net.Conn, error) {
		c, err := net.DialTimeout(network, addr, 15*time.Second)
		if err != nil {
			return nil, err
		}

		if tcpConn, ok := c.(*net.TCPConn); ok {
			err := tcpConn.SetKeepAlive(true)
			if err != nil {
				return nil, err
			}
		}

		return c, nil
	}
}
