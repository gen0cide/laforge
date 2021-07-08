package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gen0cide/laforge/builder/vspherensxt/vsphere"
)

func main() {
	httpClient := http.Client{}

	baseUrl, urlExists := os.LookupEnv("VSPHERE_URL")
	username, usernameExists := os.LookupEnv("VSPHERE_USERNAME")
	password, passwordExists := os.LookupEnv("VSPHERE_PASSWORD")
	if !urlExists || !usernameExists || !passwordExists {
		log.Fatalf("please set VSPHERE_URL (exists? %t), VSPHERE_USERNAME (exists? %t), and VSPHERE_PASSWORD (exists? %t)", urlExists, usernameExists, passwordExists)
	}
	vshpere := vsphere.VSphere{
		HttpClient: httpClient,
		BaseUrl:    baseUrl,
		Username:   username,
		Password:   password,
	}

	datastoreList, err := vshpere.ListDatastores()
	if err != nil {
		log.Fatalf("error while getting datastores: %v", err)
	}
	for _, datastore := range datastoreList {
		fmt.Printf("%s [%s]: %d Gb / %d Gb\n", datastore.Name, datastore.Identifier, datastore.FreeSpace/1000/1000/1000, datastore.Capacity/1000/1000/1000)
	}
}
