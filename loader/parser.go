package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/competition"
	"github.com/gen0cide/laforge/ent/environment"
	"github.com/gen0cide/laforge/ent/filedelete"
	"github.com/gen0cide/laforge/ent/filedownload"
	"github.com/gen0cide/laforge/ent/fileextract"
	"github.com/gen0cide/laforge/ent/identity"
	"github.com/gen0cide/laforge/ent/user"
	"github.com/gen0cide/laforge/loader/include"
	hcl2 "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/transform"
	gohcl2 "github.com/hashicorp/hcl/v2/gohcl"
	hcl2parse "github.com/hashicorp/hcl/v2/hclparse"
	_ "github.com/mattn/go-sqlite3"
	zglob "github.com/mattn/go-zglob"
)

// Include defines a named include type
type Include struct {
	Path string `hcl:"path,attr"`
}

type fileGlobResolver struct {
	BaseDir string
	Parser  *hcl2parse.Parser
}

// DefinedConfigs is the stuct to hold in all the loading for hcl
type DefinedConfigs struct {
	Filename            string
	BaseDir             string                       `hcl:"base_dir,optional" json:"base_dir,omitempty"`
	IncludePaths        []*Include                   `hcl:"include,block" json:"include_paths,omitempty"`
	DefinedCompetitions []*ent.Competition           `hcl:"competition,block" json:"competitions,omitempty"`
	DefinedHosts        []*ent.Host                  `hcl:"host,block" json:"hosts,omitempty"`
	DefinedNetworks     []*ent.Network               `hcl:"network,block" json:"networks,omitempty"`
	DefinedScripts      []*ent.Script                `hcl:"script,block" json:"scripts,omitempty"`
	DefinedCommands     []*ent.Command               `hcl:"command,block" json:"defined_commands,omitempty"`
	DefinedDNSRecords   []*ent.DNSRecord             `hcl:"dns_record,block" json:"defined_dns_records,omitempty"`
	DefinedEnvironments []*ent.Environment           `hcl:"environment,block" json:"environments,omitempty"`
	DefinedFileDownload []*ent.FileDownload          `hcl:"file_download,block" json:"file_download,omitempty"`
	DefinedFileDelete   []*ent.FileDelete            `hcl:"file_delete,block" json:"file_delete,omitempty"`
	DefinedFileExtract  []*ent.FileExtract           `hcl:"file_extract,block" json:"file_extract,omitempty"`
	DefinedIdentities   []*ent.Identity              `hcl:"identity,block" json:"identities,omitempty"`
	Competitions        map[string]*ent.Competition  `json:"-"`
	Hosts               map[string]*ent.Host         `json:"-"`
	Networks            map[string]*ent.Network      `json:"-"`
	Scripts             map[string]*ent.Script       `json:"-"`
	Commands            map[string]*ent.Command      `json:"-"`
	DNSRecords          map[string]*ent.DNSRecord    `json:"-"`
	Environments        map[string]*ent.Environment  `json:"-"`
	FileDownload        map[string]*ent.FileDownload `json:"-"`
	FileDelete          map[string]*ent.FileDelete   `json:"-"`
	FileExtract         map[string]*ent.FileExtract  `json:"-"`
	Identities          map[string]*ent.Identity     `json:"-"`
}

// Loader defines the Laforge configuration loader object
type Loader struct {
	// Parser is the actual HCLv2 parser
	Parser *hcl2parse.Parser

	// SourceFile is the location of the first file loaded
	SourceFile string

	// ConfigMap contains all the configuration steps
	ConfigMap map[string]*DefinedConfigs
}

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

