// Package vspherensxt provides an interface to deploy hosts and networks on VSphere and NSX-T
package vspherensxt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gen0cide/laforge/builder/vspherensxt/nsxt"
	"github.com/gen0cide/laforge/builder/vspherensxt/vsphere"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/logging"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
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
	VSphereDatastore          *types.DatastoreSummary
	VSphereResourcePool       *object.ResourcePool
	VSphereFolder             *object.Folder
	Logger                    *logging.Logger
}

func (builder VSphereNSXTBuilder) ID() string {
	return ID
}

func (builder VSphereNSXTBuilder) Name() string {
	return Name
}

func (builder VSphereNSXTBuilder) Description() string {
	return Description
}

func (builder VSphereNSXTBuilder) Author() string {
	return Author
}

func (builder VSphereNSXTBuilder) Version() string {
	return Version
}

func (builder VSphereNSXTBuilder) generateBuildID(build *ent.Build) string {
	buildId, err := build.ID.MarshalText()
	if err != nil {
		buildId = []byte(fmt.Sprint(build.Revision))
	}
	return fmt.Sprintf("%s", buildId)
}

func (builder VSphereNSXTBuilder) generateVmName(competition *ent.Competition, team *ent.Team, host *ent.Host, build *ent.Build) string {
	return (competition.HclID + "-Team-" + fmt.Sprintf("%02d", team.TeamNumber) + "-" + host.Hostname + "-" + builder.generateBuildID(build))
}

func (builder VSphereNSXTBuilder) generateRouterName(competition *ent.Competition, team *ent.Team, build *ent.Build) string {
	return (competition.HclID + "-Team-" + fmt.Sprintf("%02d", team.TeamNumber) + "-" + builder.generateBuildID(build))
}

func (builder VSphereNSXTBuilder) generateNetworkName(competition *ent.Competition, team *ent.Team, network *ent.Network, build *ent.Build) string {
	return (competition.HclID + "-Team-" + fmt.Sprintf("%02d", team.TeamNumber) + "-" + network.Name + "-" + builder.generateBuildID(build))
}

