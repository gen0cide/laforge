package laforge

import "path/filepath"

// CallFile is a debug type for tracking any file a configuration object was referenced in
type CallFile struct {
	CallerFile string
	CallerDir  string
}

// Caller represents a call chain in FIFO order of all a configuration object's CallFiles
type Caller []CallFile

// NewCaller returns a first generation Caller with a origin CallFile embedded inside
func NewCaller(src string) Caller {
	return Caller{
		CallFile{
			CallerFile: src,
			CallerDir:  filepath.Dir(src),
		},
	}
}

// Stack pushes a new Call stack (n) down onto the original Call stack (c)
func (c Caller) Stack(n Caller) Caller {
	return append(n, c...)
}

// Current retrieves the latest configuration file that has touched an object
func (c Caller) Current() CallFile {
	return c[0]
}