// Bind enumerates the Loader's original file, performing recursive include loads to the
// Loader, generating ASTs for each dependency. Bind finishes with a call to Deconflict().
func (l *Loader) Bind() (*DefinedConfigs, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	root, err := filepath.Abs(cwd)
	if err != nil {
		return nil, err
	}
	transformer := include.Transformer("include", nil, FileGlobResolver(root, l.Parser))
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
			newLF := &DefinedConfigs{}
			diags := gohcl2.DecodeBody(f.Body, nil, newLF)
			if diags.HasErrors() {
				for _, e := range diags.Errs() {
					ne, ok := e.(*hcl2.Diagnostic)
					if ok {
						log.Fatalf("Laforge failed to parse a config file:\n Location: %v\n    Issue: %v\n   Detail: %v", ne.Subject, ne.Summary, ne.Detail)
					}
				}
				return nil, diags
			}
			newLF.Filename = name
			l.ConfigMap[name] = newLF
		}
		newLen := len(l.Parser.Files())
		if currLen == newLen {
			break
		}
		currLen = newLen
	}
	return l.merger(filenames)
}

// NewLoader returns a default Loader type
func NewLoader() *Loader {
	return &Loader{
		Parser:     hcl2parse.NewParser(),
		ConfigMap:  map[string]*DefinedConfigs{},
		SourceFile: "",
	}
}

func (l *Loader) merger(filenames []string) (*DefinedConfigs, error) {
	combinedConfigs := &DefinedConfigs{
		Filename:     l.SourceFile,
		Competitions: map[string]*ent.Competition{},
		Hosts:        map[string]*ent.Host{},
		Networks:     map[string]*ent.Network{},
		Scripts:      map[string]*ent.Script{},
		Commands:     map[string]*ent.Command{},
		DNSRecords:   map[string]*ent.DNSRecord{},
		Environments: map[string]*ent.Environment{},
		FileDownload: map[string]*ent.FileDownload{},
		FileDelete:   map[string]*ent.FileDelete{},
		FileExtract:  map[string]*ent.FileExtract{},
		Identities:   map[string]*ent.Identity{},
	}
	for _, filename := range filenames {
		element := l.ConfigMap[filename]
		for _, x := range element.DefinedCompetitions {
			obj, found := combinedConfigs.Competitions[x.HclID]
			if !found {
				combinedConfigs.Competitions[x.HclID] = x
				continue
			}
			fmt.Println("Stored: ", obj)
			fmt.Println("New: ", x)
			if x.RootPassword != "" {
				obj.RootPassword = x.RootPassword
			}
			if x.Config != nil {
				obj.Config = x.Config
			}
			if x.Tags != nil {
				obj.Tags = x.Tags
			}
			if x.HCLCompetitionToDNS != nil {
				obj.HCLCompetitionToDNS = x.HCLCompetitionToDNS
			}
			combinedConfigs.Competitions[x.HclID] = obj
		}
		for _, x := range element.DefinedHosts {
			obj, found := combinedConfigs.Hosts[x.HclID]
			if !found {
				combinedConfigs.Hosts[x.HclID] = x
				continue
			}
			fmt.Println("Stored: ", obj)
			fmt.Println("New: ", x)
		}
		for _, x := range element.DefinedNetworks {
			obj, found := combinedConfigs.Networks[x.HclID]
			if !found {
				combinedConfigs.Networks[x.HclID] = x
				continue
			}
			fmt.Println("Stored: ", obj)
			fmt.Println("New: ", x)
		}
		for _, x := range element.DefinedScripts {
			obj, found := combinedConfigs.Scripts[x.HclID]
			if !found {
				combinedConfigs.Scripts[x.HclID] = x
				continue
			}
			fmt.Println("Stored: ", obj)
			fmt.Println("New: ", x)
		}
		for _, x := range element.DefinedCommands {
			obj, found := combinedConfigs.Commands[x.HclID]
			if !found {
				combinedConfigs.Commands[x.HclID] = x
				continue
			}
			fmt.Println("Stored: ", obj)
			fmt.Println("New: ", x)
		}
		for _, x := range element.DefinedDNSRecords {
			obj, found := combinedConfigs.DNSRecords[x.HclID]
			if !found {
				combinedConfigs.DNSRecords[x.HclID] = x
				continue
			}
			fmt.Println("Stored: ", obj)
			fmt.Println("New: ", x)
		}
		for _, x := range element.DefinedEnvironments {
			obj, found := combinedConfigs.Environments[x.HclID]
			if !found {
				combinedConfigs.Environments[x.HclID] = x
				continue
			}
			fmt.Println("Stored: ", obj)
			fmt.Println("New: ", x)
		}
		for _, x := range element.DefinedFileDownload {
			obj, found := combinedConfigs.FileDownload[x.HclID]
			if !found {
				combinedConfigs.FileDownload[x.HclID] = x
				continue
			}
			fmt.Println("Stored: ", obj)
			fmt.Println("New: ", x)
		}
		for _, x := range element.DefinedFileDelete {
			element.FileDelete[x.HclID] = x
		}
		for _, x := range element.DefinedFileExtract {
			element.FileExtract[x.HclID] = x
		}
		for _, x := range element.DefinedIdentities {
			obj, found := combinedConfigs.Identities[x.HclID]
			if !found {
				combinedConfigs.Identities[x.HclID] = x
				continue
			}
			fmt.Println("Stored: ", obj)
			fmt.Println("New: ", x)
		}
	}
	return combinedConfigs, nil
}

