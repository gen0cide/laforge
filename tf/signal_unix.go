// +build !windows

package tf

import (
	"os"
	"syscall"
)

var forwardSignals = []os.Signal{syscall.SIGTERM, syscall.SIGINT}
