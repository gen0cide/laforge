package command

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/deprecated/competition"
	"github.com/shiena/ansicolor"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func CmdSsh(c *cli.Context) {
	hostname := c.Args().Get(0)
	if len(hostname) < 1 {
		competition.LogFatal("You did not provide a hostname to use.")
	}

	comp, env := InitConfig()
	sshHosts := env.NewSSHConfig()

	var (
		ip     string
		socket string
		port   = "22"
		user   = "root"
	)

	if val, ok := sshHosts.Hosts[hostname]; ok {
		ip = val
		socket = ip + ":" + port
	} else {
		competition.LogFatal("Unknown host: " + hostname)
	}

	publicKey, err := PublicKeyFile(comp.SSHPrivateKeyPath())
	if err != nil {
		panic(err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			publicKey,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", socket, config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}
	defer conn.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		panic(err)
	}

	session, err := conn.NewSession()
	if err != nil {
		log.Fatal("unable to create session: ", err)
	}
	defer session.Close()

	session.Stdout = ansicolor.NewAnsiColorWriter(os.Stdout)
	session.Stderr = ansicolor.NewAnsiColorWriter(os.Stderr)
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.ECHOCTL:       1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", termHeight, termWidth, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}

	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}

	session.Wait()
	terminal.Restore(fd, oldState)
}

func PublicKeyFile(file string) (ssh.AuthMethod, error) {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}
