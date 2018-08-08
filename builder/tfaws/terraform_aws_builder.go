package tfaws

import (
	"github.com/gen0cide/gscript/logger"
	"github.com/gen0cide/laforge/builder/buildutil"
	"github.com/gen0cide/laforge/core"
)

const (
	_id          = `tfaws`
	_name        = `Terraform AWS Builder`
	_description = `generates terraform configurations that isolate teams into VPCs`
	_author      = `Alex Levinson <github.com/gen0cide>`
	_version     = `0.0.1`
)

var (
	validations = buildutil.Validations{
		buildutil.Requirement{
			Name:       "Environment maintainer not defined",
			Resolution: "add a maintainer block to your environment configuration",
			Check:      buildutil.FieldNotEmpty(core.Environment{}, "Maintainer"),
		},
	}
)

// TerraformAWSBuilder implements a laforge builder that packages an environment into
// a terraform configuration targeting AWS with each team isolated into their own VPC.
type TerraformAWSBuilder struct {
	Base   *core.Laforge
	Logger logger.Logger
}

// New creates an empty TerraformAWSBuilder
func New() *TerraformAWSBuilder {
	return &TerraformAWSBuilder{}
}

// ID implements the Builder interface (returns the ID of the builder - usually the go package name)
func (t *TerraformAWSBuilder) ID() string {
	return _id
}

// Name implements the Builder interface (returns the name of the builder - usually titleized version of the type)
func (t *TerraformAWSBuilder) Name() string {
	return _name
}

// Description implements the Builder interface (returns the builder's description)
func (t *TerraformAWSBuilder) Description() string {
	return _description
}

// Author implements the Builder interface (author's name and contact info)
func (t *TerraformAWSBuilder) Author() string {
	return _author
}

// Version implements the Builder interface (builder version)
func (t *TerraformAWSBuilder) Version() string {
	return _version
}

// Validations implements the Builder interface (builder checks)
func (t *TerraformAWSBuilder) Validations() buildutil.Validations {
	return validations
}

// SetLaforge implements the Builder interface
func (t *TerraformAWSBuilder) SetLaforge(base *core.Laforge) error {
	t.Base = base
	return nil
}

// CheckRequirements implements the Builder interface
func (t *TerraformAWSBuilder) CheckRequirements() error {
	return nil
}

// PrepareAssets implements the Builder interface
func (t *TerraformAWSBuilder) PrepareAssets() error {
	return nil
}

// GenerateScripts implements the Builder interface
func (t *TerraformAWSBuilder) GenerateScripts() error {
	return nil
}

// StageDependencies implements the Builder interface
func (t *TerraformAWSBuilder) StageDependencies() error {
	return nil
}

// Render implements the Builder interface
func (t *TerraformAWSBuilder) Render() error {
	return nil
}
