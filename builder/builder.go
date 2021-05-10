package builder

import (
	"context"
	"errors"
	"net/http"

	"github.com/gen0cide/laforge/builder/vspherensxt"
	"github.com/gen0cide/laforge/builder/vspherensxt/vsphere"
	"github.com/gen0cide/laforge/ent"
)

type Builder interface {
	ID() string
	Name() string
	Description() string
	Author() string
	DeployHost(ctx context.Context, host *ent.Host) (err error)
}

// NewVSphereNSXTBuilder creates a builder instance to deploy environments to VSphere and NSX-T
func NewVSphereNSXTBuilder(env *ent.Environment) (builder vspherensxt.VSphereNSXTBuilder, err error) {
	vsphereUsername, exists := env.Config["vsphere_username"]
	if !exists {
		err = errors.New("vsphere_username doesn't exist in the environment configuration")
		return
	}
	vspherePassword, exists := env.Config["vsphere_password"]
	if !exists {
		err = errors.New("vsphere_password doesn't exist in the environment configuration")
		return
	}
	vsphereBaseUrl, exists := env.Config["vsphere_base_url"]
	if !exists {
		err = errors.New("vsphere_base_url doesn't exist in the environment configuration")
		return
	}
	nsxtBaseUrl, exists := env.Config["nsxt_base_url"]
	if !exists {
		err = errors.New("nsxt_base_url doesn't exist in the environment configuration")
		return
	}
	contentLibraryName, exists := env.Config["vsphere_content_library"]
	if !exists {
		err = errors.New("vsphere_content_library doesn't exist in the environment configuration")
		return
	}
	datastoreName, exists := env.Config["vsphere_datastore"]
	if !exists {
		err = errors.New("vsphere_datastore doesn't exist in the environment configuration")
		return
	}
	resourcePoolName, exists := env.Config["vsphere_resource_pool"]
	if !exists {
		err = errors.New("vsphere_resource_pool doesn't exist in the environment configuration")
		return
	}
	folderName, exists := env.Config["vsphere_folder"]
	if !exists {
		err = errors.New("vsphere_folder doesn't exist in the environment configuration")
		return
	}
	templatePrefix, exists := env.Config["vsphere_template_prefix"]
	if !exists {
		err = errors.New("vsphere_template_prefix doesn't exist in the environment configuration")
		return
	}

	httpClient := http.Client{}

	client := vsphere.VSphere{
		HttpClient: httpClient,
		BaseUrl:    vsphereBaseUrl,
		Username:   builder.Username,
		Password:   builder.Password,
	}

	datastore, err := client.GetDatastoreByName(datastoreName)
	if err != nil {
		return
	}

	resourcePool, err := client.GetResourcePoolByName(resourcePoolName)
	if err != nil {
		return
	}

	folder, err := client.GetFolderByName(folderName)
	if err != nil {
		return
	}

	builder = vspherensxt.VSphereNSXTBuilder{
		HttpClient:                httpClient,
		Username:                  vsphereUsername,
		Password:                  vspherePassword,
		NsxtUrl:                   nsxtBaseUrl,
		TemplatePrefix:            templatePrefix,
		VSphereClient:             client,
		VSphereContentLibraryName: contentLibraryName,
		VSphereDatastore:          datastore,
		VSphereResourcePool:       resourcePool,
		VSphereFolder:             folder,
	}
	return
}
