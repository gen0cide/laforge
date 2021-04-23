package grpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gen0cide/laforge/ent"
)

func BuildAgent(agentID string, serverAddress string, binarypath string, isWindows bool) {
	command := ""
	if isWindows {
		command = "CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=\"zcc\" go build -ldflags=\" -X 'main.clientID=" + agentID + "' -X 'main.address=" + serverAddress + "'\" -o " + binarypath + " github.com/gen0cide/laforge/grpc/agent"
	} else {
		command = "CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags=\" -X 'main.clientID=" + agentID + "' -X 'main.address=" + serverAddress + "'\" -o " + binarypath + " github.com/gen0cide/laforge/grpc/agent"
	}
	cmd := exec.Command("bash", "-c", command)
	stdoutStderr, err := cmd.CombinedOutput()
	cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created %s, Output %s\n", binarypath, stdoutStderr)
}

func main() {

	client := &ent.Client{}

	pgHost, ok := os.LookupEnv("PG_HOST")
	if !ok {
		client = ent.PGOpen("postgresql://laforger:laforge@127.0.0.1/laforge")
	} else {
		client = ent.PGOpen(pgHost)
	}

	serverAddress, ok := os.LookupEnv("GRPC_SERVER")
	if !ok {
		serverAddress = "localhost:50051"
	}

	ctx := context.Background()
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	phs, err := client.ProvisionedHost.Query().All(ctx)
	if err != nil {
		log.Fatalf("Failed to Query All Provisioned Hosts: %v", err)
	}

	for _, ph := range phs {
		host, err := ph.QueryProvisionedHostToHost().Only(ctx)
		if err != nil {
			log.Fatalf("Failed to Query Host: %v", err)
		}
		hostName := host.Hostname

		switch runtime.GOOS {
		case "windows":
			if !strings.Contains(host.OS, "w2k") {
				continue
			}
		case "linux":
			if strings.Contains(host.OS, "w2k") {
				continue
			}
		}

		pn, err := ph.QueryProvisionedHostToProvisionedNetwork().Only(ctx)
		if err != nil {
			log.Fatalf("Failed to Query Provisioned Network: %v", err)
		}
		networkName := pn.Name

		team, err := pn.QueryProvisionedNetworkToTeam().Only(ctx)
		if err != nil {
			log.Fatalf("Failed to Query Team: %v", err)
		}
		teamName := team.TeamNumber
		env, err := team.QueryTeamToBuild().QueryBuildToEnvironment().Only(ctx)
		if err != nil {
			log.Fatalf("Failed to Query Enviroment: %v", err)
		}
		envName := env.Name

		binaryName := filepath.Join(envName, "team", fmt.Sprint(teamName), networkName, hostName)
		BuildAgent(fmt.Sprint(ph.ID), serverAddress, binaryName, false)
	}
}
