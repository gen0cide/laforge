// Package builder defines the types and interfaces that are used to transpile laforge states
// into other infrastructure and automation tools.
package builder

import (
	"fmt"
	"os"

	"github.com/gen0cide/laforge/builder/tfgcp"

	"github.com/gen0cide/laforge/builder/buildutil/valdations"

	"github.com/fatih/color"
	"github.com/gen0cide/laforge/builder/buildutil"
	"github.com/gen0cide/laforge/builder/null"
	"github.com/gen0cide/laforge/core"
	"github.com/pkg/errors"
)

var (
	// ValidBuilders retains a map of ID to empty Builder objects.
	ValidBuilders = map[string]Builder{
		"tfgcp": tfgcp.New(),
		// "tfaws": tfaws.New(),
		// "tfibm": tfibm.New(),
		"null": null.New(),
	}

	// ErrNoBuilderFound is thrown when a builder is not a known valid builder parameter
	ErrNoBuilderFound = errors.New("builder type defined in environment was not recognized")
)

// BuildEngine is the primary interface for building components in gscript.
type BuildEngine struct {
	Base    *core.Laforge
	Builder Builder
}

// Steps
// CheckRequirements(l *core.Laforge) error

// Builder is a generic interface used to implement various builders.
type Builder interface {
	// ID of the builder - usually the go package name
	// - Must be unique and be a valid golang package name.
	ID() string

	// Name of the builder - usually titleized version of the type
	Name() string

	// Human readable description information
	Description() string

	// Author name and contact information
	Author() string

	// Version of the builder implementation
	Version() string

	// Validations returns a list of requirements that need to be met for a builder to be successful
	Validations() validations.Validations

	// Set's a Laforge base instance to work off
	SetLaforge(*core.Laforge) error

	// Check the provided base instance for any dynamic
	// configuration and ensure build requirements are met.
	// While it might seem redundant with Validations(), the point
	// of CheckRequirements is to be able to provide more
	// stateful analysis of the context instead of just true/false fields.
	CheckRequirements() error

	// Gather remote assets, upload to S3, create a unique names, etc.
	PrepareAssets() error

	// Run all the scripts through a rendering step, templating them out
	// if needed.
	GenerateScripts() error

	// Stage any dependencies (or preform pre-rendering tasks)
	StageDependencies() error

	// Attempt to render the actual configuration
	Render() error
}

// New attempts to create a new BuildEngine based on the laforge state parameters
func New(base *core.Laforge, overwrite, update bool) (*BuildEngine, error) {
	err := base.AssertMinContext(core.EnvContext)
	if err != nil {
		return nil, buildutil.Throw(err, "Cannot perform a build without being in an EnvContext.", nil)
	}

	setup := core.InitializeBuildDirectory(base, overwrite, update)
	if setup != nil {
		return nil, buildutil.Throw(setup, "Cannot initialize build directory", nil)
	}

	bldr, found := ValidBuilders[base.CurrentEnv.Builder]
	if !found {
		return nil, buildutil.Throw(ErrNoBuilderFound, "Invalid builder defined", &buildutil.V{"provided": base.CurrentEnv.Builder})
	}

	core.SetLogName(fmt.Sprintf("%s/%s", color.WhiteString("builder"), color.HiGreenString(base.CurrentEnv.Builder)))
	return &BuildEngine{
		Base:    base,
		Builder: bldr,
	}, nil
}

// Do performs the waterfall of functions on the given context with it's designated builder
func (b *BuildEngine) Do() error {
	err := b.Builder.SetLaforge(b.Base)
	if err != nil {
		return buildutil.Throw(err, "failed to set laforge", nil)
	}
	core.Logger.Infof("Injected context into builder")
	vals := b.Builder.Validations()
	for _, x := range vals {
		core.Logger.Debugf("Checking validation: %s", x.Name)
		if !x.Check(b.Base) {
			core.Logger.Errorf("Build Requirement Failed: %s", x.Name)
			core.Logger.Errorf("  Resolution > %s", x.Resolution)
			os.Exit(1)
		}
	}
	err = b.Builder.CheckRequirements()
	if err != nil {
		return buildutil.Throw(err, "failed checking requirements", nil)
	}
	core.Logger.Infof("Requirements checks passed")
	err = b.Builder.PrepareAssets()
	if err != nil {
		return buildutil.Throw(err, "failed preparing assets", nil)
	}
	core.Logger.Infof("Resolved and cached required assets")
	err = b.Builder.StageDependencies()
	if err != nil {
		return buildutil.Throw(err, "failed staging dependencies", nil)
	}
	core.Logger.Infof("Staged dependencies for rendering")
	err = b.Builder.GenerateScripts()
	if err != nil {
		return buildutil.Throw(err, "failed generating scripts", nil)
	}
	core.Logger.Infof("Generated scripts and templates")
	err = b.Builder.Render()
	if err != nil {
		return buildutil.Throw(err, "failed rendering build", nil)
	}
	core.Logger.Infof("Successfully rendered build")
	return nil
}
