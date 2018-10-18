// +build windows

package runner

import (
	"os"
)

var forwardSignals []os.Signal = []os.Signal{}
