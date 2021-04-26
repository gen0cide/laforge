package vsphere

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vapi/library"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25"
	"gopkg.in/guregu/null.v4"
)

type VSphere struct {
	BaseUrl  string
	Client   http.Client
	Username string
	Password string
}

type PowerState string

const (
	POWER_STATE_ON        = "POWERED_ON"
	POWER_STATE_OFF       = "POWERED_OFF"
	POWER_STATE_SUSPENDED = "SUSPENDED"
)

type AuthorizationResponse struct {
	Value string `json:"value"`
}

type Identifier string

type VirtualMachine struct {
	Identifier Identifier `json:"vm"`
	MemorySize int        `json:"memory-size_MiB"`
	Name       string     `json:"name"`
	PowerState PowerState `json:"power_state"`
	CpuCount   int        `json:"cpu_count"`
}

type VirtualMachineList struct {
	Value []VirtualMachine `json:"value"`
}

type Datastore struct {
	Identifier Identifier `json:"datastore"`
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	FreeSpace  int        `json:"free_space"`
	Capacity   int        `json:"capacity"`
}

type DatastoreList struct {
	Value []Datastore `json:"value"`
}

type Folder struct {
	Identifier Identifier `json:"folder"`
	Name       string     `json:"name"`
	Type       string     `json:"type"`
}

type FolderList struct {
	Value []Folder `json:"value"`
}

type ResourcePool struct {
	Identifier Identifier `json:"resource_pool"`
	Name       string     `json:"name"`
}

type ResourcePoolList struct {
	Value []ResourcePool `json:"value"`
}

type Memory struct {
	Size int `json:"size_MiB"`
}

type TemplateDiskStorage struct {
	DatastoreIdentifier string      `json:"datastore"`
	StoragePolicyId     null.String `json:"storage_policy"`
}

type TemplateDisk struct {
	DiskStorage TemplateDiskStorage `json:"disk_storage"`
	Capacity    int                 `json:"capacity"`
}

type TemplateDiskEntry struct {
	Value TemplateDisk `json:"value"`
	Key   string       `json:"key"`
}

type MacType string
type BackingType string

const (
	NIC_MANUAL                = "MANUAL"
	NIC_GENERATED             = "GENERATED"
	NIC_ASSIGNED              = "ASSIGNED"
	NIC_STANDARD_PORTGROUP    = "STANDARD_PORTGROUP"
	NIC_HOST_DEVICE           = "HOST_DEVICE"
	NIC_DISTRIBUTED_PORTGROUP = "DISTRIBUTED_PORTGROUP"
	NIC_OPAQUE_NETWORK        = "OPAQUE_NETWORK"
)

type TemplateNic struct {
	Identifier  Identifier  `json:"network"`
	MacType     MacType     `json:"mac_type"`
	BackingType BackingType `json:"backing_type"`
}

type TemplateNicEntry struct {
	Value TemplateNic `json:"value"`
	Key   string      `json:"key"`
}

type TemplateCpu struct {
	Count          int `json:"count"`
	CoresPerSocket int `json:"cores_per_socket"`
}

type TemplateHomeStorage struct {
	DatastoreIdentifier string      `json:"datastore"`
	StoragePolicyId     null.String `json:"storage_policy"`
}

type Template struct {
	Identifier    Identifier          `json:"vm_template"`
	Memory        Memory              `json:"memory"`
	Disks         []TemplateDiskEntry `json:"disks"`
	Nics          []TemplateNicEntry  `json:"nics"`
	Cpu           TemplateCpu         `json:"cpu"`
	VmHomeStorage TemplateHomeStorage `json:"vm_home_storage"`
	GuestOS       string              `json:"guest_OS"`
}

type TemplateResponse struct {
	Value Template `json:"value"`
}

type OperatingSystem string

const (
	OS_UBUNTU_64 = "UBUNTU_64"
)

type VirtualMachinePlacement struct {
	Datastore    Identifier `json:"datastore"`
	Folder       Identifier `json:"folder"`
	ResourcePool Identifier `json:"resource_pool"`
}

type VirtualMachineCpu struct {
	Count          int `json:"count"`
	CoresPerSocket int `json:"cores_per_socket"`
}

type VirtualMachineVmdk struct {
	Capacity int `json:"capacity"`
}

type VirtualMachineScsi struct {
	Bus  int `json:"bus"`
	Unit int `json:"unit"`
}

type VirtualMachineDiskType string

const (
	VM_DISK_TYPE_SCSI VirtualMachineDiskType = "SCSI"
)

