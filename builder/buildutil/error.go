package buildutil

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// V is a short type for declaring vars
type V map[string]interface{}

// BuildIssue is specially crafted error to help in debugging
type BuildIssue struct {
	Source error
	Reason string
	Vars   V
}

// Error implements the error interface
func (b BuildIssue) Error() string {
	return b.String()
}

// String implements the stringer interface
func (b BuildIssue) String() string {
	lines := []string{}
	first := fmt.Sprintf("%s%s%s %s", color.HiYellowString("["), color.HiRedString("BUILD ISSUE"), color.HiYellowString("]"), color.HiWhiteString(b.Reason))
	lines = append(lines, first)
	second := fmt.Sprintf("  %s: %v", color.HiRedString("ERROR"), b.Source)
	lines = append(lines, second)
	if len(b.Vars) > 0 {
		header := fmt.Sprintf("  %s:", color.HiRedString("VARS"))
		lines = append(lines, header)
		for k, v := range b.Vars {
			lines = append(lines, fmt.Sprintf("    %s %s %v", color.HiWhiteString(k), color.HiYellowString("=>"), v))
		}
	}
	return strings.Join(lines, "\n")
}

// BuildError represents a chain of build errors
type BuildError []BuildIssue

// Error implements the error interface
func (b BuildError) Error() string {
	return b.String()
}

// String implements the stringer interface
func (b BuildError) String() string {
	lines := []string{}
	for _, x := range b {
		lines = append(lines, x.String())
	}
	return fmt.Sprintf("%s\n\n%s\n", color.HiYellowString("[ *** BUILD ERRORS *** ]"), strings.Join(lines, "\n\n"))
}

// IsBuildError is a convenience function to tell you if an Error is a BuildErro
func IsBuildError(e error) bool {
	_, ok := e.(BuildError)
	return ok
}

// Throw creates a build error with some additional information that will be passed to the user.
func Throw(err error, reason string, vars *V) error {
	if vars == nil {
		vars = &V{}
	}
	be := BuildIssue{
		Source: err,
		Reason: reason,
		Vars:   *vars,
	}
	return BuildError{be}
}

// Stack pushes a new Call stack (n) down onto the original Call stack (c)
func (b BuildError) Stack(addition BuildError) BuildError {
	return append(addition, b...)
}

// Top retrieves the the first element in the build error stack
func (b BuildError) Top() *BuildIssue {
	if len(b) == 0 {
		return nil
	}
	return &b[0]
}

// Bottom retrieves the the last element in the build error stack
func (b BuildError) Bottom() *BuildIssue {
	if len(b) == 0 {
		return nil
	}
	return &b[len(b)-1]
}
