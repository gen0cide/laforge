package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/hashicorp/hcl2/ext/include"
	"github.com/hashicorp/hcl2/ext/transform"
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	hcl2 "github.com/hashicorp/hcl2/hcl"
	hcl2parse "github.com/hashicorp/hcl2/hclparse"
	zglob "github.com/mattn/go-zglob"
	"github.com/xlab/treeprint"
)

// Loader defines the Laforge configuration loader object
type Loader struct {
	// Parser is the actual HCLv2 parser
	Parser *hcl2parse.Parser

	// SourceFile is the location of the first file loaded
	SourceFile string

	// ConfigMap contains all the configuration steps
	ConfigMap map[string]*Laforge

	// CallerMap contains a reference of what files call what other files
	CallerMap map[string]Caller

	// Includes is a map of the dependency graph
	Includes treeprint.Tree

	// FileTree is used to map various parts of the call field
	FileTree map[string]treeprint.Tree
}

// AddToTree effectively tracks the filetree as it grows from dependencies for the Loader
func (l *Loader) AddToTree(filename, parentname string) treeprint.Tree {
	if t, ok := l.FileTree[filename]; ok {
		return t
	}
	parent, found := l.FileTree[parentname]
	var child treeprint.Tree
	if !found {
		switch filepath.Base(parentname) {
		case "team.laforge":
			parent = l.Includes.AddMetaNode(cli.Boldc("TEAM"), cli.Boldw(parentname))
		case "build.laforge":
			parent = l.Includes.AddMetaNode(cli.Boldg("BUILD"), cli.Boldw(parentname))
		case envFile:
			parent = l.Includes.AddMetaNode(cli.Boldy("ENV"), cli.Boldw(parentname))
		case "base.laforge":
			parent = l.Includes.AddMetaNode(cli.Boldr("BASE"), cli.Boldw(parentname))
		case "global.laforge":
			parent = l.Includes.AddMetaNode(cli.Boldb("GLOBAL"), cli.Boldw(parentname))
		case ".":
			l.Includes.SetValue(cli.Boldw("."))
			parent = l.Includes
		default:
			parent = l.Includes.AddNode(filename)
		}
		l.FileTree[parentname] = parent
	}
	switch filepath.Base(filename) {
	case "team.laforge":
		child = parent.AddMetaBranch(cli.Boldc("TEAM"), cli.Boldw(filename))
	case "build.laforge":
		child = parent.AddMetaBranch(cli.Boldg("BUILD"), cli.Boldw(filename))
	case envFile:
		child = parent.AddMetaBranch(cli.Boldy("ENV"), cli.Boldw(filename))
	case "base.laforge":
		child = parent.AddMetaBranch(cli.Boldr("BASE"), cli.Boldw(filename))
	case "global.laforge":
		child = parent.AddMetaBranch(cli.Boldb("GLOBAL"), cli.Boldw(filename))
	default:
		child = parent.AddNode(filename)
	}
	l.FileTree[filename] = child
	return child
}

// ParseConfigFile loads a root file into Loader
func (l *Loader) ParseConfigFile(filename string) error {
	var diags hcl2.Diagnostics
	if strings.HasSuffix(filename, ".json") {
		_, diags = l.Parser.ParseJSONFile(filename)
	} else {
		_, diags = l.Parser.ParseHCLFile(filename)
	}
	if diags.HasErrors() {
		for _, e := range diags.Errs() {
			ne, ok := e.(*hcl2.Diagnostic)
			if ok {
				cli.Logger.Errorf("Laforge failed to parse a config file:\n Location: %v\n    Issue: %v\n   Detail: %v", ne.Subject, ne.Summary, ne.Detail)
			}
		}
		return diags
	}
	l.AddToTree(filename, ".")
	l.SourceFile = filename
	l.CallerMap[filename] = NewCaller(filename)
	return nil
}

// NewLoader returns a default Loader type
func NewLoader() *Loader {
	return &Loader{
		Parser:    hcl2parse.NewParser(),
		ConfigMap: map[string]*Laforge{},
		CallerMap: map[string]Caller{},
		Includes:  treeprint.New(),
		FileTree:  map[string]treeprint.Tree{},
	}
}