func main() {
	client, err := ent.Open("sqlite3", "file:test.sqlite?_loc=auto&cache=shared&_fk=1")
	ctx := context.Background()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	tloader := NewLoader()
	tloader.ParseConfigFile("/home/red/Documents/infra/envs/fred/env.laforge")
	loadedConfig, err := tloader.Bind()
	identities, _ := createIdentities(ctx, client, loadedConfig.Identities, "/envs/fred12")

	// envs, _ := createEnviroments(ctx, client, loadedConfig.Environments)
	// user, _ := envs[0].QueryEnvironmentToUser().Only(ctx)
	// backenv := user.QueryUserToEnvironment().OnlyX(ctx) // Will Fail with multiple users connected to the Env
	// log.Println(envs[0])
	// log.Println(backenv)
	fmt.Println(identities)
}

func createEnviroments(ctx context.Context, client *ent.Client, configEnvs map[string]*ent.Environment) ([]*ent.Environment, error) {
	bulk := []*ent.EnvironmentCreate{}
	returnedEnvironment := []*ent.Environment{}
	for _, cEnviroment := range configEnvs {
		entEnv, err := client.Environment.
			Query().
			Where(environment.HclIDEQ(cEnviroment.HclID)).
			Only(ctx)
		if err != nil {
			if err == err.(*ent.NotFoundError) {
				createdQuery := client.Environment.Create().
					SetHclID(cEnviroment.HclID).
					SetAdminCidrs(cEnviroment.AdminCidrs).
					SetBuilder(cEnviroment.Builder).
					SetCompetitionID(cEnviroment.CompetitionID).
					SetConfig(cEnviroment.Config).
					SetDescription(cEnviroment.Description).
					SetExposedVdiPorts(cEnviroment.ExposedVdiPorts).
					SetName(cEnviroment.Name).
					SetRevision(cEnviroment.Revision).
					SetTags(cEnviroment.Tags).
					SetTeamCount(cEnviroment.TeamCount)
				bulk = append(bulk, createdQuery)
				continue
			}
		}
		updatedEnv, err := entEnv.Update().
			SetHclID(cEnviroment.HclID).
			SetAdminCidrs(cEnviroment.AdminCidrs).
			SetBuilder(cEnviroment.Builder).
			SetCompetitionID(cEnviroment.CompetitionID).
			SetConfig(cEnviroment.Config).
			SetDescription(cEnviroment.Description).
			SetExposedVdiPorts(cEnviroment.ExposedVdiPorts).
			SetName(cEnviroment.Name).
			SetRevision(cEnviroment.Revision).
			SetTags(cEnviroment.Tags).
			SetTeamCount(cEnviroment.TeamCount).
			Save(ctx)
		if err != nil {
			log.Fatalf("failed creating user: %v", err)
			return nil, err
		}
		returnedEnvironment = append(returnedEnvironment, updatedEnv)
	}
	if len(bulk) > 0 {
		dbEnv, err := client.Environment.CreateBulk(bulk...).Save(ctx)
		if err != nil {
			log.Fatalf("failed creating user: %v", err)
			return nil, err
		}
		returnedEnvironment = append(returnedEnvironment, dbEnv...)
	}
	return returnedEnvironment, nil
}

