package vsphere

import (
	"encoding/json"
	"errors"
	"net/http"
)

type VSphere struct {
	BaseUrl string
	Client http.Client
	Username string
	Password string
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

type VirtualMachine struct {
	Identifier string `json:"vm"`
	MemorySize int `json:"memory-size_MiB"`
	Name string `json:"name"`
	PowerState PowerState `json:"power_state"`
	CpuCount int `json:"cpu_count"`
}

type VirtualMachineList struct {
	Value []VirtualMachine `json:"value"`
}

type Datastore struct {
	Identifier string `json:"datastore"`
	Name string `json:"name"`
	Type string `json:"type"`
	FreeSpace int `json:"free_space"`
	Capacity int `json:"capacity"`
}

type DatastoreList struct {
	Value []Datastore `json:"value"`
}

type Folder struct {
	Identifier string `json:"folder"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type FolderList struct {
	Value []Folder `json:"value"`
}

type ResourcePool struct {
	Identifier string `json:"resource_pool"`
	Name string `json:"name"`
}

type ResourcePoolList struct {
	Value []ResourcePool `json:"value"`
}

type TemplateMemory struct {
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
	Identifier string `json:"network"`
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
	Identifier string `json:"vm_template"`
	Memory TemplateMemory `json:"memory"`
	Disks []TemplateDiskEntry `json:"disks"`
	Nics []TemplateNicEntry `json:"nics"`
	Cpu TemplateCpu `json:"cpu"`
	VmHomeStorage TemplateHomeStorage `json:"vm_home_storage"`
	GuestOS string `json:"guest_OS"`
}

type TemplateResponse struct {
	Value Template `json:"value"`
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