package main

func main() {
	panic("deprecated")
	// log.SetLevel(log.DebugLevel)
	// httpClient := http.Client{}

	// baseUrl, urlExists := os.LookupEnv("VSPHERE_URL")
	// username, usernameExists := os.LookupEnv("VSPHERE_USERNAME")
	// password, passwordExists := os.LookupEnv("VSPHERE_PASSWORD")
	// libraryName, libraryNameExists := os.LookupEnv("VSPHERE_CONTENT_LIBRARY")
	// templateName, templateNameExists := os.LookupEnv("VSPHERE_TEMPLATE_NAME")
	// datastoreName, datastoreNameExists := os.LookupEnv("VSPHERE_DATASTORE")
	// resourcePoolName, resourcePoolNameExists := os.LookupEnv("VSPHERE_RESOURCE_POOL")
	// folderName, folderNameExists := os.LookupEnv("VSPHERE_FOLDER")
	// networkName, networkNameExists := os.LookupEnv("VSPHERE_NETWORK")
	// if !urlExists || !usernameExists || !passwordExists || !templateNameExists || !libraryNameExists || !datastoreNameExists || !resourcePoolNameExists || !folderNameExists || !networkNameExists {
	// 	log.WithFields(log.Fields{
	// 		"VSPHERE_URL":             urlExists,
	// 		"VSPHERE_USERNAME":        usernameExists,
	// 		"VSPHERE_PASSWORD":        passwordExists,
	// 		"VSPHERE_CONTENT_LIBRARY": libraryNameExists,
	// 		"VSPHERE_TEMPLATE_NAME":   templateNameExists,
	// 		"VSPHERE_DATASTORE":       datastoreNameExists,
	// 		"VSPHERE_RESOURCE_POOL":   resourcePoolNameExists,
	// 		"VSPHERE_FOLDER":          folderNameExists,
	// 		"VSPHERE_NETWORK":         networkNameExists,
	// 	}).Fatal("some env variables are missing, these currently exist:")
	// }
	// vs := vsphere.VSphere{
	// 	HttpClient: httpClient,
	// 	BaseUrl:    baseUrl,
	// 	Username:   username,
	// 	Password:   password,
	// }

	// vsphere.InitializeGovmomi(&vs, baseUrl, username, password)

	// ctx := context.Background()

	// datastore, exists, err := vs.GetDatastoreByName(ctx, datastoreName)
	// if err != nil {
	// 	log.Fatalf("error while finding the datastore: %v", err)
	// }
	// if !exists {
	// 	log.Fatalf("datastore doesn't exist")
	// }
	// log.WithFields(log.Fields{
	// 	"name":       datastore.Name,
	// 	"identifier": datastore.Datastore.Value,
	// }).Info("Found datastore")

	// folder, err := vs.GetFolderByName(ctx, folderName)
	// if err != nil {
	// 	log.Fatalf("error while finding the folder: %v", err)
	// }
	// log.WithFields(log.Fields{
	// 	"name":       folder.Name,
	// 	"identifier": folder.Reference().Value,
	// }).Info("Found folder")

	// resourcePool, err := vs.GetResourcePoolByName(ctx, resourcePoolName)
	// if err != nil {
	// 	log.Fatalf("error while finding the resource pool: %v", err)
	// }
	// log.WithFields(log.Fields{
	// 	"name":       resourcePool.Name,
	// 	"identifier": resourcePool.Reference().Value,
	// }).Info("Found resource pool")

	// network, exists, err := vs.GetNetworkByName(ctx, networkName)
	// if err != nil {
	// 	log.Fatalf("error while finding the network: %v", err)
	// }
	// if !exists {
	// 	log.Fatalf("datastore doesn't exist")
	// }
	// log.WithFields(log.Fields{
	// 	"name":       network.Name,
	// 	"identifier": network.Network.Value,
	// }).Info("Found network")

	// // templateId, err := vs.GetTemplateIDByName(ctx, libraryName, templateName)
	// // if err != nil {
	// // 	log.Fatalf("error while searching content library for \"%s\": %v", templateName, err)
	// // }
	// template, err := vs.GetTemplate(ctx, libraryName, templateName)
	// if err != nil {
	// 	log.Fatalf("error while getting template: %v", err)
	// }
	// log.WithFields(log.Fields{
	// 	"name":       templateName,
	// 	"identifier": template.ID,
	// }).Info("Found template")

	// deploymentName := "Agent-Deploy-Test-kalibuntu"

	// log.Infof("Deploying %s [%s] as %s\n", template.Name, template.ID, deploymentName)

	// deploymentSpec := vsphere.DeployTemplateSpec{
	// 	Description: "Kalibuntu Agent Deploy Test",
	// 	Name:        deploymentName,
	// 	DiskStorage: vsphere.TemplateDiskStorage{
	// 		DatastoreIdentifier: vsphere.Identifier(datastore.Datastore.Value),
	// 	},
	// 	GuestCustomization: vsphere.DeployGuestCustomization{
	// 		Name: "LaForge-Test-Linux",
	// 	},
	// 	DiskStorageOverrides: map[string]vsphere.DeploySpecDiskStorage{},
	// 	HardwareCustomization: vsphere.HardwareCustomization{
	// 		CpuUpdate: vsphere.CpuUpdate{
	// 			NumCoresPerSocket: 2,
	// 			NumCpus:           4,
	// 		},
	// 		DisksToRemove: []string{},
	// 		DisksToUpdate: map[string]vsphere.DiskUpdateSpec{},
	// 		MemoryUpdate: vsphere.MemoryUpdate{
	// 			Memory: 8192,
	// 		},
	// 		Nics: map[string]vsphere.HCNicUpdateSpec{
	// 			"4000": {
	// 				Identifier: vsphere.Identifier(network.Network.Value),
	// 			},
	// 		},
	// 	},
	// 	Placement: vsphere.DeployPlacement{
	// 		FolderId:       null.StringFrom(string(folder.Reference().Value)),
	// 		ResourcePoolId: null.StringFrom(string(resourcePool.Reference().Value)),
	// 	},
	// 	PoweredOn: true,
	// 	VmHomeStorage: vsphere.DeployHomeStorage{
	// 		DatastoreId: vsphere.Identifier(datastore.Datastore.Value),
	// 	},
	// }

	// err = vs.DeployTemplate(template.ID, deploymentSpec)
	// if err != nil {
	// 	log.Fatalf("error while deploying template \"%s\": %v", template.Name, err)
	// }

	// log.Infof("Successcfully deployed VM \"%s\"!\n", deploymentSpec.Name)
}