// FileGlobResolver is a modified FileResolver in the HCLv2 include extension that accounts for globbed
// includes:
//	include {
//		path = "./foo/*.laforge"
//	}
func FileGlobResolver(baseDir string, parser *hcl2parse.Parser, loader *Loader) include.Resolver {
	return &fileGlobResolver{
		BaseDir: baseDir,
		Parser:  parser,
		Loader:  loader,
	}
}

type fileGlobResolver struct {
	BaseDir string
	Parser  *hcl2parse.Parser
	Loader  *Loader
}

func (r fileGlobResolver) ResolveBodyPath(path string, refRange hcl2.Range) (hcl2.Body, hcl2.Diagnostics) {
	callerFile := filepath.Join(refRange.Filename)
	callerDir := filepath.Dir(callerFile)
	targetFile := filepath.Join(callerDir, path)
	body := hcl2.EmptyBody()
	var diags hcl2.Diagnostics
	if strings.Contains(targetFile, `*`) {
		matches, err := zglob.Glob(targetFile)
		if err != nil {
			return body, hcl2.Diagnostics{&hcl2.Diagnostic{
				Severity: hcl2.DiagError,
				Summary:  "directory glob error",
				Detail:   fmt.Sprintf("could not glob on %s: %v", targetFile, err),
			}}
		}
		for _, m := range matches {
			switch {
			case strings.HasSuffix(m, ".json"):
				r.Loader.AddToTree(m, refRange.Filename)
				r.Loader.CallerMap[m] = NewCaller(m)
				_, newDiags := r.Parser.ParseJSONFile(m)
				diags = diags.Extend(newDiags)
			case strings.HasSuffix(m, ".laforge"):
				r.Loader.AddToTree(m, refRange.Filename)
				r.Loader.CallerMap[m] = NewCaller(m)
				_, newDiags := r.Parser.ParseHCLFile(m)
				diags = diags.Extend(newDiags)
			default:
				newDiag := &hcl2.Diagnostic{
					Severity: hcl2.DiagWarning,
					Summary:  "invalid file in glob",
					Detail:   fmt.Sprintf("%s is not a valid JSON or Laforge file (glob=%s)", m, targetFile),
				}
				diags = diags.Append(newDiag)
			}
		}
	} else {
		if strings.HasSuffix(targetFile, ".json") {
			_, diags = r.Parser.ParseJSONFile(targetFile)
		} else {
			_, diags = r.Parser.ParseHCLFile(targetFile)
		}
		r.Loader.AddToTree(targetFile, refRange.Filename)
		r.Loader.CallerMap[targetFile] = NewCaller(targetFile)
		if diags.HasErrors() {
			for _, e := range diags.Errs() {
				ne, ok := e.(*hcl2.Diagnostic)
				if ok {
					cli.Logger.Errorf("Laforge failed to parse a config file:\n Location: %v\n    Issue: %v\n   Detail: %v", ne.Subject, ne.Summary, ne.Detail)
				}
			}
		}
		return nil, diags
	}
	if diags.HasErrors() {
		for _, e := range diags.Errs() {
			ne, ok := e.(*hcl2.Diagnostic)
			if ok {
				cli.Logger.Errorf("Laforge failed to parse a config file:\n Location: %v\n    Issue: %v\n   Detail: %v", ne.Subject, ne.Summary, ne.Detail)
			}
		}
	}
	return body, diags
}

// Deconflict attempts to perform a state differential on all referenced files by
// traversing the config files in LIFO order.
func (l *Loader) Deconflict(filenames []string) (*Laforge, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	lf := &Laforge{
		CurrDir: cwd,
		Caller:  l.CallerMap[l.SourceFile],
	}
	lf.CreateIndex()
	for _, fname := range filenames {
		lf, err = Mask(lf, l.ConfigMap[fname])
		if err != nil {
			return lf, err
		}
		lf.Includes = append(lf.Includes, fname)
		cli.Logger.Debugf("Config Imported From: %s", fname)
	}
	lf.DependencyGraph = l.Includes
	return lf, nil
}

// NewMergeConflict is a helper function that creates a nicely formatted error
// when a merge conflict fails during an object differential update.
func NewMergeConflict(
	base, layer interface{},
	baseid, layerid string,
	baseCaller, layerCaller CallFile,
) error {
	return fmt.Errorf(
		"merge conflict between %T.%s (%s) and %T.%s (%s)",
		base,
		baseid,
		baseCaller.CallerFile,
		layer,
		layerid,
		layerCaller.CallerFile,
	)
}

