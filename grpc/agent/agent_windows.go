//go:build windows
// +build windows

package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	wapi "github.com/iamacarpet/go-win64api"
)

// RebootSystem Reboots Host Operating System
func RebootSystem() {
	// This is how to properlly rebot windows
	user32 := syscall.MustLoadDLL("user32")
	defer user32.Release()

	exitwin := user32.MustFindProc("ExitWindowsEx")

	r1, _, _ := exitwin.Call(0x02, 0)
	if r1 != 1 {
		ExecuteCommand("cmd", "/C", "shutdown", "/r", "/f")
	}

	time.Sleep(1 * time.Hour) // sleep forever bc we need to restart
}

// CreateSystemUser Creates User with specified password.
func CreateSystemUser(username string, password string) error {
	_, err := wapi.UserAdd(username, username, password)
	return err
}

// ChangeSystemUserPassword Change user password.
func ChangeSystemUserPassword(username string, password string) error {
	_, err := wapi.ChangePassword(username, password)
	return err
}

// AddSystemUserGroup Add user to group.
func AddSystemUserGroup(groupname string, username string) error {
	_, err := wapi.LocalGroupAddMembers(groupname, []string{username})
	return err
}

func SystemDownloadFile(path, url string) error {
	retryCount := 5
	var resp *http.Response
	var err error
	for i := 0; i < retryCount; i++ {
		// Get the data
		resp, err = http.Get(url)
		if err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(absolutePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Convert Unix line endings to windows line endings
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	body := strings.Replace(string(bodyBytes), "\n", "\r\n", -1)

	// Write the body to file
	_, err = io.WriteString(out, body)
	return err
}

// SystemExecuteCommand Runs the Command that is inputted and either returns the error or output
func SystemExecuteCommand(command string, args ...string) (string, error) {
	var err error
	_, err = os.Stat(command)
	output := ""
	if err == nil {
		// Make sure we have rwx permissions if it's a script
		err = os.Chmod(command, 0700)
		if err != nil {
			return output, err
		}
	}
	// Execute the command
	arguments := []string{}
	arguments = append(arguments, command)
	arguments = append(arguments, args...)
	cmd := exec.Command("powershell.exe", arguments...)
	out, err := cmd.CombinedOutput()
	return string(out), err
	// retryCount := 5
	// for i := 0; i < retryCount; i++ {
	// 	// Get the data
	// 	cmd := exec.Command("powershell.exe", arguments...)
	// 	out, err := cmd.CombinedOutput()
	// 	if err == nil {
	// 		output = string(out)
	// 		break
	// 	}
	// 	time.Sleep(1 * time.Minute)
	// }
	// if err != nil {
	// 	return output, err
	// }
	// return output, nil
}

func GetSystemDependencies() []string {
	return []string{}
}
