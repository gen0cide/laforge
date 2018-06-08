package competition

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	mr "math/rand"

	"golang.org/x/crypto/ssh"

	"github.com/fatih/color"
)

const (
	LF_HOME       = "LF_HOME"
	LF_HOME_FILE  = ".lf_home"
	LF_ENV        = "LF_ENV"
	LF_ENV_FILE   = ".lf_env"
	SPECIAL_CHARS = "!@#$%^"
)

func GetHome() string {
	u, err := user.Current()
	if err != nil {
		LogFatal("Error getting current user: " + err.Error())
	}
	hf := filepath.Join(u.HomeDir, LF_HOME_FILE)
	if PathExists(hf) {
		d, err := ioutil.ReadFile(hf)
		if err != nil {
			LogFatal("Error reading LF_HOME_FILE: " + err.Error())
		}
		if len(d) > 1 {
			return string(d)
		}
	}

	return ""
}

func GetEnv() string {
	u, err := user.Current()
	if err != nil {
		LogFatal("Error getting current user: " + err.Error())
	}
	ef := filepath.Join(u.HomeDir, LF_ENV_FILE)
	if PathExists(ef) {
		d, err := ioutil.ReadFile(ef)
		if err != nil {
			LogFatal("Error reading LF_ENV_FILE: " + err.Error())
		}
		if len(d) > 1 {
			return string(d)
		}
	}

	return ""
}

func SetHome(val string) {
	u, err := user.Current()
	if err != nil {
		LogFatal("Error getting current user: " + err.Error())
	}
	ef := filepath.Join(u.HomeDir, LF_HOME_FILE)
	err = ioutil.WriteFile(ef, []byte(val), 0644)
	if err != nil {
		LogFatal("Could not write LF_HOME_FILE! (~/.lf_home): " + err.Error())
	}
}

func SetEnv(val string) {
	u, err := user.Current()
	if err != nil {
		LogFatal("Error getting current user: " + err.Error())
	}
	ef := filepath.Join(u.HomeDir, LF_ENV_FILE)
	err = ioutil.WriteFile(ef, []byte(val), 0644)
	if err != nil {
		LogFatal("Could not write LF_ENV_FILE! (~/.lf_env): " + err.Error())
	}
}

func EnvSet() bool {
	if len(GetEnv()) > 1 {
		return true
	}
	return false
}

func HomeSet() bool {
	if len(GetHome()) > 1 {
		return true
	}
	return false
}

func ValidateHome() {
	if !HomeSet() {
		LogFatal("LF_HOME is not defined. Run the init command to configure this.")
	}
	if !HomeExists() || !HomeValid() {
		LogFatal("LF_HOME is corrupted or not set to a valid laforge specification. Check the docs!")
	}
}

func ValidateEnv() {
	ValidateHome()
	if !EnvSet() {
		LogFatal("LF_ENV is not defined. List known environments with the ls subcommand or create a new one with create.")
	}
	if !EnvExists() || !EnvValid() {
		LogFatal("LF_ENV is corrupted or not set to a valid laforge specification. Check the docs!")
	}
}

func HomeExists() bool {
	return PathExists(GetHome())
}

func EnvExists() bool {
	return EnvDirExistsByName(GetEnv())
}

func EnvDirByName(name string) string {
	return filepath.Join(GetHome(), "environments", name)
}

