package templates

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/gen0cide/laforge/builder/buildutil"
	"github.com/gen0cide/laforge/core"
)

var (
	// ErrUnsupportedContextType is thrown when an unsupported type is attempted to be wired into a template context.
	ErrUnsupportedContextType = errors.New("unsupported type added to context")
)

// Context is a meta structure that wraps pointers to any type that needs to be accessed
// inside of a build template. This uinfied struct helps keep a single convention
// on how to get to and access data inside of build templates.
type Context struct {
	Local              *core.Local
	Build              *core.Build
	Competition        *core.Competition
	Command            *core.Command
	DNS                *core.DNS
	DNSRecord          *core.DNSRecord
	Environment        *core.Environment
	Host               *core.Host
	Identity           *core.Identity
	Network            *core.Network
	RemoteFile         *core.RemoteFile
	Script             *core.Script
	Team               *core.Team
	User               *core.User
	Remote             *core.Remote
	AMI                *core.AMI
	Laforge            *core.Laforge
	Metadata           *core.Metadata
	ProvisionedNetwork *core.ProvisionedNetwork
	ProvisionedHost    *core.ProvisionedHost
	ProvisioningStep   *core.ProvisioningStep
	Revision           *core.Revision
	Connection         *core.Connection
	Dict               Dict
}

// Dict is a temporary dictionary to be used for context
type Dict map[string]string

// NewContext takes a varadic list of objects to be embedded into the returned template context
func NewContext(i ...interface{}) (*Context, error) {
	c := &Context{Dict: map[string]string{}}
	err := c.Attach(i...)
	c.Localize()
	return c, err
}

// Clone duplicates the current contexts including pointer references to the dependent objects. Note it does not clone the Dict - that remains locally resident.
func (c *Context) Clone() *Context {
	newC := &Context{Dict: map[string]string{}}
	newC.Build = c.Build
	newC.Competition = c.Competition
	newC.Command = c.Command
	newC.DNS = c.DNS
	newC.DNSRecord = c.DNSRecord
	newC.Environment = c.Environment
	newC.Host = c.Host
	newC.Identity = c.Identity
	newC.Network = c.Network
	newC.RemoteFile = c.RemoteFile
	newC.Script = c.Script
	newC.Team = c.Team
	newC.User = c.User
	newC.Remote = c.Remote
	newC.AMI = c.AMI
	newC.Laforge = c.Laforge
	newC.Local = c.Local
	newC.Metadata = c.Metadata
	newC.ProvisionedNetwork = c.ProvisionedNetwork
	newC.ProvisionedHost = c.ProvisionedHost
	newC.ProvisioningStep = c.ProvisioningStep
	newC.Revision = c.Revision
	newC.Connection = c.Connection
	return newC
}

// Set attachs the val string to the key in the Context dictionary and returns the val for pipelining.
func (c *Context) Set(key, val string) string {
	c.Dict[key] = val
	return val
}

// Get returns a temp value set in the context's dictionary.
func (c *Context) Get(key string) string {
	return c.Dict[key]
}

// Localize is a function to set the context's localization
func (c *Context) Localize() {
	if c.Local != nil {
		return
	}
	locale := &core.Local{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
	c.Local = locale
	return
}

// Attach wires up all the core Laforge types into a cohesive Context bundle for template rendering
func (c *Context) Attach(i ...interface{}) error {
	for _, o := range i {
		switch v := o.(type) {
		case *core.Build:
			c.Build = v
		case *core.Competition:
			c.Competition = v
		case *core.Command:
			c.Command = v
		case *core.DNS:
			c.DNS = v
		case *core.DNSRecord:
			c.DNSRecord = v
		case *core.Environment:
			c.Environment = v
		case *core.Host:
			c.Host = v
		case *core.Identity:
			c.Identity = v
		case *core.Network:
			c.Network = v
		case *core.RemoteFile:
			c.RemoteFile = v
		case *core.Script:
			c.Script = v
		case *core.Team:
			c.Team = v
		case *core.User:
			c.User = v
		case *core.Remote:
			c.Remote = v
		case *core.AMI:
			c.AMI = v
		case *core.Laforge:
			c.Laforge = v
		case *core.Metadata:
			c.Metadata = v
		case *core.ProvisionedNetwork:
			c.ProvisionedNetwork = v
		case *core.ProvisionedHost:
			c.ProvisionedHost = v
		case *core.ProvisioningStep:
			c.ProvisioningStep = v
		case *core.Revision:
			c.Revision = v
		case *core.Connection:
			c.Connection = v
		default:
			return buildutil.Throw(ErrUnsupportedContextType, "Cannot associate an object of this type to a template context.", &buildutil.V{
				"type": fmt.Sprintf("%T", v),
			})
		}
	}
	return nil
}
