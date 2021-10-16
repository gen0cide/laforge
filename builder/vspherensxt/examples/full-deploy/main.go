package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gen0cide/laforge/builder"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/environment"
	"github.com/gen0cide/laforge/logging"
)

func main() {

	pgHost, ok := os.LookupEnv("PG_URI")
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

	defaultLogger := logging.CreateNewLogger("./output.lfglog")

	fmt.Println("Creating vSphere/NSX-T builder...")
	vsphereNsxt, err := builder.NewVSphereNSXTBuilder(env, &defaultLogger)
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

	for _, team := range teams {
		fmt.Printf("Networks for Team %d\n", team.TeamNumber)
		pnets, err := team.QueryTeamToProvisionedNetwork().All(ctx)
		if err != nil {
			log.Fatalf("error while querying provisioned netowrks from team: %v", err)
		}
		for _, pnet := range pnets {
			fmt.Printf("\t%s | %s\n", pnet.Name, pnet.Cidr)

			err = vsphereNsxt.DeployNetwork(ctx, pnet)
			if err != nil {
				fmt.Printf("\tERROR: %v\n", err)
				continue
			} else {
				fmt.Println("\tOK")
			}

			// phosts, err := pnet.QueryProvisionedNetworkToProvisionedHost().All(ctx)
			// if err != nil {
			// 	log.Fatalf("error while querying provisioned hosts from provisioned network: %v", err)
			// }
			// if len(phosts) > 0 {
			// 	fmt.Println("Syncing...")
			// 	time.Sleep(30 * time.Second)

			// 	for _, phost := range phosts {
			// 		host, err := phost.QueryProvisionedHostToHost().Only(ctx)
			// 		if err != nil {
			// 			log.Fatalf("error while querying host from provisioned host")
			// 		}
			// 		fmt.Printf("\t\t%s | %s\n", host.Hostname, host.OS)

			// 		err = vsphereNsxt.DeployHost(ctx, phost)
			// 		if err != nil {
			// 			fmt.Printf("\t\tERROR: %v\n", err)
			// 			continue
			// 		} else {
			// 			fmt.Println("\t\tOK")
			// 		}
			// 	}
			// }
		}
	}
}
