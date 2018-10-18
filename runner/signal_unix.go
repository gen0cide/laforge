// +build !windows

package runner

import (
	"os"
	"syscall"
)

var forwardSignals = []os.Signal{syscall.SIGTERM, syscall.SIGINT}
