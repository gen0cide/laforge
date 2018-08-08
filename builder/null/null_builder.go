package null

import (
	"github.com/gen0cide/gscript/logger"
	"github.com/gen0cide/laforge/core"
)

const (
	_id          = `null`
	_name        = `Null Builder`
	_description = `NOOP builder used for testing, debugging, and research.`
	_author      = `Alex Levinson <github.com/gen0cide>`
	_version     = `0.0.1`
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
	return _id
}

// Name implements the Builder interface (returns the name of the builder - usually titleized version of the type)
func (b *Builder) Name() string {
	return _name
}

// Description implements the Builder interface (returns the builder's description)
func (b *Builder) Description() string {
	return _description
}

// Author implements the Builder interface (author's name and contact info)
func (b *Builder) Author() string {
	return _author
}

// Version implements the Builder interface (builder version)
func (b *Builder) Version() string {
	return _version
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
