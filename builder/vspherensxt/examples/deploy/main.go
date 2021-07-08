package main

import (
	"net/http"
	"os"

	"github.com/gen0cide/laforge/builder/vspherensxt/vsphere"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v4"
)

func main() {
	log.SetLevel(log.DebugLevel)
	httpClient := http.Client{}

	baseUrl, urlExists := os.LookupEnv("VSPHERE_URL")
	username, usernameExists := os.LookupEnv("VSPHERE_USERNAME")
	password, passwordExists := os.LookupEnv("VSPHERE_PASSWORD")
	libraryName, libraryNameExists := os.LookupEnv("VSPHERE_CONTENT_LIBRARY")
	templateName, templateNameExists := os.LookupEnv("VSPHERE_TEMPLATE_NAME")
	datastoreName, datastoreNameExists := os.LookupEnv("VSPHERE_DATASTORE")
	resourcePoolName, resourcePoolNameExists := os.LookupEnv("VSPHERE_RESOURCE_POOL")
	folderName, folderNameExists := os.LookupEnv("VSPHERE_FOLDER")
	networkName, networkNameExists := os.LookupEnv("VSPHERE_NETWORK")
	if !urlExists || !usernameExists || !passwordExists || !templateNameExists || !libraryNameExists || !datastoreNameExists || !resourcePoolNameExists || !folderNameExists || !networkNameExists {
		log.WithFields(log.Fields{
			"VSPHERE_URL":             urlExists,
			"VSPHERE_USERNAME":        usernameExists,
			"VSPHERE_PASSWORD":        passwordExists,
			"VSPHERE_CONTENT_LIBRARY": libraryNameExists,
			"VSPHERE_TEMPLATE_NAME":   templateNameExists,
			"VSPHERE_DATASTORE":       datastoreNameExists,
			"VSPHERE_RESOURCE_POOL":   resourcePoolNameExists,
			"VSPHERE_FOLDER":          folderNameExists,
			"VSPHERE_NETWORK":         networkNameExists,
		}).Fatal("some env variables are missing, these currently exist:")
	}
	vs := vsphere.VSphere{
		HttpClient: httpClient,
		BaseUrl:    baseUrl,
		Username:   username,
		Password:   password,
	}

	datastore, err := vs.GetDatastoreByName(datastoreName)
	if err != nil {
		log.Fatalf("error while finding the datastore: %v", err)
	}
	log.WithFields(log.Fields{
		"name":       datastore.Name,
		"identifier": datastore.Identifier,
	}).Info("Found datastore")

	folder, err := vs.GetFolderByName(folderName)
	if err != nil {
		log.Fatalf("error while finding the folder: %v", err)
	}
	log.WithFields(log.Fields{
		"name":       folder.Name,
		"identifier": folder.Identifier,
	}).Info("Found folder")

	resourcePool, err := vs.GetResourcePoolByName(resourcePoolName)
	if err != nil {
		log.Fatalf("error while finding the resource pool: %v", err)
	}
	log.WithFields(log.Fields{
		"name":       resourcePool.Name,
		"identifier": resourcePool.Identifier,
	}).Info("Found resource pool")

	network, err := vs.GetNetworkByName(networkName)
	if err != nil {
		log.Fatalf("error while finding the network: %v", err)
	}
	log.WithFields(log.Fields{
		"name":       network.Name,
		"identifier": network.Identifier,
	}).Info("Found network")

	templateId, err := vs.GetTemplateIDByName(libraryName, templateName)
	if err != nil {
		log.Fatalf("error while searching content library for \"%s\": %v", templateName, err)
	}
	log.WithFields(log.Fields{
		"name":       templateName,
		"identifier": templateId,
	}).Info("Found template")

	template, err := vs.GetTemplate(templateId)
	if err != nil {
		log.Fatalf("error while getting resource pools: %v", err)
	}

	deploymentName := "Agent-Deploy-Test-kalibuntu"

	log.Infof("Deploying %s [%s] as %s\n", template.GuestOS, template.Identifier, deploymentName)

	deploymentSpec := vsphere.DeployTemplateSpec{
		Description: "Kalibuntu Agent Deploy Test",
		Name:        deploymentName,
		DiskStorage: vsphere.TemplateDiskStorage{
			DatastoreIdentifier: datastore.Identifier,
		},
		GuestCustomization: vsphere.DeployGuestCustomization{
			Name: "LaForge-Test-Linux",
		},
		DiskStorageOverrides: map[string]vsphere.DeploySpecDiskStorage{},
		HardwareCustomization: vsphere.HardwareCustomization{
			CpuUpdate: vsphere.CpuUpdate{
				NumCoresPerSocket: 2,
				NumCpus:           4,
			},
			DisksToRemove: []string{},
			DisksToUpdate: map[string]vsphere.DiskUpdateSpec{},
			MemoryUpdate: vsphere.MemoryUpdate{
				Memory: 8192,
			},
			Nics: map[string]vsphere.HCNicUpdateSpec{
				"4000": {
					Identifier: network.Identifier,
				},
			},
		},
		Placement: vsphere.DeployPlacement{
			FolderId:       null.StringFrom(string(folder.Identifier)),
			ResourcePoolId: null.StringFrom(string(resourcePool.Identifier)),
		},
		PoweredOn: true,
		VmHomeStorage: vsphere.DeployHomeStorage{
			DatastoreId: datastore.Identifier,
		},
	}

	err = vs.DeployTemplate(templateId, deploymentSpec)
	if err != nil {
		log.Fatalf("error while deploying template \"%s\": %v", template.GuestOS, err)
	}

	log.Infof("Successcfully deployed VM \"%s\"!\n", deploymentSpec.Name)
}
