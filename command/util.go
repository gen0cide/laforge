package command

import (
	"github.com/gen0cide/laforge/competition"
)

func InitConfig() (*competition.Competition, *competition.Environment) {
	competition.ValidateEnv()
	comp, err := competition.LoadCompetition()
	if err != nil {
		competition.LogFatal("Cannot Load LF_HOME Config: " + err.Error())
	}
	env := comp.CurrentEnv()
	if env == nil {
		competition.LogFatal("Cannot load environment! (Check ~/.lf_env)")
	}
	return comp, env
}

// // Config is the structure of the configuration for the Terraform CLI.
// //
// // This is not the configuration for Terraform itself. That is in the
// // "config" package.
// type TFConfig struct {
// 	Providers    map[string]string
// 	Provisioners map[string]string

// 	DisableCheckpoint          bool `hcl:"disable_checkpoint"`
// 	DisableCheckpointSignature bool `hcl:"disable_checkpoint_signature"`
// }

// // BuiltinConfig is the built-in defaults for the configuration. These
// // can be overridden by user configurations.
// var BuiltinConfig TFConfig

// // PluginOverrides are paths that override discovered plugins, set from
// // the config file.
// var PluginOverrides tfcmd.PluginOverrides

// // ConfigFile returns the default path to the configuration file.
// //
// // On Unix-like systems this is the ".terraformrc" file in the home directory.
// // On Windows, this is the "terraform.rc" file in the application data
// // directory.
// func TFConfigFile() (string, error) {
// 	return configFile()
// }

// // ConfigDir returns the configuration directory for Terraform.
// func TFConfigDir() (string, error) {
// 	return configDir()
// }

// // LoadConfig loads the CLI configuration from ".terraformrc" files.
// func TFLoadConfig(path string) (*Config, error) {
// 	// Read the HCL file and prepare for parsing
// 	d, err := ioutil.ReadFile(path)
// 	if err != nil {
// 		return nil, fmt.Errorf(
// 			"Error reading %s: %s", path, err)
// 	}

// 	// Parse it
// 	obj, err := hcl.Parse(string(d))
// 	if err != nil {
// 		return nil, fmt.Errorf(
// 			"Error parsing %s: %s", path, err)
// 	}

// 	// Build up the result
// 	var result Config
// 	if err := hcl.DecodeObject(&result, obj); err != nil {
// 		return nil, err
// 	}

// 	// Replace all env vars
// 	for k, v := range result.Providers {
// 		result.Providers[k] = os.ExpandEnv(v)
// 	}
// 	for k, v := range result.Provisioners {
// 		result.Provisioners[k] = os.ExpandEnv(v)
// 	}

// 	return &result, nil
// }

// // Merge merges two configurations and returns a third entirely
// // new configuration with the two merged.
// func (c1 *TFConfig) TFMerge(c2 *TFConfig) *TFConfig {
// 	var result TFConfig
// 	result.Providers = make(map[string]string)
// 	result.Provisioners = make(map[string]string)
// 	for k, v := range c1.Providers {
// 		result.Providers[k] = v
// 	}
// 	for k, v := range c2.Providers {
// 		if v1, ok := c1.Providers[k]; ok {
// 			log.Printf("[INFO] Local %s provider configuration '%s' overrides '%s'", k, v, v1)
// 		}
// 		result.Providers[k] = v
// 	}
// 	for k, v := range c1.Provisioners {
// 		result.Provisioners[k] = v
// 	}
// 	for k, v := range c2.Provisioners {
// 		if v1, ok := c1.Provisioners[k]; ok {
// 			log.Printf("[INFO] Local %s provisioner configuration '%s' overrides '%s'", k, v, v1)
// 		}
// 		result.Provisioners[k] = v
// 	}
// 	result.DisableCheckpoint = c1.DisableCheckpoint || c2.DisableCheckpoint
// 	result.DisableCheckpointSignature = c1.DisableCheckpointSignature || c2.DisableCheckpointSignature

// 	return &result
// }

// func RunTF(args []string) int {
// 	// We always need to close the DebugInfo before we exit.
// 	defer terraform.CloseDebugInfo()

// 	log.SetOutput(os.Stderr)

// 	// Load the configuration
// 	config := BuiltinConfig

// 	// Load the configuration file if we have one, that can be used to
// 	// define extra providers and provisioners.
// 	clicfgFile, err := cliConfigFile()
// 	if err != nil {
// 		Ui.Error(fmt.Sprintf("Error loading CLI configuration: \n\n%s", err))
// 		return 1
// 	}

// 	if clicfgFile != "" {
// 		usrcfg, err := LoadConfig(clicfgFile)
// 		if err != nil {
// 			Ui.Error(fmt.Sprintf("Error loading CLI configuration: \n\n%s", err))
// 			return 1
// 		}

// 		config = *config.Merge(usrcfg)
// 	}

// 	// Run checkpoint
// 	go runCheckpoint(&config)

// 	// Make sure we clean up any managed plugins at the end of this
// 	defer plugin.CleanupClients()

// 	// Build the CLI so far, we do this so we can query the subcommand.
// 	cliRunner := &cli.CLI{
// 		Args:       args,
// 		Commands:   Commands,
// 		HelpFunc:   helpFunc,
// 		HelpWriter: os.Stdout,
// 	}

// 	// Prefix the args with any args from the EnvCLI
// 	args, err = mergeEnvArgs(EnvCLI, cliRunner.Subcommand(), args)
// 	if err != nil {
// 		Ui.Error(err.Error())
// 		return 1
// 	}