func createUsers(ctx context.Context, client *ent.Client, configUsers []*ent.User, envHclID string) ([]*ent.User, error) {
	bulk := []*ent.UserCreate{}
	returnedUsers := []*ent.User{}
	for _, cuser := range configUsers {
		entUser, err := client.User.
			Query().
			Where(
				user.And(
					user.HclIDEQ(cuser.HclID),
					user.HasUserToEnvironmentWith(environment.HclIDEQ(envHclID)),
				),
			).
			Only(ctx)
		if err != nil {
			if err == err.(*ent.NotFoundError) {
				createdQuery := client.User.Create().
					SetHclID(cuser.HclID).
					SetEmail(cuser.Email).
					SetUUID(cuser.UUID).
					SetName(cuser.Name)
				bulk = append(bulk, createdQuery)
				continue
			}
		}
		_, err = entUser.Update().
			SetHclID(cuser.HclID).
			SetEmail(cuser.Email).
			SetUUID(cuser.UUID).
			SetName(cuser.Name).
			Save(ctx)
		if err != nil {
			log.Fatalf("failed creating user: %v", err)
			return nil, err
		}
	}
	if len(bulk) > 0 {
		dbUsers, err := client.User.CreateBulk(bulk...).Save(ctx)
		if err != nil {
			log.Fatalf("failed creating user: %v", err)
			return nil, err
		}
		returnedUsers = append(returnedUsers, dbUsers...)
	}
	return returnedUsers, nil
}

func createCompetitions(ctx context.Context, client *ent.Client, configCompetitions map[string]*ent.Competition, envHclID string) ([]*ent.Competition, error) {
	bulk := []*ent.CompetitionCreate{}
	returnedCompetitions := []*ent.Competition{}
	for _, cCompetition := range configCompetitions {
		entCompetition, err := client.Competition.
			Query().
			Where(
				competition.And(
					competition.HclIDEQ(cCompetition.HclID),
					competition.HasCompetitionToEnvironmentWith(environment.HclIDEQ(envHclID)),
				),
			).
			Only(ctx)
		if err != nil {
			if err == err.(*ent.NotFoundError) {
				createdQuery := client.Competition.Create().
					SetConfig(cCompetition.Config).
					SetHclID(cCompetition.HclID).
					SetRootPassword(cCompetition.RootPassword)
				bulk = append(bulk, createdQuery)
				continue
			}
		}
		_, err = entCompetition.Update().
			SetConfig(cCompetition.Config).
			SetHclID(cCompetition.HclID).
			SetRootPassword(cCompetition.RootPassword).
			Save(ctx)
		if err != nil {
			log.Fatalf("failed creating competition: %v", err)
			return nil, err
		}
	}
	if len(bulk) > 0 {
		dbCompetitions, err := client.Competition.CreateBulk(bulk...).Save(ctx)
		if err != nil {
			log.Fatalf("failed creating user: %v", err)
			return nil, err
		}
		returnedCompetitions = append(returnedCompetitions, dbCompetitions...)
	}
	return returnedCompetitions, nil
}