type VirtualMachineDisk struct {
	NewVdmk VirtualMachineVmdk     `json:"new_vmdk"`
	Scsi    VirtualMachineScsi     `json:"scsi"`
	Type    VirtualMachineDiskType `json:"type"`
}

type VirtualMachineNicType string

const (
	VM_NIC_E1000E VirtualMachineNicType = "E1000E"
)

type VirtualMachineNic struct {
	AllowGuestControl       bool                  `json:"allow_guest_control"`
	StartConnected          bool                  `json:"start_connected"`
	Type                    VirtualMachineNicType `json:"type"`
	UptCompatibilityEnabled bool                  `json:"upt_compatibility_enabled"`
	WakeOnLanEnabled        bool                  `json:"wake_on_lan_enabled"`
}

type VirtualMachineCdromType string

const (
	VM_CDROM_ISO_FILE      = "ISO_FILE"
	VM_CDROM_HOST_DEVICE   = "HOST_DEVICE"
	VM_CDROM_CLIENT_DEVICE = "CLIENT_DEVICE"
)

type VirtualMachineCdromBacking struct {
	DeviceAccessType string                  `json:"device_access_type"`
	Type             VirtualMachineCdromType `json:"type"`
}

type VirtualMachineCdromSata struct {
	Bus  int `json:"bus"`
	Unit int `json:"unit"`
}

type VirtualMachineCdrom struct {
	AllowGuestControl bool                       `json:"allow_guest_control"`
	Backing           VirtualMachineCdromBacking `json:"backing"`
	Sata              VirtualMachineCdromSata    `json:"sata"`
	StartConnected    bool                       `json:"start_connected"`
	Type              string                     `json:"type"`
}

type VirtualMachineSpec struct {
	GuestOS   OperatingSystem         `json:"guest_OS"`
	Name      string                  `json:"name"`
	Placement VirtualMachinePlacement `json:"placement"`
	Cpu       VirtualMachineCpu       `json:"cpu"`
	Memory    Memory                  `json:"memory"`
	Disks     []VirtualMachineDisk    `json:"disks"`
	Nics      []VirtualMachineNic     `json:"nics"`
	Cdroms    []VirtualMachineCdrom   `json:"cdroms"`
}

type CreateVirtualMachineData struct {
	Spec VirtualMachineSpec `json:"spec"`
}

type Network struct {
	Identifier Identifier `json:"network"`
	Name       string     `json:"name"`
	Type       string     `json:"type"`
}

type ListNetworkResponse struct {
	Value []Network `json:"value"`
}

type CpuUpdate struct {
	NumCoresPerSocket int `json:"num_cores_per_socket"`
	NumCpus           int `json:"num_cpus"`
}

type MemoryUpdate struct {
	Memory int `json:"memory"`
}

type HCNicValue struct {
	Identifier string `json:"network"`
}

type HCNic struct {
	Key   string     `json:"key"`
	Value HCNicValue `json:"value"`
}

type HardwareCustomization struct {
	CpuUpdate     CpuUpdate    `json:"cpu_update"`
	DisksToRemove []string     `json:"disks_to_remove"`
	DisksToUpdate []string     `json:"disks_to_update"`
	MemoryUpdate  MemoryUpdate `json:"memory_update"`
	Nics          []HCNic      `json:"nics"`
}

type DeployPlacement struct {
	ClusterId      null.String `json:"cluster"`
	FolderId       null.String `json:"folder"`
	HostId         null.String `json:"host"`
	ResourcePoolId null.String `json:"resource_pool"`
}

type DeployHomeStorage struct {
	DatastoreId     Identifier  `json:"datastore"`
	StoragePolicyId null.String `json:"storage_policy"`
}

type DeployTemplateSpec struct {
	Description           string                `json:"description"`
	DiskStorage           TemplateDiskStorage   `json:"disk_storage"`
	DiskStorageOverrides  []string              `json:"disk_storage_overrides"`
	HardwareCustomization HardwareCustomization `json:"hardware_customization"`
	Name                  string                `json:"name"`
	Placement             DeployPlacement       `json:"placement"`
	PoweredOn             bool                  `json:"powered_on"`
	VmHomeStorage         DeployHomeStorage     `json:"vm_home_storage"`
}

type DeployTemplatePayload struct {
	Spec DeployTemplateSpec `json:"spec"`
}

