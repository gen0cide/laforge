// Package null implements a Laforge Builder module that is effectively a NOOP. It should be used
// as the spec reference for building builder modules as it's probably got the easiest learning curve.
package null

import (
	"github.com/gen0cide/gscript/logger"
	"github.com/gen0cide/laforge/builder/buildutil/valdations"
	"github.com/gen0cide/laforge/core"
)

// Definition of builder meta-data.
const (
	ID          = `null`
	Name        = `Null Builder`
	Description = `NOOP builder used for testing, debugging, and research.`
	Author      = `Alex Levinson <github.com/gen0cide>`
	Version     = `0.0.1`
)

var (
	rules = validations.Validations{
		validations.Requirement{
			Name:       "environment maintainer not defined",
			Resolution: "add a maintainer block to your environment configuration",
			Check:      validations.FieldNotEmpty(core.Environment{}, "Maintainer"),
		},
	}
)

// Builder implements a laforge builder that packages an environment into
// a terraform configuration targeting AWS with each team isolated into their own VPC.
type Builder struct {
	Base   *core.Laforge
	Logger logger.Logger
}

// New creates an empty Builder
func New() *Builder {
	return &Builder{}
}

// ID implements the Builder interface (returns the ID of the builder - usually the go package name)
func (b *Builder) ID() string {
	return ID
}

// Name implements the Builder interface (returns the name of the builder - usually titleized version of the type)
func (b *Builder) Name() string {
	return Name
}

// Description implements the Builder interface (returns the builder's description)
func (b *Builder) Description() string {
	return Description
}

// Author implements the Builder interface (author's name and contact info)
func (b *Builder) Author() string {
	return Author
}

// Version implements the Builder interface (builder version)
func (b *Builder) Version() string {
	return Version
}

// Validations implements the Builder interface (builder checks)
func (b *Builder) Validations() validations.Validations {
	return rules
}

// SetLaforge implements the Builder interface
func (b *Builder) SetLaforge(base *core.Laforge) error {
	b.Base = base
	return nil
}

// CheckRequirements implements the Builder interface
func (b *Builder) CheckRequirements() error {
	return nil
}

// PrepareAssets implements the Builder interface
func (b *Builder) PrepareAssets() error {
	return nil
}

// GenerateScripts implements the Builder interface
func (b *Builder) GenerateScripts() error {
	return nil
}

// StageDependencies implements the Builder interface
func (b *Builder) StageDependencies() error {
	return nil
}

// Render implements the Builder interface
func (b *Builder) Render() error {
	return nil
}

// Set implements the Builder interface
func (b *Builder) Set(key string, val interface{}) {
	return
}

// Get implements the Builder interface
func (b *Builder) Get(key string) string {
	return ""
}
