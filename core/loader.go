package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl2/ext/include"
	"github.com/hashicorp/hcl2/ext/transform"
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	hcl2 "github.com/hashicorp/hcl2/hcl"
	hcl2parse "github.com/hashicorp/hcl2/hclparse"
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
			parent = l.Includes.AddMetaNode(boldc("TEAM"), boldw(parentname))
		case "build.laforge":
			parent = l.Includes.AddMetaNode(boldg("BUILD"), boldw(parentname))
		case "env.laforge":
			parent = l.Includes.AddMetaNode(boldy("ENV"), boldw(parentname))
		case "base.laforge":
			parent = l.Includes.AddMetaNode(boldr("BASE"), boldw(parentname))
		case "global.laforge":
			parent = l.Includes.AddMetaNode(boldb("GLOBAL"), boldw(parentname))
		case ".":
			l.Includes.SetValue(boldw("."))
			parent = l.Includes
		default:
			parent = l.Includes.AddNode(filename)
		}
		l.FileTree[parentname] = parent
	}
	switch filepath.Base(filename) {
	case "team.laforge":
		child = parent.AddMetaBranch(boldc("TEAM"), boldw(filename))
	case "build.laforge":
		child = parent.AddMetaBranch(boldg("BUILD"), boldw(filename))
	case "env.laforge":
		child = parent.AddMetaBranch(boldy("ENV"), boldw(filename))
	case "base.laforge":
		child = parent.AddMetaBranch(boldr("BASE"), boldw(filename))
	case "global.laforge":
		child = parent.AddMetaBranch(boldb("GLOBAL"), boldw(filename))
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
		matches, err := filepath.Glob(targetFile)
		if err != nil {
			return body, hcl2.Diagnostics{&hcl2.Diagnostic{
				Severity: hcl2.DiagError,
				Summary:  "directory glob error",
				Detail:   fmt.Sprintf("could not glob on %s: %v", targetFile, err),
			}}
		}
		for _, m := range matches {
			if strings.HasSuffix(m, ".json") {
				r.Loader.AddToTree(m, refRange.Filename)
				r.Loader.CallerMap[m] = NewCaller(m)
				_, newDiags := r.Parser.ParseJSONFile(m)
				diags = diags.Extend(newDiags)
			} else if strings.HasSuffix(m, ".laforge") {
				r.Loader.AddToTree(m, refRange.Filename)
				r.Loader.CallerMap[m] = NewCaller(m)
				_, newDiags := r.Parser.ParseHCLFile(m)
				diags = diags.Extend(newDiags)
			} else {
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
		return nil, diags
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
		Logger.Infof("Config Imported From: %s", fname)
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
			gohcl2.DecodeBody(f.Body, nil, newLF)
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