// DeployHost deploys a given host from the environment to VSphere
func (builder VSphereNSXTBuilder) DeployHost(ctx context.Context, provisionedHost *ent.ProvisionedHost) (err error) {
	host, err := provisionedHost.QueryProvisionedHostToHost().Only(ctx)
	if err != nil {
		return
	}
	cpuCount := int32(0)
	memorySize := int64(0)
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

	build, err := provisionedHost.QueryProvisionedHostToPlan().QueryPlanToBuild().Only(ctx)
	if err != nil {
		return
	}
	competition, err := build.QueryBuildToCompetition().Only(ctx)
	if err != nil {
		return
	}
	network, err := provisionedHost.QueryProvisionedHostToProvisionedNetwork().QueryProvisionedNetworkToNetwork().Only(ctx)
	if err != nil {
		return
	}
	team, err := provisionedHost.QueryProvisionedHostToProvisionedNetwork().QueryProvisionedNetworkToTeam().Only(ctx)
	if err != nil {
		return
	}

	networkName := builder.generateNetworkName(competition, team, network, build)

	// nsxtNetwork, err := builder.VSphereClient.GetNetworkByName(networkName)
	// if err != nil {
	// 	return
	// }

	// templateId, err := builder.VSphereClient.GetTemplateIDByName(builder.VSphereContentLibraryName, )
	// if err != nil {
	// 	return
	// }
	// template, err := builder.VSphereClient.GetTemplate(templateId)
	// if err != nil {
	// 	return
	// }

	// nicId, exists := host.Vars["nic_id"]
	// if !exists {
	// 	if len(template.Nics) > 0 {
	// 		nicIds := make([]string, 0, len(template.Nics))
	// 		for key := range template.Nics {
	// 			nicIds = append(nicIds, key)
	// 		}
	// 		if len(nicIds) < 1 {
	// 			return fmt.Errorf("no nics exist in the vm template for template %s", (builder.TemplatePrefix + host.OS))
	// 		}
	// 		nicId = nicIds[0]
	// 	} else {
	// 		err = errors.New("nic_id doesn't exist in the host vars for " + host.Hostname)
	// 		return err
	// 	}
	// }

	templateName := (builder.TemplatePrefix + host.OS)
	vmName := builder.generateVmName(competition, team, host, build)
	guestCustomizationName := (vmName + "-Customization-Spec")

	var guestCustomizationSpec *types.CustomizationSpecItem

	// specAlreadyExists, err := builder.VSphereClient.GuestCustomizationExists(ctx, guestCustomizationName)
	// if err != nil {
	// 	return
	// }
	// if !specAlreadyExists {
	guestCustomizationSpec, err = builder.VSphereClient.GenerateGuestCustomization(ctx, guestCustomizationName, templateName, provisionedHost)
	if err != nil {
		return err
	}

	err = builder.VSphereClient.CreateGuestCustomization(ctx, *guestCustomizationSpec)
	if err != nil {
		return err
	}
	// }

	// logrus.Debug(guestCustomizationSpec)

	// templateSpec := vsphere.DeployTemplateSpec{
	// 	Description: host.Description,
	// 	DiskStorage: vsphere.TemplateDiskStorage{
	// 		DatastoreIdentifier: builder.VSphereDatastore.Identifier,
	// 	},
	// 	DiskStorageOverrides: map[string]vsphere.DeploySpecDiskStorage{},
	// 	GuestCustomization: vsphere.DeployGuestCustomization{
	// 		Name: guestCustomizationName,
	// 	},
	// 	HardwareCustomization: vsphere.HardwareCustomization{
	// 		CpuUpdate: vsphere.CpuUpdate{
	// 			NumCoresPerSocket: 2,
	// 			NumCpus:           cpuCount,
	// 		},
	// 		DisksToRemove: []string{},
	// 		DisksToUpdate: map[string]vsphere.DiskUpdateSpec{},
	// 		MemoryUpdate: vsphere.MemoryUpdate{
	// 			Memory: memorySize,
	// 		},
	// 		Nics: map[string]vsphere.HCNicUpdateSpec{
	// 			nicId: {
	// 				Identifier: nsxtNetwork.Identifier,
	// 			},
	// 		},
	// 	},
	// 	Name: vmName,
	// 	Placement: vsphere.DeployPlacement{
	// 		ClusterId:      null.String{},
	// 		FolderId:       null.StringFrom(string(builder.VSphereFolder.Identifier)),
	// 		HostId:         null.String{},
	// 		ResourcePoolId: null.StringFrom(string(builder.VSphereResourcePool.Identifier)),
	// 	},
	// 	PoweredOn: true,
	// 	VmHomeStorage: vsphere.DeployHomeStorage{
	// 		DatastoreId: builder.VSphereDatastore.Identifier,
	// 	},
	// }

	// err = builder.VSphereClient.DeployTemplate(templateId, templateSpec)
	err = builder.VSphereClient.DeployLinkedClone(ctx, templateName, vmName, networkName, cpuCount, memorySize, builder.VSphereFolder, builder.VSphereResourcePool, guestCustomizationSpec)
	if err != nil {
		return
	}
	return
}

