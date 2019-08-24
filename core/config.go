package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/pkg/errors"
)

var (
	// ErrNoGlobalConfig is thrown when there is no global configuration for a user
	ErrNoGlobalConfig = errors.New("no global configuration exists")

	// ErrNoConfigRootReached is thrown when the entire filesystem has been traversed to the root without locating a valid config file
	ErrNoConfigRootReached = errors.New("no valid config file was located between current directory and the filesystem root")
)

// CreateGlobalConfig writes a User configuration to ~/.laforge/global.laforge
func CreateGlobalConfig(u User) error {
	home := os.Getenv("HOME")
	if home == "" {
		u, err := user.Current()
		if err != nil {
			return err
		}
		home = u.HomeDir
	}
	configDir := filepath.Join(home, ".laforge")
	err := os.MkdirAll(configDir, 0700)
	if err != nil {
		return err
	}
	data, err := RenderHCLv2Object(u)
	if err != nil {
		return err
	}
	globconf := filepath.Join(configDir, "global.laforge")
	return ioutil.WriteFile(globconf, data, 0600)
}

// PathExists is a convenience function to determine if a path exists at it's location
func PathExists(loc string) bool {
	absp, err := filepath.Abs(loc)
	if err != nil {
		cli.Logger.Debugf("Error determining absolute path of %s: %v", loc, err)
		return false
	}
	if _, err := os.Stat(absp); err == nil {
		return true
	}
	return false
}

// GlobalConfigDir attempts to retrieve the user's global configuration directory (dotfile in $HOME)
func GlobalConfigDir() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		u, err := user.Current()
		if err != nil {
			return "", err
		}
		home = u.HomeDir
	}
	configDir := filepath.Join(home, ".laforge")
	return configDir, nil
}

// LocateGlobalConfig attempts to detect a global configuration in the GlobalConfigDir
func LocateGlobalConfig() (string, error) {
	gcd, err := GlobalConfigDir()
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(gcd); err != nil && os.IsNotExist(err) {
		return "", ErrNoGlobalConfig
	}
	globconf := filepath.Join(gcd, "global.laforge")
	if _, err := os.Stat(globconf); err != nil && os.IsNotExist(err) {
		return "", ErrNoGlobalConfig
	}
	return globconf, nil
}

// LocateTeamConfig attempts to locate a valid team.laforge in the current directory or it's parents
func LocateTeamConfig() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return TraverseUpForFile("team.laforge", cwd)
}

// LocateBuildConfig attempts to locate a valid env.laforge in the current directory or it's parents
func LocateBuildConfig() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return TraverseUpForFile("build.laforge", cwd)
}

// LocateEnvConfig attempts to locate a valid env.laforge in the current directory or it's parents
func LocateEnvConfig() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return TraverseUpForFile(envFile, cwd)
}

// LocateBaseConfig attempts to locate a valid base.laforge in the current directory or one of it's parents
func LocateBaseConfig() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return TraverseUpForFile("base.laforge", cwd)
}

// BuildDirFromEnvRoot is a convenience function to get a prebuilt filepath to the build directory for an env root.
func BuildDirFromEnvRoot(envroot string) string {
	return filepath.Join(envroot, "build")
}

// LoadFiles loads a configuration from configuration files passed to it
func LoadFiles(envpath ...string) (*Laforge, error) {
	loader := NewLoader()
	for _, c := range envpath {
		err := loader.ParseConfigFile(c)
		if err != nil {
			return nil, err
		}
	}
	return loader.Bind()
}

// TraverseUpForFile looks for a file by the given filename upwards in the file structure
func TraverseUpForFile(filename, startdir string) (string, error) {
	absPath, err := filepath.Abs(startdir)
	if err != nil {
		return "", err
	}
	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if file.Name() == filename {
			return filepath.Join(absPath, filename), nil
		}
	}
	switch filepath.Clean(absPath) {
	case `/`:
		return "", ErrNoConfigRootReached
	case `C:\`:
		return "", ErrNoConfigRootReached
	case `.`:
		return "", ErrNoConfigRootReached
	}
	return TraverseUpForFile(filename, filepath.Dir(absPath))
}

// TouchGitKeep is a helper function to touch a .gitkeep file within a given directory.
func TouchGitKeep(p string) error {
	keeper := filepath.Join(p, ".gitkeep")
	newFile, err := os.Create(keeper)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("cannot touch .gitkeep inside directory %s", p))
	}
	err = newFile.Close()
	if err != nil {
		return err
	}
	return nil
}
