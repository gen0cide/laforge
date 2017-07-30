package competition

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"

	yaml "gopkg.in/yaml.v2"

	"golang.org/x/crypto/ssh"

	"github.com/fatih/color"
)

const (
	// LF_HOME is the base directory for competition development
	LF_HOME = "LF_HOME"

	// LF_ENV is the currently assigned LF environment
	LF_ENV = "LF_ENV"
)

func GetHome() string {
	return os.Getenv(LF_HOME)
}

func GetEnv() string {
	return os.Getenv(LF_ENV)
}

func EnvSet() bool {
	if len(os.Getenv(LF_ENV)) > 1 {
		return true
	}
	return false
}

func HomeSet() bool {
	if len(os.Getenv(LF_HOME)) > 1 {
		return true
	}
	return false
}

func ValidateHome() {
	if !HomeSet() {
		LogFatal("LF_HOME environment variable is not set. Run the init command to configure this.")
	}
	if !HomeExists() || !HomeValid() {
		LogFatal("LF_HOME is corrupted or not set to a valid laforge specification. Check the docs!")
	}
}

func ValidateEnv() {
	ValidateHome()
	if !EnvSet() {
		LogFatal("LF_ENV environment variable is not set. List known environments with the ls subcommand or create a new one with create.")
	}
	if !EnvExists() || !EnvValid() {
		LogFatal("LF_ENV is corrupted or not set to a valid laforge specification. Check the docs!")
	}
}

func HomeExists() bool {
	return PathExists(GetHome())
}

func EnvExists() bool {
	return PathExists(GetEnv())
}

func EnvDirExistsByName(name string) bool {
	return PathExists(filepath.Join(GetHome(), "environments", name))
}

func MakeSSHKeyPair(pubKeyPath, privateKeyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// generate and write private key as PEM
	privateKeyFile, err := os.Create(privateKeyPath)
	defer privateKeyFile.Close()
	if err != nil {
		return err
	}
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pubKeyPath, ssh.MarshalAuthorizedKey(pub), 0644)
}

func CreateHomeConfig() {
	starterComp := Competition{}

	yamlOutput, err := yaml.Marshal(&starterComp)
	if err != nil {
		LogFatal("YAML Marshaling Error: " + err.Error())
	}
	err = ioutil.WriteFile(filepath.Join(GetHome(), "config", "config.yml"), yamlOutput, 0644)
	if err != nil {
		LogError("Error generating SSH Key: " + err.Error())
	}
}

func ValidName(name string) bool {
	match, _ := regexp.MatchString("^[a-z0-9]{1,16}$", name)
	return match
}

func ValidPrefix(prefix string) bool {
	match, _ := regexp.MatchString("^[a-z]{1,6}$", prefix)
	return match
}

func CreateHome() {
	os.MkdirAll(GetHome(), os.ModePerm)
	os.MkdirAll(filepath.Join(GetHome(), "config"), os.ModePerm)
	os.MkdirAll(filepath.Join(GetHome(), "scripts"), os.ModePerm)
	os.MkdirAll(filepath.Join(GetHome(), "files"), os.ModePerm)
	os.MkdirAll(filepath.Join(GetHome(), "apps"), os.ModePerm)
	os.MkdirAll(filepath.Join(GetHome(), "templates"), os.ModePerm)
	os.MkdirAll(filepath.Join(GetHome(), "users"), os.ModePerm)
	os.MkdirAll(filepath.Join(GetHome(), "environments"), os.ModePerm)
	err := MakeSSHKeyPair(filepath.Join(GetHome(), "config", "infra.pem.pub"), filepath.Join(GetHome(), "config", "infra.pem"))
	if err != nil {
		LogError("Error generating SSH Key: " + err.Error())
	}
	if err != nil {
		LogError("Error generating SSH Key: " + err.Error())
	}
	CreateHomeConfig()
}