func (builder VSphereNSXTBuilder) DeployNetwork(ctx context.Context, provisionedNetwork *ent.ProvisionedNetwork) (err error) {
	build, err := provisionedNetwork.QueryProvisionedNetworkToBuild().Only(ctx)
	if err != nil {
		return fmt.Errorf("couldn't query build from network \"%s\": %v", provisionedNetwork.Name, err)
	}
	environment, err := provisionedNetwork.QueryProvisionedNetworkToBuild().QueryBuildToEnvironment().Only(ctx)
	if err != nil {
		return fmt.Errorf("couldn't query build from network \"%s\": %v", provisionedNetwork.Name, err)
	}
	competition, err := environment.EnvironmentToCompetition(ctx)
	if err != nil {
		return fmt.Errorf("couldn't query build from environment \"%s\": %v", environment.Name, err)
	}
	network, err := provisionedNetwork.QueryProvisionedNetworkToNetwork().Only(ctx)
	if err != nil {
		return fmt.Errorf("couldn't query build from network \"%s\": %v", provisionedNetwork.Name, err)
	}
	team, err := provisionedNetwork.QueryProvisionedNetworkToTeam().Only(ctx)
	if err != nil {
		return fmt.Errorf("couldn't query build from network \"%s\": %v", provisionedNetwork.Name, err)
	}

	tier1Name := builder.generateRouterName(competition[0], team, build)
	tier1Exists, nsxtError, err := builder.NsxtClient.CheckExistsTier1(tier1Name)
	if err != nil {
		return
	}
	if nsxtError != nil {
		return fmt.Errorf("nsx-t error %s (%d): %s", nsxtError.HttpStatus, nsxtError.ErrorCode, nsxtError.Message)
	}
	if !tier1Exists {
		// ------------------------------------
		// Stagger the teams network deployment
		time.Sleep(time.Duration(team.TeamNumber*5) * time.Second)
		// ------------------------------------

		tier0Path, tier0PathExists := network.Vars["nsxt_tier0_path"]
		if !tier0PathExists {
			tier0s, nsxtError, err := builder.NsxtClient.GetTier0s()
			if err != nil {
				return err
			}
			if nsxtError != nil {
				return fmt.Errorf("nsx-t error %s (%d): %s", nsxtError.HttpStatus, nsxtError.ErrorCode, nsxtError.Message)
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

		fmt.Println("Tier-1 router not found for Team " + fmt.Sprint(team.TeamNumber) + ", creating one...")
		nsxtError, err = builder.NsxtClient.CreateTier1(tier1Name, tier0Path, edgeClusterPath)
		if err != nil {
			return err
		}
		if nsxtError != nil {
			return fmt.Errorf("nsx-t error %s (%d): %s", nsxtError.HttpStatus, nsxtError.ErrorCode, nsxtError.Message)
		}

		addressParts := strings.Split(network.Cidr, "/")
		var natSourceAddress string

		switch addressParts[1] {
		case "24":
			octets := strings.Split(addressParts[0], ".")
			natSourceAddress = octets[0] + "." + octets[1] + ".0.0/16"
		default:
			octets := strings.Split(addressParts[0], ".")
			natSourceAddress = octets[0] + ".0.0.0/8"
		}

		ipPoolSubnets, nsxtError, err := builder.NsxtClient.GetIpPoolSubnets(builder.NsxtClient.IpPoolName)
		if err != nil {
			return err
		}
		if nsxtError != nil {
			return fmt.Errorf("nsx-t error %s (%d): %s", nsxtError.HttpStatus, nsxtError.ErrorCode, nsxtError.Message)
		}
		if len(ipPoolSubnets) <= 0 {
			return fmt.Errorf("error: no ip subnets found under the IP Pool \"%s\"", builder.NsxtClient.IpPoolName)
		}

		startingIp := ipPoolSubnets[0].AllocationRanges[0].Start
		endingIp := ipPoolSubnets[0].AllocationRanges[0].End
		octets := strings.Split(string(startingIp), ".")
		endOctets := strings.Split(string(endingIp), ".")
		lastOctet, err := strconv.Atoi(octets[3])
		if err != nil {
			return fmt.Errorf("error while reading last octet: %v", err)
		}
		endLastOctet, err := strconv.Atoi(endOctets[3])
		if err != nil {
			return fmt.Errorf("error while reading last endoctet: %v", err)
		}
		lastOctet = lastOctet + team.TeamNumber
		octets[3] = strconv.Itoa(lastOctet)
		natTranslatedAddress := strings.Join(octets, ".")
		if lastOctet > endLastOctet {
			return fmt.Errorf("NAT IP %s is out of the range %s-%s", natTranslatedAddress, startingIp, endingIp)
		}

		nsxtError, err = builder.NsxtClient.CreateNATRule(tier1Name, nsxt.NSXTIPElementList(natSourceAddress), nsxt.NSXTIPElementList(natTranslatedAddress))
		if err != nil {
			return err
		}
		if nsxtError != nil {
			return fmt.Errorf("nsx-t error %s (%d): %s", nsxtError.HttpStatus, nsxtError.ErrorCode, nsxtError.Message)
		}
	}

	networkName := builder.generateNetworkName(competition[0], team, network, build)
	cidrParts := strings.Split(network.Cidr, "/")
	gatewayAddress, gatewayAddressExists := network.Vars["gateway_address"]

	var octets []string
	if gatewayAddressExists {
		octets = strings.Split(gatewayAddress, ".")
	} else {
		octets = strings.Split(cidrParts[0], ".")
	}

	switch cidrParts[1] {
	case "24":
		octets[3] = "254"
		gatewayAddress = (strings.Join(octets, ".") + "/" + cidrParts[1])
	}

	// fmt.Println("deploying segment \"" + networkName + "\" w/ gateway_addr = " + gatewayAddress)
	nsxtError, err = builder.NsxtClient.CreateSegment(networkName, ("/infra/tier-1s/" + tier1Name), gatewayAddress)
	if nsxtError != nil {
		return fmt.Errorf("nsx-t error %s (%d): %s", nsxtError.HttpStatus, nsxtError.ErrorCode, nsxtError.Message)
	}
	return
}

func (builder VSphereNSXTBuilder) TeardownHost(ctx context.Context, provisionedHost *ent.ProvisionedHost) (err error) {
	host, err := provisionedHost.QueryProvisionedHostToHost().Only(ctx)
	if err != nil {
		return fmt.Errorf("couldn't query host from provisioned host \"%s\": %v", provisionedHost.ID, err)
	}
	build, err := provisionedHost.QueryProvisionedHostToProvisionedNetwork().QueryProvisionedNetworkToBuild().Only(ctx)
	if err != nil {
		return fmt.Errorf("couldn't query build from provisioned host \"%s\": %v", provisionedHost.ID, err)
	}
	competition, err := build.QueryBuildToCompetition().Only(ctx)
	if err != nil {
		return fmt.Errorf("couldn't query competition from build \"%s\": %v", build.ID, err)
	}
	team, err := provisionedHost.QueryProvisionedHostToProvisionedNetwork().QueryProvisionedNetworkToTeam().Only(ctx)
	if err != nil {
		return fmt.Errorf("couldn't query build from network \"%s\": %v", provisionedHost.ID, err)
	}

	vmName := builder.generateVmName(competition, team, host, build)
	guestCustomizationName := (vmName + "-Customization-Spec")

	specExists, err := builder.VSphereClient.GuestCustomizationExists(ctx, guestCustomizationName)
	if err != nil {
		return fmt.Errorf("error while checking if guest customization spec exists: %v", err)
	}
	if specExists {
		err = builder.VSphereClient.DeleteGuestCustomization(ctx, guestCustomizationName)
		if err != nil {
			return fmt.Errorf("error while deleting guest customization spec: %v", err)
		}
	}

	_, vmExists, err := builder.VSphereClient.GetVmSummary(ctx, vmName)
	if err != nil {
		return fmt.Errorf("error while checking if vm exists: %v", err)
	}
	if vmExists {
		vsphereErr := builder.VSphereClient.DeleteVM(ctx, vmName)
		if err != nil {
			return fmt.Errorf("error while tearing down VM \"%s\": %v", vmName, vsphereErr)
		}
		// Let vSphere sync with NSX-T
		time.Sleep(1 * time.Minute)
	}
	return
}

func (builder VSphereNSXTBuilder) TeardownNetwork(ctx context.Context, provisionedNetwork *ent.ProvisionedNetwork) (err error) {
	build, err := provisionedNetwork.QueryProvisionedNetworkToBuild().Only(ctx)
	if err != nil {
		return fmt.Errorf("couldn't query build from network \"%s\": %v", provisionedNetwork.Name, err)
	}
	environment, err := provisionedNetwork.QueryProvisionedNetworkToBuild().QueryBuildToEnvironment().Only(ctx)
	if err != nil {
		return fmt.Errorf("couldn't query build from network \"%s\": %v", provisionedNetwork.Name, err)
	}
	competition, err := environment.EnvironmentToCompetition(ctx)
	if err != nil {
		return fmt.Errorf("couldn't query build from environment \"%s\": %v", environment.Name, err)
	}
	network, err := provisionedNetwork.QueryProvisionedNetworkToNetwork().Only(ctx)
	if err != nil {
		return fmt.Errorf("couldn't query build from network \"%s\": %v", provisionedNetwork.Name, err)
	}
	team, err := provisionedNetwork.QueryProvisionedNetworkToTeam().Only(ctx)
	if err != nil {
		return fmt.Errorf("couldn't query build from network \"%s\": %v", provisionedNetwork.Name, err)
	}

	// Teardown Segment
	networkName := builder.generateNetworkName(competition[0], team, network, build)
	nsxtError, err := builder.NsxtClient.DeleteSegment(networkName)
	if err != nil {
		return
	}
	// Gotta remove those VMs first
	if nsxtError != nil {
		if nsxtError.ErrorCode == nsxt.NSXT_Segment_Has_VMs {
			return fmt.Errorf("nsx-t error: please remove all VMs from segment \"%s\" before deletion", networkName)
		}
		return fmt.Errorf("nsx-t error %s (%d): %s", nsxtError.HttpStatus, nsxtError.ErrorCode, nsxtError.Message)
	}

	tier1Name := builder.generateRouterName(competition[0], team, build)

	// Remove NAT Rules
	nsxtError, err = builder.NsxtClient.DeleteNATRule(tier1Name)
	if err != nil {
		return
	}
	// Random NSX-T Error
	if nsxtError != nil {
		return fmt.Errorf("nsx-t error %s (%d): %s", nsxtError.HttpStatus, nsxtError.ErrorCode, nsxtError.Message)
	}

	// Try to teardown Tier-1
	nsxtError, err = builder.NsxtClient.DeleteTier1(tier1Name)
	if err != nil {
		return
	}
	// If this wasn't the last segment on the Tier 1, don't worry about it
	if nsxtError != nil && nsxtError.ErrorCode != nsxt.NSXT_Tier1_Has_Children {
		return fmt.Errorf("nsx-t error %s (%d): %s", nsxtError.HttpStatus, nsxtError.ErrorCode, nsxtError.Message)
	}
	return
}
