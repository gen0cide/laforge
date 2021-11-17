package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gen0cide/laforge/builder/vspherensxt/vsphere"
	"github.com/sirupsen/logrus"
)

func main() {
	httpClient := http.Client{}

	baseUrl, urlExists := os.LookupEnv("VSPHERE_URL")
	username, usernameExists := os.LookupEnv("VSPHERE_USERNAME")
	password, passwordExists := os.LookupEnv("VSPHERE_PASSWORD")
	resourcePoolName, resourcePoolExists := os.LookupEnv("VSPHERE_RESOURCE_POOL")
	if !urlExists || !usernameExists || !passwordExists || !resourcePoolExists {
		log.Fatalf("please set VSPHERE_URL (exists? %t), VSPHERE_USERNAME (exists? %t), VSPHERE_PASSWORD (exists? %t), and VSPHERE_RESOURCE_POOL (exists? %t)", urlExists, usernameExists, passwordExists, resourcePoolExists)
	}
	vs := vsphere.VSphere{
		HttpClient: httpClient,
		BaseUrl:    baseUrl,
		Username:   username,
		Password:   password,
	}

	vsphere.InitializeGovmomi(&vs, baseUrl, username, password)

	ctx := context.Background()

	resourcePool, err := vs.GetResourcePoolByName(ctx, resourcePoolName)
	if err != nil {
		log.Fatalf("error while getting folders: %v", err)
	}
	logrus.WithFields(logrus.Fields{
		"identifier": resourcePool.Reference().Value,
		"path":       resourcePool.InventoryPath,
	}).Info(resourcePool.Name())
}
