package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gen0cide/laforge/builder"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/environment"
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

	builds, err := env.EnvironmentToBuild(ctx)
	if err != nil {
		log.Fatalf("error querying build from env: %v", err)
	}
	teams, err := builds[0].BuildToTeam(ctx)
	if err != nil {
		log.Fatalf("error querying teams from build: %v", err)
	}
	pnets, err := teams[0].TeamToProvisionedNetwork(ctx)
	if err != nil {
		log.Fatalf("error querying provisioned networks from team: %v", err)
	}

	fmt.Println("Tearing down network \"" + pnets[3].Name + "\"")
	err = vsphereNsxt.TeardownNetwork(ctx, pnets[3])
	if err != nil {
		log.Fatalf("error while tearing down network: %v", err)
	}
	fmt.Println("Tearing down network \"" + pnets[4].Name + "\"")
	err = vsphereNsxt.TeardownNetwork(ctx, pnets[4])
	if err != nil {
		log.Fatalf("error while tearing down network: %v", err)
	}
}
