// +build unix linux

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"syscall"
	"time"
)

// RebootSystem Reboots Host Operating System
func RebootSystem() {
	syscall.Sync()
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)

	time.Sleep(1 * time.Hour) // sleep forever bc we need to restart
}

// CreateSystemUser Create a new User
func CreateSystemUser(username string, password string) error {
	_, err := user.Lookup(username)
	if err != nil {
		ExecuteCommand("useradd", username)
		ChangeSystemUserPassword(username, password)
	}
	return nil
}

// ChangeSystemUserPassword Change user password.
func ChangeSystemUserPassword(username string, password string) error {
	cmd := exec.Command("passwd", username)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Error(err)
	}
	defer stdin.Close()

	passnew := fmt.Sprintf("%s\n%s\n", password, password)

	io.WriteString(stdin, passnew)

	if err = cmd.Start(); err != nil {
		logger.Errorf("An error occured: ", err)
	}

	cmd.Wait()

	return nil
}

// AddSystemUserGroup Change user password.
func AddSystemUserGroup(groupname string, username string) error {
	ExecuteCommand("usermod", "-a", "-G", groupname, username)
	return nil
}

// SystemDownloadFile Download a file with OS specific file endings
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

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// SystemExecuteCommand Runs the Command that is inputted and either returns the error or output
func SystemExecuteCommand(command string, args ...string) (string, error) {
	var err error
	_, err = os.Stat(command)
	// output := ""
	if err == nil {
		// Make sure we have rwx permissions if it's a script
		err = os.Chmod(command, 0700)
		if err != nil {
			return "", err
		}
	}
	// Execute the command
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
	// retryCount := 5
	// for i := 0; i < retryCount; i++ {
	// 	// Get the data
	// 	cmd := exec.Command(command, args...)
	// 	stdout, err := cmd.StdoutPipe()
	// 	if err != nil {
	// 		fmt.Printf("error piping stdout: %v", err)
	// 		continue
	// 	}
	// 	stderr, err := cmd.StderrPipe()
	// 	if err != nil {
	// 		fmt.Printf("error piping stderr: %v", err)
	// 		continue
	// 	}
	// 	err = cmd.Run()
	// 	// out, err := cmd.CombinedOutput()
	// 	if err == nil {
	// 		// output = string(out)
	// 		combinedOutput := io.MultiReader(stdout, stderr)
	// 		var buff []byte
	// 		_, err = combinedOutput.Read(buff)
	// 		if err != nil {
	// 			fmt.Printf("error reading combined output: %v", err)
	// 			continue
	// 		}
	// 		output = string(buff)
	// 		break
	// 	}
	// 	time.Sleep(1 * time.Minute)
	// }
	// if err != nil {
	// 	return output, err
	// }
	// _, err = cmd.Output()
	// if err != nil {
	// 	return err
	// }
	// return string(output)
	// return output, nil
}

func GetSystemDependencies() []string {
	return []string{
		"Requires=network.target",
		"After=network-online.target"}
}
