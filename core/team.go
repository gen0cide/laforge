package core

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Team represents a team specific object existing within an environment
type Team struct {
	ID               string             `hcl:"id,label" json:"id,omitempty"`
	TeamNumber       int                `hcl:"team_number,attr" json:"team_number,omitempty"`
	BuildID          string             `hcl:"build_id,attr" json:"build_id,omitempty"`
	EnvironmentID    string             `hcl:"environment_id,attr" json:"environment_id,omitempty"`
	Config           map[string]string  `hcl:"config,attr" json:"config,omitempty"`
	Tags             map[string]string  `hcl:"tags,attr" json:"tags,omitempty"`
	ProvisionedHosts []*ProvisionedHost `hcl:"provisioned_host,block" json:"provisioned_hosts,omitempty"`
	Maintainer       *User              `hcl:"maintainer,block" json:"maintainer,omitempty"`
	Build            *Build             `json:"build,omitempty"`
	Environment      *Environment       `json:"environment,omitempty"`
	RelBuildPath     string             `json:"-"`
	Caller           Caller             `json:"-"`
}

// ProvisionedHost defines a provisioned host within a team's environment (network neutral)
type ProvisionedHost struct {
	ID              string           `hcl:"id,label" json:"id,omitempty"`
	RemoteAddr      string           `hcl:"remote_addr,attr" json:"remote_addr,omitempty"`
	SSHAuthConfig   *SSHAuthConfig   `hcl:"ssh_config,block" json:"ssh_config,omitempty"`
	WinRMAuthConfig *WinRMAuthConfig `hcl:"winrm_config,block" json:"winrm_config,omitempty"`
}

// SSHAuthConfig defines how Laforge should connect via SSH to a provisioned host
type SSHAuthConfig struct {
	Hostname        string        `hcl:"hostname,attr" json:"hostname,omitempty"`
	Port            int           `hcl:"port,attr" json:"port,omitempty"`
	User            string        `hcl:"user,attr" json:"user,omitempty"`
	Password        string        `hcl:"password,attr" json:"password,omitempty"`
	IdentityFile    string        `hcl:"identity_file,attr" json:"identity_file"`
	IdentityFileRef *LocalFileRef `json:"-"`
}

// WinRMAuthConfig defines how Laforge should connect via WinRM to a provisioned host
type WinRMAuthConfig struct {
	Hostname      string        `hcl:"ip,attr" json:"ip,omitempty"`
	Port          int           `hcl:"port,attr" json:"port,omitempty"`
	HTTPS         bool          `hcl:"https,attr" json:"https,omitempty"`
	SkipVerify    bool          `hcl:"skip_verify,attr" json:"skip_verify,omitempty"`
	TLSServerName string        `hcl:"tls_server_name,attr" json:"tls_server_name,omitempty"`
	CAFile        string        `hcl:"ca_file,attr" json:"ca_file,omitempty"`
	CAFileRef     *LocalFileRef `json:"-"`
	CertFile      string        `hcl:"cert_file,attr" json:"cert_file,omitempty"`
	CertFileRef   *LocalFileRef `json:"-"`
	KeyFile       string        `hcl:"key_file,attr" json:"key_file,omitempty"`
	KeyFileRef    *LocalFileRef `json:"-"`
	User          string        `hcl:"user,attr" json:"user,omitempty"`
	Password      string        `hcl:"password,attr" json:"password,omitempty"`
}

// IsSSH is a convenience method for checking if the provisioned host is setup for remote SSH
func (p *ProvisionedHost) IsSSH() bool {
	return p.SSHAuthConfig != nil
}

// IsWinRM is a convenience method for checking if the provisioned host is setup for remote WinRM
func (p *ProvisionedHost) IsWinRM() bool {
	return p.WinRMAuthConfig != nil
}

// LoadFileDeps attempts ot load important key material in the team configuration for connecting to remote team hosts
func (t *Team) LoadFileDeps(base *Laforge, pr *PathResolver, caller CallFile) error {
	for _, ph := range t.ProvisionedHosts {
		if ph.SSHAuthConfig != nil {
			err := ph.SSHAuthConfig.LoadIdentityFile(base, pr, caller)
			if err != nil {
				return errors.Wrapf(errors.WithStack(err), "could not load ssh identity_file for host %s team %s", ph.ID, t.ID)
			}
		}
		if ph.WinRMAuthConfig != nil {
			err := ph.WinRMAuthConfig.LoadCAFile(base, pr, caller)
			if err != nil {
				return errors.Wrapf(errors.WithStack(err), "could not load winrm ca_file for host %s team %s", ph.ID, t.ID)
			}
			err = ph.WinRMAuthConfig.LoadCertFile(base, pr, caller)
			if err != nil {
				return errors.Wrapf(errors.WithStack(err), "could not load winrm cert_file for host %s team %s", ph.ID, t.ID)
			}
			err = ph.WinRMAuthConfig.LoadKeyFile(base, pr, caller)
			if err != nil {
				return errors.Wrapf(errors.WithStack(err), "could not load winrm key_file for host %s team %s", ph.ID, t.ID)
			}
		}
	}
	return nil
}

