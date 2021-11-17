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
	folderName, folderExists := os.LookupEnv("VSPHERE_FOLDER")
	if !urlExists || !usernameExists || !passwordExists || !folderExists {
		log.Fatalf("please set VSPHERE_URL (exists? %t), VSPHERE_USERNAME (exists? %t), VSPHERE_PASSWORD (exists? %t), and VSPHERE_FOLDER (exists? %t)", urlExists, usernameExists, passwordExists, folderExists)
	}
	vs := vsphere.VSphere{
		HttpClient: httpClient,
		BaseUrl:    baseUrl,
		Username:   username,
		Password:   password,
	}

	vsphere.InitializeGovmomi(&vs, baseUrl, username, password)

	ctx := context.Background()

	folder, err := vs.GetFolderSummaryByName(ctx, folderName)
	if err != nil {
		log.Fatalf("error while getting folders: %v", err)
	}
	logrus.WithFields(logrus.Fields{
		"identifier": folder.Reference().Value,
		"path":       folder.InventoryPath,
	}).Info(folder.Name())
}