func (vs *VSphere) authorize() (sessionToken string, err error) {
	authRequest, err := http.NewRequest("POST", vs.BaseUrl+"/rest/com/vmware/cis/session", nil)
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
	request.Header.Add("Content-Type", "application/json")
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

func (vs *VSphere) GetDatastoreByName(name string) (datastore Datastore, err error) {
	request, err := vs.generateAuthorizedRequest("GET", "/rest/vcenter/datastore?filter.names.1="+name)
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
	if len(datastoreList.Value) < 1 {
		err = errors.New("no datastore found for the name\"" + name + "\"")
		return
	}
	if len(datastoreList.Value) > 1 {
		err = errors.New("more than one (" + fmt.Sprint(len(datastoreList.Value)) + ") datastore found for the name\"" + name + "\"")
		return
	}
	datastore = datastoreList.Value[0]
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

func (vs *VSphere) GetFolderByName(name string) (folder Folder, err error) {
	request, err := vs.generateAuthorizedRequest("GET", "/rest/vcenter/folder?filter.names.1="+name)
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
	if len(folderList.Value) < 1 {
		err = errors.New("no folder found for the name\"" + name + "\"")
		return
	}
	if len(folderList.Value) > 1 {
		err = errors.New("more than one (" + fmt.Sprint(len(folderList.Value)) + ") folder found for the name\"" + name + "\"")
		return
	}
	folder = folderList.Value[0]
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

func (vs *VSphere) GetResourcePoolByName(name string) (resourcePool ResourcePool, err error) {
	request, err := vs.generateAuthorizedRequest("GET", "/rest/vcenter/resource-pool?filter.names.1="+name)
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
	if len(resourcePoolList.Value) < 1 {
		err = errors.New("no resource pool found for the name\"" + name + "\"")
		return
	}
	if len(resourcePoolList.Value) > 1 {
		err = errors.New("more than one (" + fmt.Sprint(len(resourcePoolList.Value)) + ") resource pool found for the name\"" + name + "\"")
		return
	}
	resourcePool = resourcePoolList.Value[0]
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

func (vs *VSphere) GetNetworkByName(name string) (network Network, err error) {
	request, err := vs.generateAuthorizedRequest("GET", "/rest/vcenter/network?filter.names.1="+name)
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
	if len(networkList.Value) < 1 {
		err = errors.New("no network found for the name\"" + name + "\"")
		return
	}
	if len(networkList.Value) > 1 {
		err = errors.New("more than one (" + fmt.Sprint(len(networkList.Value)) + ") network found for the name\"" + name + "\"")
		return
	}
	network = networkList.Value[0]
	return
}
func (vs *VSphere) GetTemplateIDByName(contentLibraryName string, templateName string) (templateId string, err error) {
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

	// Find the content library by name
	cl := library.Find{
		Name: contentLibraryName,
	}
	contentLibrary, err := clm.FindLibrary(ctx, cl)
	if err != nil {
		return
	}
	if len(contentLibrary) < 1 {
		err = errors.New("no content libaries found with the name \"" + contentLibraryName + "\"")
		return
	}
	if len(contentLibrary) > 1 {
		err = errors.New("more than one content library matches the name \"" + contentLibraryName + "\"")
		return
	}

	fi := library.FindItem{
		LibraryID: contentLibrary[0],
		Name:      templateName,
	}
	items, err := clm.FindLibraryItems(ctx, fi)
	if err != nil {
		return
	}
	if len(items) < 1 {
		err = errors.New("no templates were found with the name \"" + templateName + "\"")
		return
	}
	if len(items) > 1 {
		err = errors.New("found more than one (" + fmt.Sprint(len(items)) + ") for the template \"" + templateName + "\"")
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
	request, err := vs.generateAuthorizedRequest("GET", "/rest/vcenter/vm-template/library-items/"+templateId)
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

func (vs *VSphere) DeployTemplate(templateId string, spec DeployTemplateSpec) (err error) {
	requestData := DeployTemplatePayload{
		Spec: spec,
	}
	requestDataString, err := json.Marshal(requestData)
	if err != nil {
		return
	}
	fmt.Printf("%s\n", requestDataString)
	request, err := vs.generateAuthorizedRequestWithData("POST", "/rest/vcenter/vm-template/library-items/"+templateId+"?action=deploy", bytes.NewBuffer(requestDataString))
	if err != nil {
		return
	}
	response, err := vs.Client.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	// DEBUG: output the deployment id for debug purposes
	defer response.Body.Close()
	deploymentBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	fmt.Printf("[DEBUG] | Deployed VM \"%s\" [%s]", spec.Name, string(deploymentBytes))

	return
}
