package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
	"github.com/gen0cide/laforge/competition"
)

var (
	TeamID      int
	Parallelism int
)

func CmdTf(c *cli.Context) {
	cli.ShowAppHelpAndExit(c, 0)
}

func SetTeamID(c *cli.Context, e *competition.Environment) {
	teamID := c.Int("team")
	teams := e.TeamIDs()
	Parallelism = c.Int("parallelism")
	if teamID >= 0 && teamID <= teams[len(teams)-1] {
		TeamID = teamID
	} else {
		competition.LogFatal(fmt.Sprintf("Error: %d is not within the valid team range. (0-%d)", teamID, teams[len(teams)-1]))
	}
}

func CmdTfInit(c *cli.Context) {
	TFCheck()
	_, env := InitConfig()
	SetTeamID(c, env)
	tfDir := env.TfDirForTeam(TeamID)
	os.Chdir(tfDir)
	cmdArgs := []string{"init", tfDir}
	TFRun(cmdArgs)
}

func CmdTfPlan(c *cli.Context) {
	TFCheck()
	_, env := InitConfig()
	SetTeamID(c, env)
	tfDir := env.TfDirForTeam(TeamID)
	os.Chdir(tfDir)
	cmdArgs := []string{"plan", fmt.Sprintf("-parallelism=%d", Parallelism), tfDir}
	TFRun(cmdArgs)
}

func CmdTfApply(c *cli.Context) {
	TFCheck()
	_, env := InitConfig()
	SetTeamID(c, env)
	tfDir := env.TfDirForTeam(TeamID)
	os.Chdir(tfDir)
	cmdArgs := []string{"apply", fmt.Sprintf("-parallelism=%d", Parallelism), tfDir}
	TFRun(cmdArgs)
}

func CmdTfTaint(c *cli.Context) {
	TFCheck()
	_, env := InitConfig()
	SetTeamID(c, env)
	tfDir := env.TfDirForTeam(TeamID)
	object := c.Args().Get(0)
	if len(object) < 1 {
		competition.LogFatal("You did not specify a terraform object to taint! Example: laforge tf taint aws_instance.t0_test_host01")
	}
	os.Chdir(tfDir)
	cmdArgs := []string{"taint", object}
	TFRun(cmdArgs)
}

func CmdTfOutput(c *cli.Context) {
	TFCheck()
	_, env := InitConfig()
	SetTeamID(c, env)
	tfDir := env.TfDirForTeam(TeamID)
	os.Chdir(tfDir)
	cmdArgs := []string{"output", tfDir}
	TFRun(cmdArgs)
}

func CmdTfRefresh(c *cli.Context) {
	TFCheck()
	_, env := InitConfig()
	SetTeamID(c, env)
	tfDir := env.TfDirForTeam(TeamID)
	os.Chdir(tfDir)
	cmdArgs := []string{"refresh", "-force", fmt.Sprintf("-parallelism=%d", Parallelism), tfDir}
	TFRun(cmdArgs)
}

func CmdTfState(c *cli.Context) {
	TFCheck()
	_, env := InitConfig()
	SetTeamID(c, env)
	tfDir := env.TfDirForTeam(TeamID)
	os.Chdir(tfDir)
	cmdArgs := []string{"state", "list"}
	TFRun(cmdArgs)
}

func CmdTfDestroy(c *cli.Context) {
	TFCheck()
	_, env := InitConfig()
	SetTeamID(c, env)
	tfDir := env.TfDirForTeam(TeamID)
	os.Chdir(tfDir)
	cmdArgs := []string{"destroy", "-force", fmt.Sprintf("-parallelism=%d", Parallelism), tfDir}
	TFRun(cmdArgs)
}

func TFCheck() {
	_, err := exec.LookPath("terraform")
	if err != nil {
		competition.LogFatal("The terraform executable could not be found in your $PATH!\n\t* Download it at https://www.terraform.io/downloads.html")
	}
}

func TFRun(args []string) {
	cmdName := "terraform"
	cmd := exec.Command(cmdName, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		competition.LogError("Terraform Error: " + err.Error())
	}

	err = cmd.Wait()
	if err != nil {
		competition.LogError("Terraform Error: " + err.Error())
	}

	fmt.Printf("\n")
}
