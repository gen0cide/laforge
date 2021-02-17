package loader

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/gen0cide/laforge/loader/include"
	hcl2 "github.com/hashicorp/hcl/v2"
	hcl2parse "github.com/hashicorp/hcl/v2/hclparse"
	zglob "github.com/mattn/go-zglob"
)

// FileGlobResolver is a modified FileResolver in the HCLv2 include extension that accounts for globbed
// includes:
//	include {
//		path = "./foo/*.laforge"
//	}
func FileGlobResolver(baseDir string, parser *hcl2parse.Parser) include.Resolver {
	return &fileGlobResolver{
		BaseDir: baseDir,
		Parser:  parser,
	}
}

type fileGlobResolver struct {
	BaseDir string
	Parser  *hcl2parse.Parser
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
				_, newDiags := r.Parser.ParseJSONFile(m)
				diags = diags.Extend(newDiags)
			case strings.HasSuffix(m, ".laforge"):
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
		if diags.HasErrors() {
			for _, e := range diags.Errs() {
				ne, ok := e.(*hcl2.Diagnostic)
				if ok {
					log.Fatalf("Laforge failed to parse a config file:\n Location: %v\n    Issue: %v\n   Detail: %v", ne.Subject, ne.Summary, ne.Detail)
				}
			}
		}
		return nil, diags
	}
	if diags.HasErrors() {
		for _, e := range diags.Errs() {
			ne, ok := e.(*hcl2.Diagnostic)
			if ok {
				log.Fatalf("Laforge failed to parse a config file:\n Location: %v\n    Issue: %v\n   Detail: %v", ne.Subject, ne.Summary, ne.Detail)
			}
		}
	}
	return body, diags
}

// Loader defines the Laforge configuration loader object
type Loader struct {
	// Parser is the actual HCLv2 parser
	Parser *hcl2parse.Parser

	// SourceFile is the location of the first file loaded
	SourceFile string
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
				log.Fatalf("Laforge failed to parse a config file:\n Location: %v\n    Issue: %v\n   Detail: %v", ne.Subject, ne.Summary, ne.Detail)
			}
		}
		return diags
	}
	l.SourceFile = filename
	return nil
}

// // Bind enumerates the Loader's original file, performing recursive include loads to the
// // Loader, generating ASTs for each dependency. Bind finishes with a call to Deconflict().
// func (l *Loader) Bind() (*bool, error) {
// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		return nil, err
// 	}
// 	root, err := filepath.Abs(cwd)
// 	if err != nil {
// 		return nil, err
// 	}
// 	transformer := include.Transformer("include", nil, FileGlobResolver(root, l.Parser))
// 	filenames := []string{}
// 	for name := range l.Parser.Files() {
// 		filenames = append([]string{name}, filenames...)
// 	}
// 	currLen := len(l.Parser.Files())
// 	for {
// 		for name, f := range l.Parser.Files() {
// 			transform.Deep(f.Body, transformer)
// 			exists := false
// 			for _, i := range filenames {
// 				if i == name {
// 					exists = true
// 				}
// 			}
// 			if !exists {
// 				filenames = append([]string{name}, filenames...)
// 			}
// 			newLF := &Laforge{}
// 			diags := gohcl2.DecodeBody(f.Body, nil, newLF)
// 			if diags.HasErrors() {
// 				for _, e := range diags.Errs() {
// 					ne, ok := e.(*hcl2.Diagnostic)
// 					if ok {
// 						log.Fatalf("Laforge failed to parse a config file:\n Location: %v\n    Issue: %v\n   Detail: %v", ne.Subject, ne.Summary, ne.Detail)
// 					}
// 				}
// 				return nil, diags
// 			}
// 			newLF.Filename = name
// 			newLF.Caller = l.CallerMap[name]
// 			l.ConfigMap[name] = newLF
// 		}
// 		newLen := len(l.Parser.Files())
// 		if currLen == newLen {
// 			break
// 		}
// 		currLen = newLen
// 	}
// 	return l.Deconflict(filenames)
// }

func main() {
	tloader := &Loader{
		Parser:     hcl2parse.NewParser(),
		SourceFile: "",
	}
	tloader.ParseConfigFile("/home/red/Documents/infra/envs/fred/env.laforge")
}
