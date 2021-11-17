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
	if !urlExists || !usernameExists || !passwordExists {
		log.Fatalf("please set VSPHERE_URL (exists? %t), VSPHERE_USERNAME (exists? %t), and VSPHERE_PASSWORD (exists? %t)", urlExists, usernameExists, passwordExists)
	}
	vs := vsphere.VSphere{
		HttpClient: httpClient,
		BaseUrl:    baseUrl,
		Username:   username,
		Password:   password,
	}

	vsphere.InitializeGovmomi(&vs, baseUrl, username, password)

	ctx := context.Background()

	datastoreList, err := vs.ListDatastores(ctx)
	if err != nil {
		log.Fatalf("error while getting datastores: %v", err)
	}
	for _, datastore := range datastoreList {
		logrus.WithFields(logrus.Fields{
			"identifier":    datastore.Datastore.Value,
			"freeSpace(Gb)": datastore.FreeSpace / 1024 / 1024 / 1024,
			"capacity(Gb)":  datastore.Capacity / 1024 / 1024 / 1024,
		}).Info(datastore.Name)
	}
}
