package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cespare/xxhash"
	"github.com/packer-community/winrmcp/winrmcp"

	"github.com/pkg/errors"
)

// ShellConfig is a generic type for shell configurations
type ShellConfig interface {
	// Protocol denotes the protocol to be used
	Protocol() string
}

var (
	// ErrInvalidShellConfigType is thrown when an invalid shellconfig type is passed into a connection handler
	ErrInvalidShellConfigType = errors.New("invalid shell configuration provided to connection handler")
)

// SSHAuthConfig defines how Laforge should connect via SSH to a provisioned host
//easyjson:json
type SSHAuthConfig struct {
	RemoteAddr      string        `hcl:"remote_addr,attr" json:"remote_addr,omitempty"`
	Port            int           `hcl:"port,attr" json:"port,omitempty"`
	User            string        `hcl:"user,attr" json:"user,omitempty"`
	Password        string        `hcl:"password,optional" json:"password,omitempty"`
	IdentityFile    string        `hcl:"identity_file,optional" json:"identity_file,omitempty"`
	IdentityFileRef *LocalFileRef `json:"-"`
}

// WinRMAuthConfig defines how Laforge should connect via WinRM to a provisioned host
//easyjson:json
type WinRMAuthConfig struct {
	RemoteAddr    string        `hcl:"remote_addr,attr" json:"remote_addr,omitempty"`
	Port          int           `hcl:"port,attr" json:"port,omitempty"`
	HTTPS         bool          `hcl:"https,attr" json:"https,omitempty"`
	SkipVerify    bool          `hcl:"skip_verify,attr" json:"skip_verify,omitempty"`
	TLSServerName string        `hcl:"tls_server_name,optional" json:"tls_server_name,omitempty"`
	CAFile        string        `hcl:"ca_file,optional" json:"ca_file,omitempty"`
	CertFile      string        `hcl:"cert_file,optional" json:"cert_file,omitempty"`
	KeyFile       string        `hcl:"key_file,optional" json:"key_file,omitempty"`
	User          string        `hcl:"user,attr" json:"user,omitempty"`
	Password      string        `hcl:"password,attr" json:"password,omitempty"`
	KeyFileRef    *LocalFileRef `json:"-"`
	CertFileRef   *LocalFileRef `json:"-"`
	CAFileRef     *LocalFileRef `json:"-"`
}

// LoadFileDeps attempts ot load important key material in the team configuration for connecting to remote team hosts
func (t *Team) LoadFileDeps(base *Laforge, pr *PathResolver, caller CallFile) error {
	for _, ph := range t.ProvisionedHosts {
		if ph.Conn == nil {
			continue
		}
		if ph.Conn.SSHAuthConfig != nil {
			err := ph.Conn.SSHAuthConfig.LoadIdentityFile(base, pr, caller)
			if err != nil {
				return errors.Wrapf(errors.WithStack(err), "could not load ssh identity_file for host %s team %s", ph.ID, t.ID)
			}
		}
		if ph.Conn.WinRMAuthConfig != nil {
			err := ph.Conn.WinRMAuthConfig.LoadCAFile(base, pr, caller)
			if err != nil {
				return errors.Wrapf(errors.WithStack(err), "could not load winrm ca_file for host %s team %s", ph.ID, t.ID)
			}
			err = ph.Conn.WinRMAuthConfig.LoadCertFile(base, pr, caller)
			if err != nil {
				return errors.Wrapf(errors.WithStack(err), "could not load winrm cert_file for host %s team %s", ph.ID, t.ID)
			}
			err = ph.Conn.WinRMAuthConfig.LoadKeyFile(base, pr, caller)
			if err != nil {
				return errors.Wrapf(errors.WithStack(err), "could not load winrm key_file for host %s team %s", ph.ID, t.ID)
			}
		}
	}
	return nil
}

// Protocol implements the ShellConfig interface
func (s *SSHAuthConfig) Protocol() string {
	return "ssh"
}

// Protocol implements the ShellConfig interface
func (w *WinRMAuthConfig) Protocol() string {
	return "winrm"
}

// Hash implements the Hasher interface
func (s *SSHAuthConfig) Hash() uint64 {
	return xxhash.Sum64String(
		fmt.Sprintf(
			"rma=%v port=%v user=%v pw=%v ifile=%v",
			s.RemoteAddr,
			s.Port,
			s.User,
			s.Password,
			s.IdentityFile,
		),
	)
}

// Hash implements the Hasher interface
func (w *WinRMAuthConfig) Hash() uint64 {
	return xxhash.Sum64String(
		fmt.Sprintf(
			"rma=%v port=%v https=%v sv=%v tlsn=%v caf=%v cef=%v kef=%v user=%v pw=%v",
			w.RemoteAddr,
			w.Port,
			w.HTTPS,
			w.SkipVerify,
			w.TLSServerName,
			w.CAFile,
			w.CertFile,
			w.KeyFile,
			w.User,
			w.Password,
		),
	)
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

// ToUploadConfig returns the socket and a winrmcp config for uploading via WinRM
func (w *WinRMAuthConfig) ToUploadConfig() (string, winrmcp.Config) {
	return fmt.Sprintf("%s:%d", w.RemoteAddr, w.Port), winrmcp.Config{
		Auth: winrmcp.Auth{
			User:     w.User,
			Password: w.Password,
		},
		Https:    w.HTTPS,
		Insecure: w.SkipVerify,
	}
}
