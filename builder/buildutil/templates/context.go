package templates

import (
	"errors"
	"fmt"

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
	Build       *core.Build
	Competition *core.Competition
	Command     *core.Command
	DNS         *core.DNS
	DNSRecord   *core.DNSRecord
	Environment *core.Environment
	Host        *core.Host
	Identity    *core.Identity
	Network     *core.Network
	RemoteFile  *core.RemoteFile
	Script      *core.Script
	Team        *core.Team
	User        *core.User
	Remote      *core.Remote
	AMI         *core.AMI
	Laforge     *core.Laforge
}

// NewContext takes a varadic list of objects to be embedded into the returned template context
func NewContext(i ...interface{}) (*Context, error) {
	c := &Context{}
	err := c.Attach(i...)
	return c, err
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
		default:
			return buildutil.Throw(ErrUnsupportedContextType, "Cannot associate an object of this type to a template context.", &buildutil.V{
				"type": fmt.Sprintf("%T", v),
			})
		}
	}
	return nil
}