func CurrentEnvDir() string {
	return EnvDirByName(GetEnv())
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

func TouchFile(path string) {
	os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
}

func CreateHomeConfig() {
	err := ioutil.WriteFile(filepath.Join(GetHome(), "config", "config.yml"), MustAsset("config.yml"), 0644)
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

func Contains(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func CreateHome() {
	os.MkdirAll(GetHome(), os.ModePerm)
	os.MkdirAll(filepath.Join(GetHome(), "config"), os.ModePerm)
	os.OpenFile(filepath.Join(GetHome(), "config", ".gitkeep"), os.O_RDONLY|os.O_CREATE, 0644)
	ioutil.WriteFile(filepath.Join(GetHome(), "config", "employees.json"), MustAsset("employees.json"), 0644)
	os.MkdirAll(filepath.Join(GetHome(), "scripts"), os.ModePerm)
	os.OpenFile(filepath.Join(GetHome(), "scripts", ".gitkeep"), os.O_RDONLY|os.O_CREATE, 0644)
	os.MkdirAll(filepath.Join(GetHome(), "files"), os.ModePerm)
	os.OpenFile(filepath.Join(GetHome(), "files", ".gitkeep"), os.O_RDONLY|os.O_CREATE, 0644)
	os.MkdirAll(filepath.Join(GetHome(), "apps"), os.ModePerm)
	os.OpenFile(filepath.Join(GetHome(), "apps", ".gitkeep"), os.O_RDONLY|os.O_CREATE, 0644)
	os.MkdirAll(filepath.Join(GetHome(), "utils"), os.ModePerm)
	os.OpenFile(filepath.Join(GetHome(), "utils", ".gitkeep"), os.O_RDONLY|os.O_CREATE, 0644)
	os.MkdirAll(filepath.Join(GetHome(), "environments"), os.ModePerm)
	os.OpenFile(filepath.Join(GetHome(), "environments", ".gitkeep"), os.O_RDONLY|os.O_CREATE, 0644)
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
	if !PathExists(path.Join(GetHome(), "utils")) {
		LogError("No utils/ located in LF_HOME")
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
	envValid := true

	if !PathExists(path.Join(CurrentEnvDir(), "terraform")) {
		os.MkdirAll(path.Join(CurrentEnvDir(), "terraform"), 0755)
	}
	if !PathExists(path.Join(CurrentEnvDir(), "hosts")) {
		LogError("No terraform/ located in LF_ENV")
		envValid = false
	}
	if !PathExists(path.Join(CurrentEnvDir(), "networks")) {
		LogError("No terraform/ located in LF_ENV")
		envValid = false
	}
	if !PathExists(path.Join(CurrentEnvDir(), "env.yml")) {
		LogError("No env.yml located in LF_ENV")
		envValid = false
	}

	return envValid
}

func CalculateReversePTR(n *Network, h *Host) string {
	octets := strings.Split(n.CIDR, ".")
	return fmt.Sprintf("%d.%s", h.LastOctet, octets[2])
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

func FileToName(path string) string {
	return strings.TrimSuffix(filepath.Base(path), ".yml")
}

func Log(msg string) {
	white := color.New(color.FgHiWhite).SprintFunc()
	blue := color.New(color.FgHiBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("%s%s%s %s\n", white("["), blue("LAFORGE"), white("]"), green(msg))
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

func DictionaryWords() (words []string) {
	file := MustAsset("dictionary.dat")
	words = strings.Split(string(file), "\n")
	return
}

func DeterminedPassword(seed string) string {
	data := []byte{}
	hVal := sha512.Sum512([]byte(seed))
	for _, i := range hVal {
		data = append(data, i)
	}
	genesis := binary.BigEndian.Uint64(data)
	return RandomPasswordFromSeed(int64(genesis))
}

func RandomPasswordFromSeed(seed int64) string {
	r := mr.New(mr.NewSource(seed))
	dict := DictionaryWords()
	passWords := []string{}
	idx := 0
	for idx < 4 {
		word := dict[r.Intn(len(dict))]
		if len(word) > 4 && len(word) < 10 {
			passWords = append(passWords, word)
			idx++
		}
	}
	newPw := strings.Title(strings.Join(passWords, "-"))
	return fmt.Sprintf("%s%s%d", newPw, string(SPECIAL_CHARS[r.Intn(len(SPECIAL_CHARS))]), r.Intn(9))
}

func GetPublicIP() string {
	resp, err := http.Get("http://ipv4.myexternalip.com/raw")
	if err != nil {
		LogFatal("Cannot connect to the internet: " + err.Error())
	}
	defer resp.Body.Close()
	ipData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		LogFatal("Could not read body of IP: " + err.Error())
	}
	return strings.TrimSpace(string(ipData))
}

func NewTermScanner(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		fmt.Printf("nn\n")
		return i + 1, data[0:i], nil
	}
	if i := bytes.IndexByte(data, '?'); i >= 0 {
		return i + 1, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

func ExecInteractiveCommand(command string, args []string) {
	cmd := exec.Command(command, args...)

	out, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(out)
	scanner.Split(NewTermScanner)

	in, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "rm: remove regular empty file 'somefile.txt'" {
			in.Write([]byte("y\n"))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func CustomInternalCNAME(e *Environment, n *Network, c string) string {
	return fmt.Sprintf("%s.%s.%s", c, n.Subdomain, e.Domain)
}

func CustomExternalCNAME(e *Environment, c string) string {
	return fmt.Sprintf("%s.%s", c, e.Domain)
}

func CustomIP(cidr string, offset, id int) string {
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		LogFatal("Not a valid CIDR: " + cidr)
	}
	newIP := ip.To4()
	lastOctet := offset + id
	newIP[3] = byte(lastOctet)
	return newIP.String()
}
