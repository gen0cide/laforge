package shells

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/shiena/ansicolor"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/pkg/errors"

	"github.com/gen0cide/laforge/core"
	"golang.org/x/crypto/ssh"
)

// SSH is a type that implements SSH connections to hosts
type SSH struct {
	Config *core.SSHAuthConfig
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Kind implements the Sheller interface
func (s *SSH) Kind() string {
	return "ssh"
}

// SetIO implements the Sheller interface
func (s *SSH) SetIO(stdout io.Writer, stderr io.Writer, stdin io.Reader) error {
	s.Stdin = stdin
	s.Stdout = stdout
	s.Stderr = stderr
	return nil
}

// SetConfig implements the Sheller interface
func (s *SSH) SetConfig(sc core.ShellConfig) error {
	ac, ok := sc.(*core.SSHAuthConfig)
	if !ok {
		return core.ErrInvalidShellConfigType
	}
	s.Config = ac
	if ac == nil {
		return errors.New("nil auth config provised")
	}
	return nil
}

// LaunchInteractiveShell attempts to launch a full functional TTY shell through SSH
func (s *SSH) LaunchInteractiveShell() error {
	pubkey, err := PublicKeyFile(s.Config.IdentityFile)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: s.Config.User,
		Auth: []ssh.AuthMethod{
			pubkey,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	socket := fmt.Sprintf("%s:%d", s.Config.RemoteAddr, s.Config.Port)
	conn, err := ssh.Dial("tcp", socket, config)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("failed to connect to %s", socket))
	}
	defer conn.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return errors.WithMessage(err, "failed to acquire raw terminal file descriptor")
	}

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return errors.WithMessage(err, "failed to get terminal dimensions")
	}

	session, err := conn.NewSession()
	if err != nil {
		return errors.WithMessage(err, "failed to create new session")
	}
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

	if err := session.RequestPty("xterm", termHeight, termWidth, modes); err != nil {
		return errors.WithMessage(err, "failed to request pseudo terminal")
	}

	defer terminal.Restore(fd, oldState)

	if err := session.Shell(); err != nil {
		return errors.WithMessage(err, "failed to start remote shell")
	}

	return session.Wait()
}

// PublicKeyFile does stuff
func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}
