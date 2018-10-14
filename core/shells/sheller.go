package shells

import (
	"io"

	"github.com/gen0cide/laforge/core"
)

// Sheller is an interface to define remote connections to hosts
type Sheller interface {
	// Kind denotes the type of configuration
	Kind() string

	// SetIO sets the I/O interfaces for the connection
	SetIO(stdout io.Writer, stderr io.Writer, stdin io.Reader) error

	// SetConfig takes the configuration type required
	SetConfig(sc core.ShellConfig) error

	// LaunchInteractiveShell spawns an interactive user session
	LaunchInteractiveShell() error
}
