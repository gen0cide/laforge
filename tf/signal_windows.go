// +build windows

package tf

import (
	"os"
)

var forwardSignals []os.Signal = []os.Signal{}