// LoadIdentityFile attempts to locate the referenced source file with a laforge base configuration
func (s *SSHAuthConfig) LoadIdentityFile(base *Laforge, pr *PathResolver, caller CallFile) error {
	if s.IdentityFile == "" {
		return nil
	}
	cwd, _ := os.Getwd()
	testSrc := s.IdentityFile
	if !filepath.IsAbs(s.IdentityFile) {
		testSrc = filepath.Join(caller.CallerDir, s.IdentityFile)
	}
	if !PathExists(testSrc) {
		pr.Unresolved[s.IdentityFile] = true
		return errors.Wrapf(ErrAbsPathDeclNotExist, "caller=%s path=%s", caller.CallerFile, s.IdentityFile)
	}
	rel, _ := filepath.Rel(cwd, testSrc)
	rel2, _ := filepath.Rel(caller.CallerDir, testSrc)
	lfr := &LocalFileRef{
		Base:          filepath.Base(testSrc),
		AbsPath:       testSrc,
		RelPath:       rel,
		Cwd:           cwd,
		DeclaredPath:  s.IdentityFile,
		RelToCallFile: rel2,
	}
	s.IdentityFileRef = lfr
	return nil
}

// LoadCAFile attempts to locate the referenced source file with a laforge base configuration
func (w *WinRMAuthConfig) LoadCAFile(base *Laforge, pr *PathResolver, caller CallFile) error {
	if w.CAFile == "" {
		return nil
	}
	cwd, _ := os.Getwd()
	testSrc := w.CAFile
	if !filepath.IsAbs(w.CAFile) {
		testSrc = filepath.Join(caller.CallerDir, w.CAFile)
	}
	if !PathExists(testSrc) {
		pr.Unresolved[w.CAFile] = true
		return errors.Wrapf(ErrAbsPathDeclNotExist, "caller=%s path=%s", caller.CallerFile, w.CAFile)
	}
	rel, _ := filepath.Rel(cwd, testSrc)
	rel2, _ := filepath.Rel(caller.CallerDir, testSrc)
	lfr := &LocalFileRef{
		Base:          filepath.Base(testSrc),
		AbsPath:       testSrc,
		RelPath:       rel,
		Cwd:           cwd,
		DeclaredPath:  w.CAFile,
		RelToCallFile: rel2,
	}
	w.CAFileRef = lfr
	return nil
}

// LoadCertFile attempts to locate the referenced source file with a laforge base configuration
func (w *WinRMAuthConfig) LoadCertFile(base *Laforge, pr *PathResolver, caller CallFile) error {
	if w.CertFile == "" {
		return nil
	}
	cwd, _ := os.Getwd()
	testSrc := w.CertFile
	if !filepath.IsAbs(w.CertFile) {
		testSrc = filepath.Join(caller.CallerDir, w.CertFile)
	}
	if !PathExists(testSrc) {
		pr.Unresolved[w.CertFile] = true
		return errors.Wrapf(ErrAbsPathDeclNotExist, "caller=%s path=%s", caller.CallerFile, w.CertFile)
	}
	rel, _ := filepath.Rel(cwd, testSrc)
	rel2, _ := filepath.Rel(caller.CallerDir, testSrc)
	lfr := &LocalFileRef{
		Base:          filepath.Base(testSrc),
		AbsPath:       testSrc,
		RelPath:       rel,
		Cwd:           cwd,
		DeclaredPath:  w.CertFile,
		RelToCallFile: rel2,
	}
	w.CertFileRef = lfr
	return nil
}

// LoadKeyFile attempts to locate the referenced source file with a laforge base configuration
func (w *WinRMAuthConfig) LoadKeyFile(base *Laforge, pr *PathResolver, caller CallFile) error {
	if w.KeyFile == "" {
		return nil
	}
	cwd, _ := os.Getwd()
	testSrc := w.KeyFile
	if !filepath.IsAbs(w.KeyFile) {
		testSrc = filepath.Join(caller.CallerDir, w.KeyFile)
	}
	if !PathExists(testSrc) {
		pr.Unresolved[w.KeyFile] = true
		return errors.Wrapf(ErrAbsPathDeclNotExist, "caller=%s path=%s", caller.CallerFile, w.KeyFile)
	}
	rel, _ := filepath.Rel(cwd, testSrc)
	rel2, _ := filepath.Rel(caller.CallerDir, testSrc)
	lfr := &LocalFileRef{
		Base:          filepath.Base(testSrc),
		AbsPath:       testSrc,
		RelPath:       rel,
		Cwd:           cwd,
		DeclaredPath:  w.KeyFile,
		RelToCallFile: rel2,
	}
	w.KeyFileRef = lfr
	return nil
}