func HomeValid() bool {
	// LF_HOME should look like this:
	// folder =>
	//   config (folder)
	//    - config.yml
	//    - infra.pem
	//    - infra.pem.pub
	//   scripts (folder)
	//   files (folder)
	//   apps (folder)
	//   templates (folder)
	//   users (folder)
	//   environments (folder)

	homeValid := true

	if !PathExists(path.Join(GetHome(), "config")) {
		LogError("No config/ folder located in LF_HOME")
		homeValid = false
	}
	if !PathExists(path.Join(GetHome(), "config", "config.yml")) {
		LogError("No config/config.yml located in LF_HOME")
		homeValid = false
	}
	if !PathExists(path.Join(GetHome(), "config", "infra.pem")) {
		LogError("No config/infra.pem located in LF_HOME")
		homeValid = false
	}
	if !PathExists(path.Join(GetHome(), "config", "infra.pem.pub")) {
		LogError("No config/infra.pem.pub located in LF_HOME")
		homeValid = false
	}
	if !PathExists(path.Join(GetHome(), "scripts")) {
		LogError("No scripts/ located in LF_HOME")
		homeValid = false
	}
	if !PathExists(path.Join(GetHome(), "files")) {
		LogError("No files/ located in LF_HOME")
		homeValid = false
	}
	if !PathExists(path.Join(GetHome(), "apps")) {
		LogError("No apps/ located in LF_HOME")
		homeValid = false
	}
	if !PathExists(path.Join(GetHome(), "templates")) {
		LogError("No templates/ located in LF_HOME")
		homeValid = false
	}
	if !PathExists(path.Join(GetHome(), "users")) {
		LogError("No users/ located in LF_HOME")
		homeValid = false
	}
	if !PathExists(path.Join(GetHome(), "environments")) {
		LogError("No environments/ located in LF_HOME")
		homeValid = false
	}

	if homeValid == false {
		LogError("Your LF_HOME directory isn't setup correctly. Check the docs or use the init subcommand to create a new one.")
	}

	return homeValid
}

func EnvValid() bool {
	// LF_ENV should look like this:
	// folder =>
	//   terraform (folder)
	//   env.yml
	//   history.log

	envValid := true

	if !PathExists(path.Join(GetEnv(), "terraform")) {
		LogError("No terraform/ located in LF_ENV")
		envValid = false
	}
	if !PathExists(path.Join(GetEnv(), "env.yml")) {
		LogError("No env.yml located in LF_ENV")
		envValid = false
	}
	if !PathExists(path.Join(GetEnv(), "history.log")) {
		LogError("No history.log located in LF_ENV")
		envValid = false
	}

	return envValid
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	LogError("OS File Error: " + err.Error())
	return false
}

func LogError(msg string) {
	white := color.New(color.FgHiWhite).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()
	fmt.Printf("%s%s%s %s\n", white("["), yellow("ERROR"), white("]"), msg)
}

func LogFatal(msg string) {
	white := color.New(color.FgHiWhite).SprintFunc()
	red := color.New(color.FgHiRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	fmt.Printf("%s%s%s %s\n", white("["), red("FATAL"), white("]"), yellow(msg))
	os.Exit(1)
}

func Log(msg string) {
	white := color.New(color.FgHiWhite).SprintFunc()
	blue := color.New(color.FgHiBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf(" %s%s%s %s\n", white("["), blue("INFO"), white("]"), green(msg))
}

func LogPlain(msg string) {
	fmt.Printf("\n        %s\n\n", msg)
}

func LogEnvs(envs map[*Environment]bool) {
	Log(" == Environment List == ")
	white := color.New(color.FgWhite).SprintFunc()
	cyan := color.New(color.FgHiCyan).SprintFunc()
	for env, curr := range envs {
		if curr == true {
			fmt.Printf("      %s %s %s\n", cyan("*"), cyan(env.Name), cyan("(current)"))
		} else {
			fmt.Printf("        %s\n", white(env.Name))
		}
	}
}
