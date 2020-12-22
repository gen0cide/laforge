// +build unix linux

package main

import (
	"fmt"
	"io"
	"os/exec"
	os "os/user"
	"syscall"
)

// RebootSystem Reboots Host Operating System
func RebootSystem() {
	syscall.Sync()
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}

// CreateSystemUser Create a new User
func CreateSystemUser(username string, password string) error {
	_, err := os.Lookup(username)
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
