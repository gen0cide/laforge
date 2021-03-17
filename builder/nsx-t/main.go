package main

import (
	"context"
	"fmt"
	"log"
	"os"

	nsxt "github.com/vmware/go-vmware-nsxt"
)

var (
	host = "nsx01.cyberrange.rit.edu"
	username = os.Getenv("GOVMOMI_USERNAME")
	password = os.Getenv("GOVMOMI_PASSWORD")
)

func main() {
	cfg := nsxt.Configuration{
		BasePath:             "/api/v1",
		Host:                 host,
		Scheme:               "https",
		UserName:             username,
		Password:             password,
		Insecure:             true,
	}

	client, err := nsxt.NewAPIClient(&cfg)
	if err != nil {
		log.Fatalf("error initializing NSX-T client: %v\n", err)
	}

	fmt.Println("NX-T initialized...")

	ctx := context.Background()
	// client.Context

	infra, res, err := client.PolicyApi.ReadInfra(ctx)
	if err != nil {
		log.Fatalf("Error reading infra: %v\n", err)
	}
	fmt.Printf("HTTP Status: %s\n", res.Status)
	fmt.Println(infra)
}