// Package vspherensxt provides an interface to deploy hosts and networks on VSphere and NSX-T
package vspherensxt

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gen0cide/laforge/builder/vspherensxt/nsxt"
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
	NsxtClient                nsxt.NSXTClient
	VSphereClient             vsphere.VSphere
	TemplatePrefix            string
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
func (builder *VSphereNSXTBuilder) DeployHost(ctx context.Context, provisionedHost *ent.ProvisionedHost) (err error) {
	host := provisionedHost.HCLProvisionedHostToHost
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

func (builder *VSphereNSXTBuilder) DeployNetwork(ctx context.Context, provisionedNetwork *ent.ProvisionedNetwork) (err error) {
	build := provisionedNetwork.HCLProvisionedNetworkToBuild
	environment := provisionedNetwork.HCLProvisionedNetworkToBuild.HCLBuildToEnvironment
	network := provisionedNetwork.HCLProvisionedNetworkToNetwork
	team := provisionedNetwork.HCLProvisionedNetworkToTeam

	tier1Name := (environment.Name + "-Team-" + fmt.Sprint(team.TeamNumber) + "-V" + fmt.Sprint(build.Revision))
	tier1Exists, err := builder.NsxtClient.CheckExistsTier1(tier1Name)
	if err != nil {
		return
	}
	if !tier1Exists {
		tier0Path, tier0PathExists := network.Vars["nsxt_tier0_path"]
		if !tier0PathExists {
			tier0s, err := builder.NsxtClient.GetTier0s()
			if err != nil {
				return err
			}
			if len(tier0s) > 1 {
				return errors.New("tier0_path doesn't exist in the network vars for " + network.Name + " and multiple (" + fmt.Sprint(len(tier0s)) + ") Tier 0s found.")
			}
			if len(tier0s) == 0 {
				return errors.New("tier0_path doesn't exist in the network vars for " + network.Name + " and no Tier 0s found.")
			}
			tier0Path = tier0s[0].Path
		}

		edgeClusterPath, edgeClusterPathExists := network.Vars["nsxt_edge_cluster_path"]
		if !edgeClusterPathExists {
			return errors.New("nsxt_edge_cluster_path doesn't exist in the network vars for " + network.Name)
		}

		err = builder.NsxtClient.CreateTier1(tier1Name, tier0Path, edgeClusterPath)
		if err != nil {
			return err
		}
	}

	networkName := (environment.Name + "-Team-" + fmt.Sprint(team.TeamNumber) + network.Name + "-V" + fmt.Sprint(build.Revision))
	// cidrParts := strings.Split(network.Cidr, "/")
	// gatewayAddress := strings.Join(strings.Split(cidrParts[0], ".")[:2], ".") + ".254/" + cidrParts
	err = builder.NsxtClient.CreateNetworkSegment(networkName, ("/infra/tier-1s/" + tier1Name), network.Cidr)
	return
}
