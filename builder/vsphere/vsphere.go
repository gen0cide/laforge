package vsphere

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vapi/library"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25"
)

type VSphere struct {
	BaseUrl string
	Client http.Client
	Username string
	Password string
	ContentLibraryId string
}

type PowerState string

const (
	POWER_STATE_ON = "POWERED_ON"
	POWER_STATE_OFF = "POWERED_OFF"
	POWER_STATE_SUSPENDED = "SUSPENDED"
)

type AuthorizationResponse struct {
	Value string `json:"value"`
}

type Identifier string

type VirtualMachine struct {
	Identifier Identifier `json:"vm"`
	MemorySize int `json:"memory-size_MiB"`
	Name string `json:"name"`
	PowerState PowerState `json:"power_state"`
	CpuCount int `json:"cpu_count"`
}

type VirtualMachineList struct {
	Value []VirtualMachine `json:"value"`
}

type Datastore struct {
	Identifier Identifier `json:"datastore"`
	Name string `json:"name"`
	Type string `json:"type"`
	FreeSpace int `json:"free_space"`
	Capacity int `json:"capacity"`
}

type DatastoreList struct {
	Value []Datastore `json:"value"`
}

type Folder struct {
	Identifier Identifier `json:"folder"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type FolderList struct {
	Value []Folder `json:"value"`
}

type ResourcePool struct {
	Identifier Identifier `json:"resource_pool"`
	Name string `json:"name"`
}

type ResourcePoolList struct {
	Value []ResourcePool `json:"value"`
}

type Memory struct {
	Size int `json:"size_MiB"`
}

type TemplateDiskStorage struct {
	DatastoreIdentifier string `json:"datastore"`
	StoragePolicyId string `json:"storage_policy"`
}

type TemplateDisk struct {
	DiskStorage TemplateDiskStorage `json:"disk_storage"`
	Capacity int `json:"capacity"`
}

type TemplateDiskEntry struct {
	Value TemplateDisk `json:"value"`
	Key string `json:"key"`
}

type MacType string
type BackingType string

const (
	NIC_MANUAL = "MANUAL"
	NIC_GENERATED = "GENERATED"
	NIC_ASSIGNED = "ASSIGNED"
	NIC_STANDARD_PORTGROUP = "STANDARD_PORTGROUP"
	NIC_HOST_DEVICE = "HOST_DEVICE"
	NIC_DISTRIBUTED_PORTGROUP = "DISTRIBUTED_PORTGROUP"
	NIC_OPAQUE_NETWORK = "OPAQUE_NETWORK"
)

type TemplateNic struct {
	Identifier Identifier `json:"network"`
	MacType MacType `json:"mac_type"`
	BackingType BackingType `json:"backing_type"`
}

type TemplateNicEntry struct {
	Value TemplateNic `json:"value"`
	Key string `json:"key"`
}

type TemplateCpu struct {
	Count int `json:"count"`
	CoresPerSocket int `json:"cores_per_socket"`
}

type TemplateHomeStorage struct {
	DatastoreIdentifier string `json:"datastore"`
	StoragePolicyId string `json:"storage_policy"`
}

type Template struct {
	Identifier Identifier `json:"vm_template"`
	Memory Memory `json:"memory"`
	Disks []TemplateDiskEntry `json:"disks"`
	Nics []TemplateNicEntry `json:"nics"`
	Cpu TemplateCpu `json:"cpu"`
	VmHomeStorage TemplateHomeStorage `json:"vm_home_storage"`
	GuestOS string `json:"guest_OS"`
}

type TemplateResponse struct {
	Value Template `json:"value"`
}

type OperatingSystem string

const (
	OS_UBUNTU_64 = "UBUNTU_64"
)

type VirtualMachinePlacement struct {
	Datastore Identifier `json:"datastore"`
	Folder Identifier `json:"folder"`
	ResourcePool Identifier `json:"resource_pool"`
}

type VirtualMachineCpu struct {
	Count int `json:"count"`
	CoresPerSocket int `json:"cores_per_socket"`
}

type VirtualMachineVmdk struct {
	Capacity int `json:"capacity"`
}

type VirtualMachineScsi struct {
	Bus int `json:"bus"`
	Unit int `json:"unit"`
}

type VirtualMachineDiskType string

const (
	VM_DISK_TYPE_SCSI VirtualMachineDiskType = "SCSI"
)

type VirtualMachineDisk struct {
	NewVdmk VirtualMachineVmdk `json:"new_vmdk"`
	Scsi VirtualMachineScsi `json:"scsi"`
	Type VirtualMachineDiskType `json:"type"`
}

type VirtualMachineNicType string

const (
	VM_NIC_E1000E VirtualMachineNicType = "E1000E"
)

type VirtualMachineNic struct {
	AllowGuestControl bool `json:"allow_guest_control"`
	StartConnected bool `json:"start_connected"`
	Type VirtualMachineNicType `json:"type"`
	UptCompatibilityEnabled bool `json:"upt_compatibility_enabled"`
	WakeOnLanEnabled bool `json:"wake_on_lan_enabled"`
}

type VirtualMachineCdromType string;

const (
	VM_CDROM_ISO_FILE = "ISO_FILE"
	VM_CDROM_HOST_DEVICE = "HOST_DEVICE"
	VM_CDROM_CLIENT_DEVICE = "CLIENT_DEVICE"
)

type VirtualMachineCdromBacking struct {
	DeviceAccessType string `json:"device_access_type"`
	Type VirtualMachineCdromType `json:"type"`
}

type VirtualMachineCdromSata struct {
	Bus int `json:"bus"`
	Unit int `json:"unit"`
}

type VirtualMachineCdrom struct {
	AllowGuestControl bool `json:"allow_guest_control"`
	Backing VirtualMachineCdromBacking `json:"backing"`
	Sata VirtualMachineCdromSata `json:"sata"`
	StartConnected bool `json:"start_connected"`
	Type string `json:"type"`
}

type VirtualMachineSpec struct {
	GuestOS OperatingSystem `json:"guest_OS"`
	Name string `json:"name"`
	Placement VirtualMachinePlacement `json:"placement"`
	Cpu VirtualMachineCpu `json:"cpu"`
	Memory Memory `json:"memory"`
	Disks []VirtualMachineDisk `json:"disks"`
	Nics []VirtualMachineNic `json:"nics"`
	Cdroms []VirtualMachineCdrom `json:"cdroms"`
}	

type CreateVirtualMachineData struct {
	Spec VirtualMachineSpec `json:"spec"`
}

type Network struct {
	Identifier Identifier `json:"network"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ListNetworkResponse struct {
	Value []Network `json:"value"`
}

func (vs *VSphere) authorize() (sessionToken string, err error) {
	authRequest, err := http.NewRequest("POST", vs.BaseUrl + "/rest/com/vmware/cis/session", nil)
	if err != nil {
		return
	}
	authRequest.SetBasicAuth(vs.Username, vs.Password)
	authResponse, err := vs.Client.Do(authRequest)
	if err != nil {
		return
	}
	if authResponse.Status != "200 OK" {
		err = errors.New("recieved status " + authResponse.Status)
		return
	}
	
	defer authResponse.Body.Close()

	var authBody AuthorizationResponse
	err = json.NewDecoder(authResponse.Body).Decode(&authBody)
	if err != nil {
		return
	}
	sessionToken = authBody.Value
	return
}

func (vs *VSphere) generateAuthorizedRequest(method string, url string) (request *http.Request, err error) {
	sessionToken, err := vs.authorize()
	if err != nil {
		return
	}
	request, err = http.NewRequest(method, (vs.BaseUrl + url), nil)
	if err != nil {
		return
	}
	request.Header.Add("vmware-api-session-id", sessionToken)
	return
}

func (vs *VSphere) generateAuthorizedRequestWithData(method string, url string, data *bytes.Buffer) (request *http.Request, err error) {
	sessionToken, err := vs.authorize()
	if err != nil {
		return
	}
	request, err = http.NewRequest(method, (vs.BaseUrl + url), data)
	if err != nil {
		return
	}
	request.Header.Add("vmware-api-session-id", sessionToken)
	return
}

func (vs *VSphere) ListVms() (vms []VirtualMachine, err error) {
	request, err := vs.generateAuthorizedRequest("GET", "/rest/vcenter/vm")
	if err != nil {
		return
	}
	response, err := vs.Client.Do(request)
	if err != nil {
		return
	}
	if response.Status != "200 OK" {
		err = errors.New("recieved status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	var vmList VirtualMachineList
	err = json.NewDecoder(response.Body).Decode(&vmList)
	if err != nil {
		return
	}
	vms = vmList.Value
	return
}

func (vs *VSphere) ListDatastores() (datastores []Datastore, err error) {
	request, err := vs.generateAuthorizedRequest("GET", "/rest/vcenter/datastore")
	if err != nil {
		return
	}
	response, err := vs.Client.Do(request)
	if err != nil {
		return
	}
	if response.Status != "200 OK" {
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	var datastoreList DatastoreList
	err = json.NewDecoder(response.Body).Decode(&datastoreList)
	if err != nil {
		return
	}
	datastores = datastoreList.Value
	return
}

func (vs *VSphere) ListFolders() (folders []Folder, err error) {
	request, err := vs.generateAuthorizedRequest("GET", "/rest/vcenter/folder")
	if err != nil {
		return
	}
	response, err := vs.Client.Do(request)
	if err != nil {
		return
	}
	if response.Status != "200 OK" {
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	var folderList FolderList
	err = json.NewDecoder(response.Body).Decode(&folderList)
	if err != nil {
		return
	}
	folders = folderList.Value
	return
}

func (vs *VSphere) ListResourcePools() (resourcePools []ResourcePool, err error) {
	request, err := vs.generateAuthorizedRequest("GET", "/rest/vcenter/resource-pool")
	if err != nil {
		return
	}
	response, err := vs.Client.Do(request)
	if err != nil {
		return
	}
	if response.Status != "200 OK" {
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	var resourcePoolList ResourcePoolList
	err = json.NewDecoder(response.Body).Decode(&resourcePoolList)
	if err != nil {
		return
	}
	resourcePools = resourcePoolList.Value
	return
}

func (vs *VSphere) ListNetworks() (networks []Network, err error) {
	request, err := vs.generateAuthorizedRequest("GET", "/rest/vcenter/network")
	if err != nil {
		return
	}
	response, err := vs.Client.Do(request)
	if err != nil {
		return
	}
	if response.Status != "200 OK" {
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	var networkList ListNetworkResponse
	err = json.NewDecoder(response.Body).Decode(&networkList)
	if err != nil {
		return
	}
	networks = networkList.Value
	return
}

func (vs *VSphere) GetTemplateIDByName(name string) (templateId string, err error) {
	u, err := url.Parse(vs.BaseUrl + "/sdk")
	if err != nil {
		return
	}
	u.User = url.UserPassword(vs.Username, vs.Password)
	
	ctx := context.TODO()
	client, err := govmomi.NewClient(ctx, u, false)
	if err != nil {
		return
	}

	v25, err := vim25.NewClient(ctx, client.RoundTripper)
	if err != nil {
		return
	}
	c := rest.NewClient(v25)
	c.Login(ctx, u.User)
	clm := library.NewManager(c)

	fi := library.FindItem{
		LibraryID: vs.ContentLibraryId,
		Name:      name,
	}
	items, err := clm.FindLibraryItems(ctx, fi)
	if err != nil {
		return
	}
	if len(items) < 1 {
		return
	}
	item, err := clm.GetLibraryItem(ctx, items[0])
	if err != nil {
		return
	}
	templateId = item.ID
	return
}

func (vs *VSphere) GetTemplate(templateId string) (template Template, err error) {
	request, err := vs.generateAuthorizedRequest("GET", "/rest/vcenter/vm-template/library-items/" + templateId)
	if err != nil {
		return
	}
	response, err := vs.Client.Do(request)
	if err != nil {
		return
	}
	if response.Status != "200 OK" {
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	var templateResponse TemplateResponse
	err = json.NewDecoder(response.Body).Decode(&templateResponse)
	if err != nil {
		return
	}
	template = templateResponse.Value
	return
}

func (vs *VSphere) CreateVM(vmSpec VirtualMachineSpec) (err error) {
	requestData := CreateVirtualMachineData{
		Spec: vmSpec,
	}
	requestDataString, err := json.Marshal(requestData)
	if err != nil {
		return
	}
	request, err := vs.generateAuthorizedRequestWithData("POST", "/rest/vcenter/vm", bytes.NewBuffer(requestDataString))
	if err != nil {
		return
	}
	response, err := vs.Client.Do(request)
	if err != nil {
		return
	}
	if response.Status != "200 OK" {
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	// TODO: Actually parse the response
	return
}