package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gen0cide/laforge/builder/vspherensxt/vsphere"
	"gopkg.in/guregu/null.v4"
)

func main() {
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
		log.Fatalf("please set VSPHERE_URL (exists? %t)"+
			", VSPHERE_USERNAME (exists? %t)"+
			", VSPHERE_PASSWORD (exists? %t)"+
			", VSPHERE_CONTENT_LIBRARY (exists? %t)"+
			", VSPHERE_TEMPLATE_NAME (exists? %t)"+
			", VSPHERE_DATASTORE (exists? %t)"+
			", VSPHERE_RESOURCE_POOL (exists? %t)"+
			", VSPHERE_FOLDER (exists? %t)"+
			", VSPHERE_NETWORK (exists? %t)",
			urlExists,
			usernameExists,
			passwordExists,
			libraryNameExists,
			templateNameExists,
			datastoreNameExists,
			resourcePoolNameExists,
			folderNameExists,
			networkNameExists)
	}
	vs := vsphere.VSphere{
		Client:   httpClient,
		BaseUrl:  baseUrl,
		Username: username,
		Password: password,
	}

	datastore, err := vs.GetDatastoreByName(datastoreName)
	if err != nil {
		log.Fatalf("error while finding the datastore: %v", err)
	}
	fmt.Printf("Found datastore\t\t%s [%s]\n", datastore.Name, datastore.Identifier)

	folder, err := vs.GetFolderByName(folderName)
	if err != nil {
		log.Fatalf("error while finding the folder: %v", err)
	}
	fmt.Printf("Found folder\t\t%s [%s]\n", folder.Name, folder.Identifier)

	resourcePool, err := vs.GetResourcePoolByName(resourcePoolName)
	if err != nil {
		log.Fatalf("error while finding the resource pool: %v", err)
	}
	fmt.Printf("Found resource pool\t%s [%s]\n", resourcePool.Name, resourcePool.Identifier)

	network, err := vs.GetNetworkByName(networkName)
	if err != nil {
		log.Fatalf("error while finding the network: %v", err)
	}
	fmt.Printf("Found network\t\t%s [%s]\n", network.Name, network.Identifier)

	templateId, err := vs.GetTemplateIDByName(libraryName, templateName)
	if err != nil {
		log.Fatalf("error while searching content library for \"%s\": %v", templateName, err)
	}
	fmt.Printf("Found template\t\t%s [%s]\n", templateName, templateId)

	template, err := vs.GetTemplate(templateId)
	if err != nil {
		log.Fatalf("error while getting resource pools: %v", err)
	}

	deploymentSpec := vsphere.DeployTemplateSpec{
		Description: "Test VM created from Golang",
		Name:        "Builder-Test",
		DiskStorage: vsphere.TemplateDiskStorage{
			DatastoreIdentifier: datastore.Identifier,
		},
		DiskStorageOverrides: []string{},
		HardwareCustomization: vsphere.HardwareCustomization{
			CpuUpdate: vsphere.CpuUpdate{
				NumCoresPerSocket: 2,
				NumCpus:           4,
			},
			DisksToRemove: []string{},
			DisksToUpdate: []string{},
			MemoryUpdate: vsphere.MemoryUpdate{
				Memory: 4096,
			},
			Nics: []vsphere.HCNic{
				{
					Key: "4000",
					Value: vsphere.HCNicValue{
						Identifier: network.Identifier,
					},
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

	fmt.Printf("Deploying %s [%s] as %s\n", template.GuestOS, template.Identifier, deploymentSpec.Name)

	err = vs.DeployTemplate(templateId, deploymentSpec)
	if err != nil {
		log.Fatalf("error while deploying template \""+template.GuestOS+"\": %v", err)
	}

	fmt.Printf("Successcfully deployed VM \"" + deploymentSpec.Name + "\"!")
}
