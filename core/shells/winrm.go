package shells

import (
	"io"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/masterzen/winrm"

	"github.com/gen0cide/laforge/core"
)

type WinRM struct {
	Config *core.WinRMAuthConfig
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Kind implements the Sheller interface
func (w *WinRM) Kind() string {
	return "winrm"
}

// SetIO implements the Sheller interface
func (w *WinRM) SetIO(stdout io.Writer, stderr io.Writer, stdin io.Reader) error {
	w.Stdin = stdin
	w.Stdout = stdout
	w.Stderr = stderr
	return nil
}

// SetConfig implements the Sheller interface
func (w *WinRM) SetConfig(sc core.ShellConfig) error {
	ac, ok := sc.(*core.WinRMAuthConfig)
	if !ok {
		return core.ErrInvalidShellConfigType
	}
	if ac == nil {
		return errors.New("nil auth config provised")
	}
	w.Config = ac
	return nil
}

// LaunchInteractiveShell implements the Sheller interface
func (w *WinRM) LaunchInteractiveShell() error {
	endpoint := winrm.NewEndpoint(
		w.Config.Hostname,
		w.Config.Port,
		w.Config.HTTPS,
		w.Config.SkipVerify,
		[]byte{},
		[]byte{},
		[]byte{},
		20*time.Second,
	)

	client, err := winrm.NewClient(endpoint, w.Config.User, w.Config.Password)
	if err != nil {
		return errors.WithMessage(err, "could not create winrm client")
	}

	_, err = client.RunWithInput("powershell", os.Stdout, os.Stderr, os.Stdin)
	if err != nil {
		return errors.WithMessage(err, "connection issue")
	}

	return nil
}