// Bind enumerates the Loader's original file, performing recursive include loads to the
// Loader, generating ASTs for each dependency. Bind finishes with a call to Deconflict().
func (l *Loader) Bind() (*Laforge, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	root, err := filepath.Abs(cwd)
	if err != nil {
		return nil, err
	}
	transformer := include.Transformer("include", nil, FileGlobResolver(root, l.Parser, l))
	filenames := []string{}
	for name := range l.Parser.Files() {
		filenames = append([]string{name}, filenames...)
	}
	currLen := len(l.Parser.Files())
	for {
		for name, f := range l.Parser.Files() {
			transform.Deep(f.Body, transformer)
			exists := false
			for _, i := range filenames {
				if i == name {
					exists = true
				}
			}
			if !exists {
				filenames = append([]string{name}, filenames...)
			}
			newLF := &Laforge{}
			diags := gohcl2.DecodeBody(f.Body, nil, newLF)
			if diags.HasErrors() {
				for _, e := range diags.Errs() {
					ne, ok := e.(*hcl2.Diagnostic)
					if ok {
						cli.Logger.Errorf("Laforge failed to parse a config file:\n Location: %v\n    Issue: %v\n   Detail: %v", ne.Subject, ne.Summary, ne.Detail)
					}
				}
				return nil, diags
			}
			newLF.Filename = name
			newLF.Caller = l.CallerMap[name]
			l.ConfigMap[name] = newLF
		}
		newLen := len(l.Parser.Files())
		if currLen == newLen {
			break
		}
		currLen = newLen
	}
	return l.Deconflict(filenames)
}

type transientContext struct {
	Build              *Build               `hcl:"build,block" json:"build,omitempty"`
	Competition        *Competition         `hcl:"competition,block" json:"competition,omitempty"`
	Command            *Command             `hcl:"command,block" json:"command,omitempty"`
	DNSRecord          *DNSRecord           `hcl:"dns_record,block" json:"dns_record,omitempty"`
	Environment        *Environment         `hcl:"environment,block" json:"environment,omitempty"`
	Host               *Host                `cty:"host" hcl:"host,block" json:"host,omitempty"`
	Identity           *Identity            `hcl:"identity,block" json:"identity,omitempty"`
	Network            *Network             `hcl:"network,block" json:"network,omitempty"`
	RemoteFile         *RemoteFile          `hcl:"remote_file,block" json:"remote_file,omitempty"`
	Script             *Script              `hcl:"script,block" json:"script,omitempty"`
	Team               *Team                `hcl:"team,block" json:"team,omitempty"`
	User               *User                `hcl:"user,block" json:"user,omitempty"`
	AMI                *AMI                 `hcl:"ami,block" json:"ami,omitempty"`
	ProvisionedHost    *ProvisionedHost     `hcl:"provisioned_host,block" json:"provisioned_host,omitempty"`
	ProvisionedNetwork *ProvisionedNetwork  `hcl:"provisioned_network,block" json:"provisioned_network,omitempty"`
	ProvisioningStep   *ProvisioningStep    `hcl:"provisioning_step,block" json:"provisioning_step,omitempty"`
	Connection         *Connection          `hcl:"connection,block" json:"connection,omitempty"`
	Includes           []*transientIncludes `hcl:"include,block" json:"includes,omitempty"`
}

type transientIncludes struct {
	Path string `hcl:"path,attr" json:"path,omitempty"`
}

type transientReverseContext struct {
	Build           []*Build           `hcl:"build,block" json:"build,omitempty"`
	Competition     []*Competition     `hcl:"competition,block" json:"competition,omitempty"`
	Command         []*Command         `hcl:"command,block" json:"command,omitempty"`
	DNSRecord       []*DNSRecord       `hcl:"dns_record,block" json:"dns_record,omitempty"`
	Environment     []*Environment     `hcl:"environment,block" json:"environment,omitempty"`
	Host            []*Host            `cty:"host" hcl:"host,block" json:"host,omitempty"`
	Identity        []*Identity        `hcl:"identity,block" json:"identity,omitempty"`
	Network         []*Network         `hcl:"network,block" json:"network,omitempty"`
	RemoteFile      []*RemoteFile      `hcl:"remote_file,block" json:"remote_file,omitempty"`
	Script          []*Script          `hcl:"script,block" json:"script,omitempty"`
	Team            []*Team            `hcl:"team,block" json:"team,omitempty"`
	User            []*User            `hcl:"user,block" json:"user,omitempty"`
	ProvisionedHost []*ProvisionedHost `hcl:"provisioned_host,block" json:"provisioned_host,omitempty"`
	AMI             []*AMI             `hcl:"ami,block" json:"ami,omitempty"`
}

