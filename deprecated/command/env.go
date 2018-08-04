package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/deprecated/competition"
	input "github.com/tcnksm/go-input"
)

func CmdEnv(c *cli.Context) {
	cli.ShowAppHelpAndExit(c, 0)
}

func CmdEnvLs(c *cli.Context) {
	comp, err := competition.LoadCompetition()
	if err != nil {
		competition.LogFatal("Cannot Load LF_HOME Config: " + err.Error())
	}

	competition.LogEnvs(comp.GetEnvs())
}

func CmdEnvUse(c *cli.Context) {
	envName := c.Args().Get(0)
	if len(envName) < 1 {
		competition.LogFatal("You did not provide an environment to use.")
	}
	comp, err := competition.LoadCompetition()
	if err != nil {
		competition.LogFatal("Cannot Load LF_HOME Config: " + err.Error())
	}
	comp.ChangeEnv(envName)

}

func CmdEnvCreate(c *cli.Context) {
	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	query := "Enter environment name"
	newEnvName, err := ui.Ask(query, &input.Options{
		// Read the default val from env var
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})
	if err != nil {
		competition.LogFatal("Fatal Error: " + err.Error())
	}
	query = "Enter environment prefix"
	newEnvPrefix, err := ui.Ask(query, &input.Options{
		// Read the default val from env var
		Required:  true,
		Loop:      true,
		HideOrder: true,
	})
	if err != nil {
		competition.LogFatal("Fatal Error: " + err.Error())
	}
	comp, err := competition.LoadCompetition()
	if err != nil {
		competition.LogFatal("Cannot Load LF_HOME: " + err.Error())
	}
	comp.CreateEnv(newEnvName, newEnvPrefix)
}

func CmdEnvPassword(c *cli.Context) {
	_, env := InitConfig()
	TFCheck()
	podID := c.Args().Get(0)
	if len(podID) < 1 {
		competition.LogFatal("You did not provide a Pod ID to use.")
	}
	podVal, err := strconv.Atoi(podID)
	if err != nil {
		competition.LogFatal("You did not supply a valid team number.")
	}
	dp := env.PodPassword(podVal)
	competition.Log(fmt.Sprintf("Determined Password: %s", dp))
}

func CmdEnvSshConfig(c *cli.Context) {
	_, env := InitConfig()
	env.GenerateSSHConfig()
	competition.Log("SSH Config successfully saved to: " + env.SSHConfigPath())
	return
}

func CmdEnvBashConfig(c *cli.Context) {
	comp, env := InitConfig()
	fmt.Println("# ---------------------------------------------")
	fmt.Println("# Laforge Environment Configuration -----------")
	fmt.Println("# ---------------------------------------------")
	fmt.Println("# To apply the changes to your current shell:")
	fmt.Println("# ")
	fmt.Println("# laforge env bashconfig > /tmp/lf.env && source /tmp/lf.env && rm /tmp/lf.env")
	fmt.Println("# ")
	fmt.Println("# ---------------------------------------------")
	fmt.Printf("export LF_HOME=\"%s\";\n", competition.GetHome())
	fmt.Printf("export LF_ENV=\"%s\";\n", env.Name)
	fmt.Printf("export LF_ENV_HOME=\"%s\";\n", filepath.Join(competition.GetHome(), "environments", env.Name))
	fmt.Printf("export LF_ENV_TF_DIR=\"%s\";\n", env.TfDir())
	fmt.Printf("export LF_SSH_CONFIG=\"%s\";\n", filepath.Join(competition.GetHome(), "environments", env.Name, "ssh.conf"))
	fmt.Printf("export LF_SSH_KEY=\"%s\";\n", comp.SSHPrivateKeyPath())
	fmt.Println("alias lfhome=\"cd $LF_HOME\";")
	fmt.Println("alias lfenv=\"cd $LF_ENV_HOME\";")
	fmt.Println("alias tfhome=\"cd $LF_ENV_TF_DIR\";")
	fmt.Println("alias lfssh=\"ssh -F $LF_SSH_CONFIG -i $LF_SSH_KEY\";")
	fmt.Println("alias lfscp=\"scp -F $LF_SSH_CONFIG -i $LF_SSH_KEY\";")
	fmt.Println("alias tfp=\"cd $LF_ENV_TF_DIR && terraform plan -parallelism=50\";")
	fmt.Println("alias tfa=\"cd $LF_ENV_TF_DIR && terraform apply -parallelism=50\";")
	fmt.Println("alias tfd=\"cd $LF_ENV_TF_DIR && terraform destroy -force -parallelism=50\";")
	fmt.Println("alias lfnm=\"nmap -sS -T4 -Pn\";")
}
