package core

import (
	"fmt"
	"os"
	"path"

	"github.com/cespare/xxhash"
	"github.com/packer-community/winrmcp/winrmcp"
	"github.com/shiena/ansicolor"

	"github.com/pkg/errors"
)

// Connection defines an access method provisioned host within a team's environment
//easyjson:json
type Connection struct {
	ID                 string              `hcl:"id,label" json:"id,omitempty"`
	Active             bool                `hcl:"active,attr" json:"active,omitempty"`
	LocalAddr          string              `hcl:"local_addr,attr" json:"local_addr,omitempty"`
	RemoteAddr         string              `hcl:"remote_addr,attr" json:"remote_addr,omitempty"`
	ResourceName       string              `hcl:"resource_name,attr" json:"resource_name,omitempty"`
	SSHAuthConfig      *SSHAuthConfig      `hcl:"ssh,block" json:"ssh,omitempty"`
	WinRMAuthConfig    *WinRMAuthConfig    `hcl:"winrm,block" json:"winrm,omitempty"`
	Revision           int64               `hcl:"revision,attr" json:"revision,omitempty"`
	OnConflict         *OnConflict         `hcl:"on_conflict,block" json:"on_conflict,omitempty"`
	Competition        *Competition        `json:"-"`
	Environment        *Environment        `json:"-"`
	Build              *Build              `json:"-"`
	Team               *Team               `json:"-"`
	Network            *Network            `json:"-"`
	Host               *Host               `json:"-"`
	ProvisionedHost    *ProvisionedHost    `json:"-"`
	ProvisionedNetwork *ProvisionedNetwork `json:"-"`
	Caller             Caller              `json:"-"`
}

// Hash implements the Hasher interface
func (c *Connection) Hash() uint64 {
	sshc := uint64(666)
	wrmc := uint64(666)
	if c.IsWinRM() {
		wrmc = c.WinRMAuthConfig.Hash()
	}
	if c.IsSSH() {
		sshc = c.SSHAuthConfig.Hash()
	}
	return xxhash.Sum64String(
		fmt.Sprintf(
			"id=%v localaddr=%v rmaddr=%v rname=%v sshc=%v wrmc=%v",
			c.ID,
			c.LocalAddr,
			c.RemoteAddr,
			c.ResourceName,
			sshc,
			wrmc,
		),
	)
}

// Path implements the Pather interface
func (c *Connection) Path() string {
	return c.ID
}

// Base implements the Pather interface
func (c *Connection) Base() string {
	return path.Base(c.ID)
}

// ValidatePath implements the Pather interface
func (c *Connection) ValidatePath() error {
	if err := ValidateGenericPath(c.Path()); err != nil {
		return err
	}
	return nil
}

// IsSSH is a convenience method for checking if the provisioned host is setup for remote SSH
func (c *Connection) IsSSH() bool {
	return c.SSHAuthConfig != nil
}

// ParentLaforgeID returns connections parent provisioned host
func (c *Connection) ParentLaforgeID() string {
	return path.Dir(c.Path())
}

// IsWinRM is a convenience method for checking if the provisioned host is setup for remote WinRM
func (c *Connection) IsWinRM() bool {
	return c.WinRMAuthConfig != nil
}

// GetCaller implements the Mergeable interface
func (c *Connection) GetCaller() Caller {
	return c.Caller
}

// LaforgeID implements the Mergeable interface
func (c *Connection) LaforgeID() string {
	return c.ID
}

// GetOnConflict implements the Mergeable interface
func (c *Connection) GetOnConflict() OnConflict {
	if c.OnConflict == nil {
		return OnConflict{
			Do: "default",
		}
	}
	return *c.OnConflict
}

// SetCaller implements the Mergeable interface
func (c *Connection) SetCaller(ca Caller) {
	c.Caller = ca
}

// SetOnConflict implements the Mergeable interface
func (c *Connection) SetOnConflict(o OnConflict) {
	c.OnConflict = &o
}

// Swap implements the Mergeable interface
func (c *Connection) Swap(m Mergeable) error {
	rawVal, ok := m.(*Connection)
	if !ok {
		return errors.Wrapf(ErrSwapTypeMismatch, "expected %T, got %T", c, m)
	}
	*c = *rawVal
	return nil
}

// SetID increments the revision and sets the ID if needed
func (c *Connection) SetID() string {
	if c.ID == "" {
		c.ID = path.Join(c.ProvisionedHost.LaforgeID(), "conn")
	}
	return c.ID
}

// RemoteShell connects your local console to a remote provisioned host
func (c *Connection) RemoteShell() error {
	if c.IsWinRM() {
		return c.InteractiveWinRM()
	}
	return c.InteractiveSSH()
}

// Upload uploads a src file/dir to a dst file/dir on the provisioned host
func (c *Connection) Upload(src, dst string) error {
	if c.IsWinRM() {
		return c.UploadWinRM(src, dst)
	}
	return c.UploadSCP(src, dst)
}

// UploadWinRM uses WinRM to upload src to dst on the provisioned host
func (c *Connection) UploadWinRM(src, dst string) error {
	addr, config := c.WinRMAuthConfig.ToUploadConfig()
	client, err := winrmcp.New(addr, &config)
	if err != nil {
		return err
	}
	return client.Copy(src, dst)
}

// InteractiveWinRM launches an interactive shell over WinRM
func (c *Connection) InteractiveWinRM() error {
	client := &WinRMClient{}
	err := client.SetConfig(c.WinRMAuthConfig)
	if err != nil {
		return err
	}
	client.SetIO(
		ansicolor.NewAnsiColorWriter(os.Stdout),
		ansicolor.NewAnsiColorWriter(os.Stderr),
		os.Stdin,
	)

	err = client.LaunchInteractiveShell()
	if err != nil {
		return err
	}
	return nil
}

// InteractiveSSH launches an interactive shell over SSH
func (c *Connection) InteractiveSSH() error {
	client, err := NewSSHClient(c.SSHAuthConfig, "")
	if err != nil {
		return err
	}

	err = client.Connect()
	if err != nil {
		return err
	}
	defer client.Disconnect()

	err = client.LaunchInteractiveShell()
	if err != nil {
		return err
	}
	return nil
}

// UploadSCP uses scp to upload src to dst on the provisioned host
func (c *Connection) UploadSCP(src, dst string) error {
	isDir := false
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		isDir = true
	}

	client, err := NewSSHClient(c.SSHAuthConfig, "")
	if err != nil {
		return err
	}

	err = client.Connect()
	if err != nil {
		return err
	}
	defer client.Disconnect()

	if isDir {
		err = client.UploadDir(dst, src)
		if err != nil {
			return err
		}
		return nil
	}

	fileInput, err := os.Open(src)
	if err != nil {
		return err
	}
	err = client.Upload(dst, fileInput)
	if err != nil {
		return err
	}
	return nil
}
