package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsx_policy "github.com/vmware/vsphere-automation-sdk-go/services/nsxt"
)

var (
	host = "https://nsx01.cyberrange.rit.edu"
	// username = os.Getenv("GOVMOMI_USERNAME")
	// password = os.Getenv("GOVMOMI_PASSWORD")
)

func main() {
	// userPassContext := security.NewUserPasswordSecurityContext(username, password)
	// security.AuthenticationHandler.Authenticate(userPassContext)
	// test := security.AuthenticationHandler{}
	// sessions

	httpClient := http.Client{}
	connector := client.NewRestConnector(host, httpClient)
	infraClient := nsx_policy.NewDefaultInfraClient(connector)

	basePath := ""
	filter := ""
	typeFilter := ""
	infra, err := infraClient.Get(&basePath, &filter, &typeFilter)
	if err != nil {
		log.Fatalf("Error getting infra: %v\n", err)
	}
	fmt.Println(infra.DisplayName)
	// cfg := nsxt.Configuration{
	// 	BasePath:             "/api/v1",
	// 	Host:                 host,
	// 	Scheme:               "https",
	// 	UserName:             username,
	// 	Password:             password,
	// 	Insecure:             true,
	// }

	// client, err := nsxt.NewAPIClient(&cfg)
	// if err != nil {
	// 	log.Fatalf("error initializing NSX-T client: %v\n", err)
	// }

	// fmt.Println("NSX-T initialized...")

	// ctx := context.Background()
	// // client.Context

	// infra, res, err := client.PolicyApi.ReadInfra(ctx)
	// if err != nil {
	// 	log.Fatalf("Error reading infra: %v\n", err)
	// }
	// fmt.Printf("HTTP Status: %s\n", res.Status)
	// fmt.Println(infra)
}