func createHosts(ctx context.Context, client *ent.Client, configHosts map[string]*ent.Host, envHclID string) ([]*ent.Host, error) {
	return nil, nil
}
func createNetworks(ctx context.Context, client *ent.Client, configNetworks map[string]*ent.Network, envHclID string) ([]*ent.Network, error) {
	return nil, nil
}
func createScripts(ctx context.Context, client *ent.Client, configScript map[string]*ent.Script, envHclID string) ([]*ent.Script, error) {
	return nil, nil
}
func createCommands(ctx context.Context, client *ent.Client, configCommands map[string]*ent.Command, envHclID string) ([]*ent.Command, error) {
	return nil, nil
}
func createDNSRecords(ctx context.Context, client *ent.Client, configDNSRecords map[string]*ent.DNSRecord, envHclID string) ([]*ent.DNSRecord, error) {
	return nil, nil
}
func createFileDownload(ctx context.Context, client *ent.Client, configFileDownloads map[string]*ent.FileDownload, envHclID string) ([]*ent.FileDownload, error) {
	bulk := []*ent.FileDownloadCreate{}
	returnedFileDownloads := []*ent.FileDownload{}
	for _, cFileDownload := range configFileDownloads {
		entFileDownload, err := client.FileDownload.
			Query().
			Where(
				filedownload.And(
					filedownload.HclIDEQ(cFileDownload.HclID),
					filedownload.HasFileDownloadToEnvironmentWith(environment.HclIDEQ(envHclID)),
				),
			).
			Only(ctx)
		if err != nil {
			if err == err.(*ent.NotFoundError) {
				createdQuery := client.FileDownload.Create().
					SetHclID(cFileDownload.HclID).
					SetSourceType(cFileDownload.SourceType).
					SetSource(cFileDownload.Source).
					SetDestination(cFileDownload.Destination).
					SetTemplate(cFileDownload.Template).
					SetPerms(cFileDownload.Perms).
					SetDisabled(cFileDownload.Disabled).
					SetMd5(cFileDownload.Md5).
					SetAbsPath(cFileDownload.AbsPath).
					SetTags(cFileDownload.Tags)
				bulk = append(bulk, createdQuery)
				continue
			}
		}
		_, err = entFileDownload.Update().
			SetHclID(cFileDownload.HclID).
			SetSourceType(cFileDownload.SourceType).
			SetSource(cFileDownload.Source).
			SetDestination(cFileDownload.Destination).
			SetTemplate(cFileDownload.Template).
			SetPerms(cFileDownload.Perms).
			SetDisabled(cFileDownload.Disabled).
			SetMd5(cFileDownload.Md5).
			SetAbsPath(cFileDownload.AbsPath).
			SetTags(cFileDownload.Tags).
			Save(ctx)
		if err != nil {
			log.Fatalf("failed creating fileextract: %v", err)
			return nil, err
		}
	}
	if len(bulk) > 0 {
		dbFileDownloads, err := client.FileDownload.CreateBulk(bulk...).Save(ctx)
		if err != nil {
			log.Fatalf("failed creating fileextract: %v", err)
			return nil, err
		}
		returnedFileDownloads = append(returnedFileDownloads, dbFileDownloads...)
	}
	return returnedFileDownloads, nil
}
func createFileDelete(ctx context.Context, client *ent.Client, configFileDeletes map[string]*ent.FileDelete, envHclID string) ([]*ent.FileDelete, error) {
	bulk := []*ent.FileDeleteCreate{}
	returnedFileDeletes := []*ent.FileDelete{}
	for _, cFileDelete := range configFileDeletes {
		entFileDelete, err := client.FileDelete.
			Query().
			Where(
				filedelete.And(
					filedelete.HclIDEQ(cFileDelete.HclID),
					filedelete.HasFileDeleteToEnvironmentWith(environment.HclIDEQ(envHclID)),
				),
			).
			Only(ctx)
		if err != nil {
			if err == err.(*ent.NotFoundError) {
				createdQuery := client.FileDelete.Create().
					SetHclID(cFileDelete.HclID).
					SetPath(cFileDelete.Path).
					SetTags(cFileDelete.Tags)
				bulk = append(bulk, createdQuery)
				continue
			}
		}
		_, err = entFileDelete.Update().
			SetHclID(cFileDelete.HclID).
			SetPath(cFileDelete.Path).
			SetTags(cFileDelete.Tags).
			Save(ctx)
		if err != nil {
			log.Fatalf("failed creating fileextract: %v", err)
			return nil, err
		}
	}
	if len(bulk) > 0 {
		dbFileDelete, err := client.FileDelete.CreateBulk(bulk...).Save(ctx)
		if err != nil {
			log.Fatalf("failed creating fileextract: %v", err)
			return nil, err
		}
		returnedFileDeletes = append(returnedFileDeletes, dbFileDelete...)
	}
	return returnedFileDeletes, nil
}
func createFileExtract(ctx context.Context, client *ent.Client, configFileExtracts map[string]*ent.FileExtract, envHclID string) ([]*ent.FileExtract, error) {
	bulk := []*ent.FileExtractCreate{}
	returnedFileExtracts := []*ent.FileExtract{}
	for _, cFileExtract := range configFileExtracts {
		entFileExtract, err := client.FileExtract.
			Query().
			Where(
				fileextract.And(
					fileextract.HclIDEQ(cFileExtract.HclID),
					fileextract.HasFileExtractToEnvironmentWith(environment.HclIDEQ(envHclID)),
				),
			).
			Only(ctx)
		if err != nil {
			if err == err.(*ent.NotFoundError) {
				createdQuery := client.FileExtract.Create().
					SetDestination(cFileExtract.Destination).
					SetHclID(cFileExtract.HclID).
					SetSource(cFileExtract.Source).
					SetTags(cFileExtract.Tags).
					SetType(cFileExtract.Type)
				bulk = append(bulk, createdQuery)
				continue
			}
		}
		_, err = entFileExtract.Update().
			SetDestination(cFileExtract.Destination).
			SetHclID(cFileExtract.HclID).
			SetSource(cFileExtract.Source).
			SetTags(cFileExtract.Tags).
			SetType(cFileExtract.Type).
			Save(ctx)
		if err != nil {
			log.Fatalf("failed creating fileextract: %v", err)
			return nil, err
		}
	}
	if len(bulk) > 0 {
		dbFileExtracts, err := client.FileExtract.CreateBulk(bulk...).Save(ctx)
		if err != nil {
			log.Fatalf("failed creating fileextract: %v", err)
			return nil, err
		}
		returnedFileExtracts = append(returnedFileExtracts, dbFileExtracts...)
	}
	return returnedFileExtracts, nil
}
func createIdentities(ctx context.Context, client *ent.Client, configIdentities map[string]*ent.Identity, envHclID string) ([]*ent.Identity, error) {
	bulk := []*ent.IdentityCreate{}
	returnedIdentities := []*ent.Identity{}
	for _, cIdentity := range configIdentities {
		entIdentity, err := client.Identity.
			Query().
			Where(
				identity.And(
					identity.HclIDEQ(cIdentity.HclID),
					identity.HasIdentityToEnvironmentWith(environment.HclIDEQ(envHclID)),
				),
			).
			Only(ctx)
		if err != nil {
			if err == err.(*ent.NotFoundError) {
				createdQuery := client.Identity.Create().
					SetAvatarFile(cIdentity.AvatarFile).
					SetDescription(cIdentity.Description).
					SetEmail(cIdentity.Email).
					SetFirstName(cIdentity.FirstName).
					SetHclID(cIdentity.HclID).
					SetLastName(cIdentity.LastName).
					SetPassword(cIdentity.Password).
					SetVars(cIdentity.Vars).
					SetTags(cIdentity.Tags)
				bulk = append(bulk, createdQuery)
				continue
			}
		}
		_, err = entIdentity.Update().
			SetAvatarFile(cIdentity.AvatarFile).
			SetDescription(cIdentity.Description).
			SetEmail(cIdentity.Email).
			SetFirstName(cIdentity.FirstName).
			SetHclID(cIdentity.HclID).
			SetLastName(cIdentity.LastName).
			SetPassword(cIdentity.Password).
			SetVars(cIdentity.Vars).
			SetTags(cIdentity.Tags).
			Save(ctx)
		if err != nil {
			log.Fatalf("failed creating competition: %v", err)
			return nil, err
		}
	}
	if len(bulk) > 0 {
		dbIdentities, err := client.Identity.CreateBulk(bulk...).Save(ctx)
		if err != nil {
			log.Fatalf("failed creating user: %v", err)
			return nil, err
		}
		returnedIdentities = append(returnedIdentities, dbIdentities...)
	}
	return returnedIdentities, nil
}
