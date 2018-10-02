package core

// Local is used to represent information about the current runtime to the user
type Local struct {
	OS   string
	Arch string
}

// IsWindows is a template helper function
func (l *Local) IsWindows() bool {
	return l.OS == "windows"
}

// IsMacOS is a template helper function
func (l *Local) IsMacOS() bool {
	return l.OS == "darwin"
}

// IsLinux is a template helper function
func (l *Local) IsLinux() bool {
	return l.OS == "linux"
}