// 	// Prefix the args with any args from the EnvCLI targeting this command
// 	suffix := strings.Replace(strings.Replace(
// 		cliRunner.Subcommand(), "-", "_", -1), " ", "_", -1)
// 	args, err = mergeEnvArgs(
// 		fmt.Sprintf("%s_%s", EnvCLI, suffix), cliRunner.Subcommand(), args)
// 	if err != nil {
// 		Ui.Error(err.Error())
// 		return 1
// 	}

// 	// We shortcut "--version" and "-v" to just show the version
// 	for _, arg := range args {
// 		if arg == "-v" || arg == "-version" || arg == "--version" {
// 			newArgs := make([]string, len(args)+1)
// 			newArgs[0] = "version"
// 			copy(newArgs[1:], args)
// 			args = newArgs
// 			break
// 		}
// 	}

// 	// Rebuild the CLI with any modified args.
// 	log.Printf("[INFO] CLI command args: %#v", args)
// 	cliRunner = &cli.CLI{
// 		Args:       args,
// 		Commands:   Commands,
// 		HelpFunc:   helpFunc,
// 		HelpWriter: os.Stdout,
// 	}

// 	// Pass in the overriding plugin paths from config
// 	PluginOverrides.Providers = config.Providers
// 	PluginOverrides.Provisioners = config.Provisioners

// 	exitCode, err := cliRunner.Run()
// 	if err != nil {
// 		Ui.Error(fmt.Sprintf("Error executing CLI: %s", err.Error()))
// 		return 1
// 	}

// 	return exitCode
// }

// func cliConfigFile() (string, error) {
// 	mustExist := true
// 	configFilePath := os.Getenv("TERRAFORM_CONFIG")
// 	if configFilePath == "" {
// 		var err error
// 		configFilePath, err = ConfigFile()
// 		mustExist = false

// 		if err != nil {
// 			log.Printf(
// 				"[ERROR] Error detecting default CLI config file path: %s",
// 				err)
// 		}
// 	}

// 	log.Printf("[DEBUG] Attempting to open CLI config file: %s", configFilePath)
// 	f, err := os.Open(configFilePath)
// 	if err == nil {
// 		f.Close()
// 		return configFilePath, nil
// 	}

// 	if mustExist || !os.IsNotExist(err) {
// 		return "", err
// 	}

// 	log.Println("[DEBUG] File doesn't exist, but doesn't need to. Ignoring.")
// 	return "", nil
// }

// // copyOutput uses output prefixes to determine whether data on stdout
// // should go to stdout or stderr. This is due to panicwrap using stderr
// // as the log and error channel.
// func copyOutput(r io.Reader, doneCh chan<- struct{}) {
// 	defer close(doneCh)

// 	pr, err := prefixedio.NewReader(r)
// 	if err != nil {
// 		panic(err)
// 	}

// 	stderrR, err := pr.Prefix(ErrorPrefix)
// 	if err != nil {
// 		panic(err)
// 	}
// 	stdoutR, err := pr.Prefix(OutputPrefix)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defaultR, err := pr.Prefix("")
// 	if err != nil {
// 		panic(err)
// 	}

// 	var stdout io.Writer = os.Stdout
// 	var stderr io.Writer = os.Stderr

// 	if runtime.GOOS == "windows" {
// 		stdout = colorable.NewColorableStdout()
// 		stderr = colorable.NewColorableStderr()

// 		// colorable is not concurrency-safe when stdout and stderr are the
// 		// same console, so we need to add some synchronization to ensure that
// 		// we can't be concurrently writing to both stderr and stdout at
// 		// once, or else we get intermingled writes that create gibberish
// 		// in the console.
// 		wrapped := synchronizedWriters(stdout, stderr)
// 		stdout = wrapped[0]
// 		stderr = wrapped[1]
// 	}

// 	var wg sync.WaitGroup
// 	wg.Add(3)
// 	go func() {
// 		defer wg.Done()
// 		io.Copy(stderr, stderrR)
// 	}()
// 	go func() {
// 		defer wg.Done()
// 		io.Copy(stdout, stdoutR)
// 	}()
// 	go func() {
// 		defer wg.Done()
// 		io.Copy(stdout, defaultR)
// 	}()

// 	wg.Wait()
// }

// func mergeEnvArgs(envName string, cmd string, args []string) ([]string, error) {
// 	v := os.Getenv(envName)
// 	if v == "" {
// 		return args, nil
// 	}

// 	log.Printf("[INFO] %s value: %q", envName, v)
// 	extra, err := shellwords.Parse(v)
// 	if err != nil {
// 		return nil, fmt.Errorf(
// 			"Error parsing extra CLI args from %s: %s",
// 			envName, err)
// 	}

// 	// Find the command to look for in the args. If there is a space,
// 	// we need to find the last part.
// 	search := cmd
// 	if idx := strings.LastIndex(search, " "); idx >= 0 {
// 		search = cmd[idx+1:]
// 	}

// 	// Find the index to place the flags. We put them exactly
// 	// after the first non-flag arg.
// 	idx := -1
// 	for i, v := range args {
// 		if v == search {
// 			idx = i
// 			break
// 		}
// 	}

// 	// idx points to the exact arg that isn't a flag. We increment
// 	// by one so that all the copying below expects idx to be the
// 	// insertion point.
// 	idx++

// 	// Copy the args
// 	newArgs := make([]string, len(args)+len(extra))
// 	copy(newArgs, args[:idx])
// 	copy(newArgs[idx:], extra)
// 	copy(newArgs[len(extra)+idx:], args[idx:])
// 	return newArgs, nil
// }
