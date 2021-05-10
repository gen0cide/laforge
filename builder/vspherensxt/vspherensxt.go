// Package vspherensxt provides an interface to deploy hosts and networks on VSphere and NSX-T
package vspherensxt

import (
	"context"
	"errors"
	"net/http"

	"github.com/gen0cide/laforge/builder/vspherensxt/vsphere"
	"github.com/gen0cide/laforge/ent"
	"gopkg.in/guregu/null.v4"
)

const (
	ID          = "vsphere-nsxt"
	Name        = "VSphere-NSXT"
	Description = "Builder that interfaces with VSphere VCenter and NSX-T"
	Author      = "Bradley harker <github.com/BradHacker>"
	Version     = "0.1"
)

type VSphereNSXTBuilder struct {
	HttpClient                http.Client
	Username                  string
	Password                  string
	NsxtUrl                   string
	TemplatePrefix            string
	VSphereClient             vsphere.VSphere
	VSphereContentLibraryName string
	VSphereDatastore          vsphere.Datastore
	VSphereResourcePool       vsphere.ResourcePool
	VSphereFolder             vsphere.Folder
}

func (builder *VSphereNSXTBuilder) ID() string {
	return ID
}

func (builder *VSphereNSXTBuilder) Name() string {
	return Name
}

func (builder *VSphereNSXTBuilder) Description() string {
	return Description
}

func (builder *VSphereNSXTBuilder) Author() string {
	return Author
}

// DeployHost deploys a given host from the environment to VSphere
func (builder *VSphereNSXTBuilder) DeployHost(ctx context.Context, host *ent.Host) (err error) {
	cpuCount := 0
	memorySize := 0
	switch host.InstanceSize {
	case "nano":
		cpuCount = 1
		memorySize = 1024
	case "micro":
		cpuCount = 1
		memorySize = 2048
	case "small":
		cpuCount = 2
		memorySize = 2048
	case "medium":
		cpuCount = 2
		memorySize = 4096
	case "large":
		cpuCount = 4
		memorySize = 4096
	case "xlarge":
		cpuCount = 4
		memorySize = 8192
	}
	if cpuCount == 0 || memorySize == 0 {
		err = errors.New("couldn't resolve host instance size: " + host.InstanceSize)
		return
	}

	nicId, exists := host.Vars["nic_id"]
	if !exists {
		err = errors.New("nic_id doesn't exist in the host vars for " + host.Hostname)
		return
	}

	includedNetwork, err := host.QueryHostToIncludedNetwork().Only(ctx)
	if err != nil {
		return
	}

	network, err := builder.VSphereClient.GetNetworkByName(includedNetwork.Name)
	if err != nil {
		return
	}

	templateSpec := vsphere.DeployTemplateSpec{
		Description: host.Description,
		DiskStorage: vsphere.TemplateDiskStorage{
			DatastoreIdentifier: builder.VSphereDatastore.Identifier,
		},
		DiskStorageOverrides: []string{},
		HardwareCustomization: vsphere.HardwareCustomization{
			CpuUpdate: vsphere.CpuUpdate{
				NumCoresPerSocket: 2,
				NumCpus:           cpuCount,
			},
			DisksToRemove: []string{},
			DisksToUpdate: []string{},
			MemoryUpdate: vsphere.MemoryUpdate{
				Memory: memorySize,
			},
			Nics: []vsphere.HCNic{
				{
					Key: nicId,
					Value: vsphere.HCNicValue{
						Identifier: network.Identifier,
					},
				},
			},
		},
		Name: host.Hostname,
		Placement: vsphere.DeployPlacement{
			ClusterId:      null.String{},
			FolderId:       null.StringFrom(string(builder.VSphereFolder.Identifier)),
			HostId:         null.String{},
			ResourcePoolId: null.StringFrom(string(builder.VSphereResourcePool.Identifier)),
		},
		PoweredOn: true,
		VmHomeStorage: vsphere.DeployHomeStorage{
			DatastoreId: builder.VSphereDatastore.Identifier,
		},
	}

	templateId, err := builder.VSphereClient.GetTemplateIDByName(builder.VSphereContentLibraryName, builder.TemplatePrefix+host.OS)
	if err != nil {
		return
	}

	err = builder.VSphereClient.DeployTemplate(templateId, templateSpec)
	if err != nil {
		return
	}
	return
}
