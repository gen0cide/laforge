package core

import (
	"io"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/masterzen/winrm"
)

// WinRMClient is a type to connection to Windows hosts remotely over the WinRM protocol
type WinRMClient struct {
	Config *WinRMAuthConfig
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Kind implements the Sheller interface
func (w *WinRMClient) Kind() string {
	return "winrm"
}

// SetIO implements the Sheller interface
func (w *WinRMClient) SetIO(stdout io.Writer, stderr io.Writer, stdin io.Reader) error {
	w.Stdin = stdin
	w.Stdout = stdout
	w.Stderr = stderr
	return nil
}

// SetConfig implements the Sheller interface
func (w *WinRMClient) SetConfig(c *WinRMAuthConfig) error {
	if c == nil {
		return errors.New("nil auth config provised")
	}
	w.Config = c
	return nil
}

// LaunchInteractiveShell implements the Sheller interface
func (w *WinRMClient) LaunchInteractiveShell() error {
	endpoint := winrm.NewEndpoint(
		w.Config.RemoteAddr,
		w.Config.Port,
		w.Config.HTTPS,
		w.Config.SkipVerify,
		[]byte{},
		[]byte{},
		[]byte{},
		20*time.Second,
	)

	if w.Stderr == nil {
		w.Stderr = os.Stderr
	}

	if w.Stdin == nil {
		w.Stdin = os.Stdin
	}

	if w.Stdout == nil {
		w.Stdout = os.Stdout
	}

	client, err := winrm.NewClient(endpoint, w.Config.User, w.Config.Password)
	if err != nil {
		return errors.WithMessage(err, "could not create winrm client")
	}

	shell, err := client.CreateShell()
	if err != nil {
		panic(err)
	}
	var cmd *winrm.Command
	cmd, err = shell.Execute("powershell -NoProfile -ExecutionPolicy Bypass")
	if err != nil {
		panic(err)
	}

	go io.Copy(cmd.Stdin, os.Stdin)
	go io.Copy(os.Stdout, cmd.Stdout)
	go io.Copy(os.Stderr, cmd.Stderr)

	cmd.Wait()
	shell.Close()

	return nil
}
