package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gen0cide/laforge/builder"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/environment"
	"github.com/gen0cide/laforge/ent/provisionednetwork"
)

func main() {

	pgHost, ok := os.LookupEnv("PG_HOST")
	client := &ent.Client{}

	if !ok {
		client = ent.PGOpen("postgresql://laforger:laforge@127.0.0.1/laforge")
	} else {
		client = ent.PGOpen(pgHost)
	}

	ctx := context.Background()
	defer ctx.Done()
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	env, err := client.Environment.Query().Where(environment.NameEQ("fred")).Only(ctx)
	if err != nil {
		log.Fatalf("error querying env: %v", err)
	}

	fmt.Println("Creating vSphere/NSX-T builder...")
	vsphereNsxt, err := builder.NewVSphereNSXTBuilder(env)
	if err != nil {
		log.Fatalf("error while creating vCenter/NSX-T builder: %v", err)
	}

	build, err := env.QueryEnvironmentToBuild().Only(ctx)
	if err != nil {
		log.Fatalf("error querying build from env: %v", err)
	}
	teams, err := build.QueryBuildToTeam().All(ctx)
	if err != nil {
		log.Fatalf("error querying teams from build: %v", err)
	}
	pnet, err := teams[0].QueryTeamToProvisionedNetwork().Where(provisionednetwork.NameEQ("corp")).Only(ctx)
	if err != nil {
		log.Fatalf("error querying provisioned networks from team: %v", err)
	}

	fmt.Printf("Found provisioned network \"%s\"\n", pnet.Name)

	phost, err := teams[0].QueryTeamToProvisionedNetwork().QueryProvisionedNetworkToProvisionedHost().Only(ctx)
	if err != nil {
		log.Fatalf("error querying provisioned networks from team: %v", err)
	}
	host, err := phost.QueryProvisionedHostToHost().Only(ctx)
	if err != nil {
		log.Fatalf("error querying provisioned networks from team: %v", err)
	}

	fmt.Printf("Found provisioned host \"%s\"\n", host.Hostname)

	fmt.Printf("Deploying network \"%s\"\n", pnet.Name)
	err = vsphereNsxt.DeployNetwork(ctx, pnet)
	if err != nil {
		log.Fatalf("error while deploying network: %v", err)
	}

	fmt.Println("Waiting 30 secs for systems to sync...")
	time.Sleep(30 * time.Second)

	fmt.Printf("Deploying host \"%s\"\n", host.Hostname)
	err = vsphereNsxt.DeployHost(ctx, phost)
	if err != nil {
		log.Fatalf("error while deploying host: %v", err)
	}
}