// GetEmptyObjByName returns a pointer to an initialized, but empty object of the specified type (camel case).
func GetEmptyObjByName(s string) (interface{}, error) {
	switch strings.ToLower(s) {
	case ObjectTypeBuild.String():
		return &Build{}, nil
	case ObjectTypeCompetition.String():
		return &Competition{}, nil
	case ObjectTypeCommand.String():
		return &Command{}, nil
	case ObjectTypeDNSRecord.String():
		return &DNSRecord{}, nil
	case ObjectTypeEnvironment.String():
		return &Environment{}, nil
	case ObjectTypeHost.String():
		return &Host{}, nil
	case ObjectTypeIdentity.String():
		return &Identity{}, nil
	case ObjectTypeNetwork.String():
		return &Network{}, nil
	case ObjectTypeRemoteFile.String():
		return &RemoteFile{}, nil
	case ObjectTypeScript.String():
		return &Script{}, nil
	case ObjectTypeTeam.String():
		return &Team{}, nil
	case ObjectTypeUser.String():
		return &User{}, nil
	case ObjectTypeAMI.String():
		return &AMI{}, nil
	case ObjectTypeProvisionedHost.String():
		return &ProvisionedHost{}, nil
	case ObjectTypeProvisionedNetwork.String():
		return &ProvisionedNetwork{}, nil
	case ObjectTypeProvisioningStep.String():
		return &ProvisioningStep{}, nil
	case ObjectTypeConnection.String():
		return &Connection{}, nil
	default:
		return nil, errors.New("specified core type name was not valid")
	}
}

func (t *transientReverseContext) Add(x interface{}) error {
	into := reflect.ValueOf(x)
	intoRaw := into.Elem()
	cont := reflect.ValueOf(t)
	cont = cont.Elem()

	aField := cont.FieldByName(intoRaw.Type().Name())
	aField.Set(reflect.Append(aField, into))
	return nil
}

// LoadHCLFromFile attempts to load and parse an HCL file with a filepath location
func LoadHCLFromFile(fileloc string, v interface{}) error {
	if _, err := os.Stat(fileloc); err != nil {
		return err
	}

	//nolint:gosec
	data, err := ioutil.ReadFile(fileloc)
	if err != nil {
		return err
	}

	err = HCLBytesToObject(data, v)
	if err != nil {
		return err
	}
	return nil
}

// HCLBytesToObject renders bytes into an object of a specific type
func HCLBytesToObject(data []byte, v interface{}) error {
	loader := NewLoader()
	into := reflect.ValueOf(v)
	intoRaw := into.Elem()
	file, diags := loader.Parser.ParseHCL(data, fmt.Sprintf("parse-%s-in-memory", intoRaw.Type().Name()))
	if diags.HasErrors() {
		return diags
	}
	if file == nil {
		return errors.New("hcl file AST returned nil from parsing")
	}
	tc := &transientContext{}
	diags = gohcl2.DecodeBody(file.Body, nil, tc)
	if diags.HasErrors() {
		return diags
	}

	cont := reflect.ValueOf(tc)
	cont = cont.Elem()

	newVal := cont.FieldByName(intoRaw.Type().Name())
	if newVal.IsNil() {
		return fmt.Errorf("a field of type %s was expected but not found in the data", intoRaw.Type().Name())
	}

	newValData := newVal.Elem()
	intoRaw.Set(newValData)
	return nil
}

func (t *transientContext) Add(v interface{}) error {
	into := reflect.ValueOf(v)
	intoRaw := into.Elem()
	cont := reflect.ValueOf(t)
	cont = cont.Elem()

	aField := cont.FieldByName(intoRaw.Type().Name())
	aField.Set(into)
	return nil
}
