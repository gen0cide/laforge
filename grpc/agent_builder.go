package main

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

func buildAgent(agentID string, serverAddress string, binarypath string){

	command:="go build -ldflags=\" -X 'main.clientID="+agentID+"' -X 'main.address="+serverAddress+"'\" -o "+binarypath+" github.com/gen0cide/laforge/grpc/agent"
	cmd := exec.Command("bash","-c",command)
	stdoutStderr, err := cmd.CombinedOutput()
	cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)

}

func main(){
	agents := []map[string]string{}
	pgHost, ok := os.LookupEnv("PG_HOST")
	client := &ent.Client{}

	agents = append(agents, map[string]string{"1": "test/test"} )
	serverAddress := "localhost:50051"

	if !ok {
		client = ent.PGOpen("postgresql://laforger:laforge@127.0.0.1/laforge")
	} else {
		client = ent.PGOpen(pgHost)
	}

	ctx := context.Background()
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	phs,err := client.ProvisionedHost.Query().All(ctx)
	if err != nil {
		log.Fatalf("failed casting UUID to int: %v", err)
	}

	for _, ph := range phs { 
		host, err := ph.QueryHost().Only(ctx)
		if err != nil {
			log.Fatalf("failed casting UUID to int: %v", err)
		}
		hostName := host.Hostname

		switch runtime.GOOS {
		case "windows":
			if !strings.Contains(host.OS,"w2k"){
				continue
			}
		case "linux":
			if strings.Contains(host.OS,"w2k"){
				continue
			}	
		}
	
		pn, err := ph.QueryProvisionedNetwork().Only(ctx)
		if err != nil {
			log.Fatalf("failed casting UUID to int: %v", err)
		}
		networkName := pn.Name

		team , err := pn.QueryProvisionedNetworkToTeam().Only(ctx)
		if err != nil {
			log.Fatalf("failed casting UUID to int: %v", err)
		}
		teamName := team.TeamNumber

		binaryName := filepath.Join("team",fmt.Sprint(teamName),networkName,hostName)
		buildAgent(fmt.Sprint(ph.ID), serverAddress, binaryName)
	}	
}
