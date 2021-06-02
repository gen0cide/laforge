// Package vpshere is for interfaceing with the VSphere REST API
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
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gen0cide/laforge/ent"
	log "github.com/sirupsen/logrus"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vapi/library"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25"
	"gopkg.in/guregu/null.v4"
)

type VSphere struct {
	BaseUrl    string
	ServerUrl  string
	HttpClient http.Client
	Username   string
	Password   string
}

type PowerState string

const (
	POWER_STATE_ON        = "POWERED_ON"
	POWER_STATE_OFF       = "POWERED_OFF"
	POWER_STATE_SUSPENDED = "SUSPENDED"
)

type Identifier string

type VirtualMachine struct {
	Identifier Identifier `json:"vm"`
	MemorySize int        `json:"memory-size_MiB"`
	Name       string     `json:"name"`
	PowerState PowerState `json:"power_state"`
	CpuCount   int        `json:"cpu_count"`
}

type Datastore struct {
	Identifier Identifier `json:"datastore"`
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	FreeSpace  int        `json:"free_space"`
	Capacity   int        `json:"capacity"`
}

type Folder struct {
	Identifier Identifier `json:"folder"`
	Name       string     `json:"name"`
	Type       string     `json:"type"`
}

type ResourcePool struct {
	Identifier Identifier `json:"resource_pool"`
	Name       string     `json:"name"`
}

type Memory struct {
	Size int `json:"size_MiB"`
}

type TemplateDiskStorage struct {
	DatastoreIdentifier Identifier  `json:"datastore"`
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
	DatastoreIdentifier Identifier  `json:"datastore"`
	StoragePolicyId     null.String `json:"storage_policy"`
}

type VmGuestOs string

const (
	OS_DOS                   VmGuestOs = "DOS"
	OS_WIN_31                VmGuestOs = "WIN_31"
	OS_WIN_95                VmGuestOs = "WIN_95"
	OS_WIN_98                VmGuestOs = "WIN_98"
	OS_WIN_ME                VmGuestOs = "WIN_ME"
	OS_WIN_NT                VmGuestOs = "WIN_NT"
	OS_WIN_2000_PRO          VmGuestOs = "WIN_2000_PRO"
	OS_WIN_2000_SERV         VmGuestOs = "WIN_2000_SERV"
	OS_WIN_2000_ADV_SERV     VmGuestOs = "WIN_2000_ADV_SERV"
	OS_WIN_XP_HOME           VmGuestOs = "WIN_XP_HOME"
	OS_WIN_XP_PRO            VmGuestOs = "WIN_XP_PRO"
	OS_WIN_XP_PRO_64         VmGuestOs = "WIN_XP_PRO_64"
	OS_WIN_NET_WEB           VmGuestOs = "WIN_NET_WEB"
	OS_WIN_NET_STANDARD      VmGuestOs = "WIN_NET_STANDARD"
	OS_WIN_NET_ENTERPRISE    VmGuestOs = "WIN_NET_ENTERPRISE"
	OS_WIN_NET_DATACENTER    VmGuestOs = "WIN_NET_DATACENTER"
	OS_WIN_NET_BUSINESS      VmGuestOs = "WIN_NET_BUSINESS"
	OS_WIN_NET_STANDARD_64   VmGuestOs = "WIN_NET_STANDARD_64"
	OS_WIN_NET_ENTERPRISE_64 VmGuestOs = "WIN_NET_ENTERPRISE_64"
	OS_WIN_LONGHORN          VmGuestOs = "WIN_LONGHORN"
	OS_WIN_LONGHORN_64       VmGuestOs = "WIN_LONGHORN_64"
	OS_WIN_NET_DATACENTER_64 VmGuestOs = "WIN_NET_DATACENTER_64"
	OS_WIN_VISTA             VmGuestOs = "WIN_VISTA"
	OS_WIN_VISTA_64          VmGuestOs = "WIN_VISTA_64"
	OS_WINDOWS_7             VmGuestOs = "WINDOWS_7"
	OS_WINDOWS_7_64          VmGuestOs = "WINDOWS_7_64"
	OS_WINDOWS_7_SERVER_64   VmGuestOs = "WINDOWS_7_SERVER_64"
	OS_WINDOWS_8             VmGuestOs = "WINDOWS_8"
	OS_WINDOWS_8_64          VmGuestOs = "WINDOWS_8_64"
	OS_WINDOWS_8_SERVER_64   VmGuestOs = "WINDOWS_8_SERVER_64"
	OS_WINDOWS_9             VmGuestOs = "WINDOWS_9"
	OS_WINDOWS_9_64          VmGuestOs = "WINDOWS_9_64"
	OS_WINDOWS_9_SERVER_64   VmGuestOs = "WINDOWS_9_SERVER_64"
	OS_WINDOWS_HYPERV        VmGuestOs = "WINDOWS_HYPERV"
	OS_WINDOWS_SERVER_2019   VmGuestOs = "WINDOWS_SERVER_2019"
	OS_WINDOWS_SERVER_2021   VmGuestOs = "WINDOWS_SERVER_2021"
	OS_FREEBSD               VmGuestOs = "FREEBSD"
	OS_FREEBSD_64            VmGuestOs = "FREEBSD_64"
	OS_FREEBSD_11            VmGuestOs = "FREEBSD_11"
	OS_FREEBSD_12            VmGuestOs = "FREEBSD_12"
	OS_FREEBSD_13            VmGuestOs = "FREEBSD_13"
	OS_FREEBSD_11_64         VmGuestOs = "FREEBSD_11_64"
	OS_FREEBSD_12_64         VmGuestOs = "FREEBSD_12_64"
	OS_FREEBSD_13_64         VmGuestOs = "FREEBSD_13_64"
	OS_REDHAT                VmGuestOs = "REDHAT"
	OS_RHEL_2                VmGuestOs = "RHEL_2"
	OS_RHEL_3                VmGuestOs = "RHEL_3"
	OS_RHEL_3_64             VmGuestOs = "RHEL_3_64"
	OS_RHEL_4                VmGuestOs = "RHEL_4"
	OS_RHEL_4_64             VmGuestOs = "RHEL_4_64"
	OS_RHEL_5                VmGuestOs = "RHEL_5"
	OS_RHEL_5_64             VmGuestOs = "RHEL_5_64"
	OS_RHEL_6                VmGuestOs = "RHEL_6"
	OS_RHEL_6_64             VmGuestOs = "RHEL_6_64"
	OS_RHEL_7                VmGuestOs = "RHEL_7"
	OS_RHEL_7_64             VmGuestOs = "RHEL_7_64"
	OS_RHEL_8_64             VmGuestOs = "RHEL_8_64"
	OS_RHEL_9_64             VmGuestOs = "RHEL_9_64"
	OS_CENTOS                VmGuestOs = "CENTOS"
	OS_CENTOS_64             VmGuestOs = "CENTOS_64"
	OS_CENTOS_6              VmGuestOs = "CENTOS_6"
	OS_CENTOS_6_64           VmGuestOs = "CENTOS_6_64"
	OS_CENTOS_7              VmGuestOs = "CENTOS_7"
	OS_CENTOS_7_64           VmGuestOs = "CENTOS_7_64"
	OS_CENTOS_8_64           VmGuestOs = "CENTOS_8_64"
	OS_CENTOS_9_64           VmGuestOs = "CENTOS_9_64"
	OS_ORACLE_LINUX          VmGuestOs = "ORACLE_LINUX"
	OS_ORACLE_LINUX_64       VmGuestOs = "ORACLE_LINUX_64"
	OS_ORACLE_LINUX_6        VmGuestOs = "ORACLE_LINUX_6"
	OS_ORACLE_LINUX_6_64     VmGuestOs = "ORACLE_LINUX_6_64"
	OS_ORACLE_LINUX_7        VmGuestOs = "ORACLE_LINUX_7"
	OS_ORACLE_LINUX_7_64     VmGuestOs = "ORACLE_LINUX_7_64"
	OS_ORACLE_LINUX_8_64     VmGuestOs = "ORACLE_LINUX_8_64"
	OS_ORACLE_LINUX_9_64     VmGuestOs = "ORACLE_LINUX_9_64"
	OS_SUSE                  VmGuestOs = "SUSE"
	OS_SUSE_64               VmGuestOs = "SUSE_64"
	OS_SLES                  VmGuestOs = "SLES"
	OS_SLES_64               VmGuestOs = "SLES_64"
	OS_SLES_10               VmGuestOs = "SLES_10"
	OS_SLES_10_64            VmGuestOs = "SLES_10_64"
	OS_SLES_11               VmGuestOs = "SLES_11"
	OS_SLES_11_64            VmGuestOs = "SLES_11_64"
	OS_SLES_12               VmGuestOs = "SLES_12"
	OS_SLES_12_64            VmGuestOs = "SLES_12_64"
	OS_SLES_15_64            VmGuestOs = "SLES_15_64"
	OS_SLES_16_64            VmGuestOs = "SLES_16_64"
	OS_NLD_9                 VmGuestOs = "NLD_9"
	OS_OES                   VmGuestOs = "OES"
	OS_SJDS                  VmGuestOs = "SJDS"
	OS_MANDRAKE              VmGuestOs = "MANDRAKE"
	OS_MANDRIVA              VmGuestOs = "MANDRIVA"
	OS_MANDRIVA_64           VmGuestOs = "MANDRIVA_64"
	OS_TURBO_LINUX           VmGuestOs = "TURBO_LINUX"
	OS_TURBO_LINUX_64        VmGuestOs = "TURBO_LINUX_64"
	OS_UBUNTU                VmGuestOs = "UBUNTU"
	OS_UBUNTU_64             VmGuestOs = "UBUNTU_64"
	OS_DEBIAN_4              VmGuestOs = "DEBIAN_4"
	OS_DEBIAN_4_64           VmGuestOs = "DEBIAN_4_64"
	OS_DEBIAN_5              VmGuestOs = "DEBIAN_5"
	OS_DEBIAN_5_64           VmGuestOs = "DEBIAN_5_64"
	OS_DEBIAN_6              VmGuestOs = "DEBIAN_6"
	OS_DEBIAN_6_64           VmGuestOs = "DEBIAN_6_64"
	OS_DEBIAN_7              VmGuestOs = "DEBIAN_7"
	OS_DEBIAN_7_64           VmGuestOs = "DEBIAN_7_64"
	OS_DEBIAN_8              VmGuestOs = "DEBIAN_8"
	OS_DEBIAN_8_64           VmGuestOs = "DEBIAN_8_64"
	OS_DEBIAN_9              VmGuestOs = "DEBIAN_9"
	OS_DEBIAN_9_64           VmGuestOs = "DEBIAN_9_64"
	OS_DEBIAN_10             VmGuestOs = "DEBIAN_10"
	OS_DEBIAN_10_64          VmGuestOs = "DEBIAN_10_64"
	OS_DEBIAN_11             VmGuestOs = "DEBIAN_11"
	OS_DEBIAN_11_64          VmGuestOs = "DEBIAN_11_64"
	OS_ASIANUX_3             VmGuestOs = "ASIANUX_3"
	OS_ASIANUX_3_64          VmGuestOs = "ASIANUX_3_64"
	OS_ASIANUX_4             VmGuestOs = "ASIANUX_4"
	OS_ASIANUX_4_64          VmGuestOs = "ASIANUX_4_64"
	OS_ASIANUX_5_64          VmGuestOs = "ASIANUX_5_64"
	OS_ASIANUX_7_64          VmGuestOs = "ASIANUX_7_64"
	OS_ASIANUX_8_64          VmGuestOs = "ASIANUX_8_64"
	OS_ASIANUX_9_64          VmGuestOs = "ASIANUX_9_64"
	OS_OPENSUSE              VmGuestOs = "OPENSUSE"
	OS_OPENSUSE_64           VmGuestOs = "OPENSUSE_64"
	OS_FEDORA                VmGuestOs = "FEDORA"
	OS_FEDORA_64             VmGuestOs = "FEDORA_64"
	OS_COREOS_64             VmGuestOs = "COREOS_64"
	OS_VMWARE_PHOTON_64      VmGuestOs = "VMWARE_PHOTON_64"
	OS_OTHER_24X_LINUX       VmGuestOs = "OTHER_24X_LINUX"
	OS_OTHER_24X_LINUX_64    VmGuestOs = "OTHER_24X_LINUX_64"
	OS_OTHER_26X_LINUX       VmGuestOs = "OTHER_26X_LINUX"
	OS_OTHER_26X_LINUX_64    VmGuestOs = "OTHER_26X_LINUX_64"
	OS_OTHER_3X_LINUX        VmGuestOs = "OTHER_3X_LINUX"
	OS_OTHER_3X_LINUX_64     VmGuestOs = "OTHER_3X_LINUX_64"
	OS_OTHER_4X_LINUX        VmGuestOs = "OTHER_4X_LINUX"
	OS_OTHER_4X_LINUX_64     VmGuestOs = "OTHER_4X_LINUX_64"
	OS_OTHER_5X_LINUX        VmGuestOs = "OTHER_5X_LINUX"
	OS_OTHER_5X_LINUX_64     VmGuestOs = "OTHER_5X_LINUX_64"
	OS_OTHER_LINUX           VmGuestOs = "OTHER_LINUX"
	OS_GENERIC_LINUX         VmGuestOs = "GENERIC_LINUX"
	OS_OTHER_LINUX_64        VmGuestOs = "OTHER_LINUX_64"
	OS_SOLARIS_6             VmGuestOs = "SOLARIS_6"
	OS_SOLARIS_7             VmGuestOs = "SOLARIS_7"
	OS_SOLARIS_8             VmGuestOs = "SOLARIS_8"
	OS_SOLARIS_9             VmGuestOs = "SOLARIS_9"
	OS_SOLARIS_10            VmGuestOs = "SOLARIS_10"
	OS_SOLARIS_10_64         VmGuestOs = "SOLARIS_10_64"
	OS_SOLARIS_11_64         VmGuestOs = "SOLARIS_11_64"
	OS_OS2                   VmGuestOs = "OS2"
	OS_ECOMSTATION           VmGuestOs = "ECOMSTATION"
	OS_ECOMSTATION_2         VmGuestOs = "ECOMSTATION_2"
	OS_NETWARE_4             VmGuestOs = "NETWARE_4"
	OS_NETWARE_5             VmGuestOs = "NETWARE_5"
	OS_NETWARE_6             VmGuestOs = "NETWARE_6"
	OS_OPENSERVER_5          VmGuestOs = "OPENSERVER_5"
	OS_OPENSERVER_6          VmGuestOs = "OPENSERVER_6"
	OS_UNIXWARE_7            VmGuestOs = "UNIXWARE_7"
	OS_DARWIN                VmGuestOs = "DARWIN"
	OS_DARWIN_64             VmGuestOs = "DARWIN_64"
	OS_DARWIN_10             VmGuestOs = "DARWIN_10"
	OS_DARWIN_10_64          VmGuestOs = "DARWIN_10_64"
	OS_DARWIN_11             VmGuestOs = "DARWIN_11"
	OS_DARWIN_11_64          VmGuestOs = "DARWIN_11_64"
	OS_DARWIN_12_64          VmGuestOs = "DARWIN_12_64"
	OS_DARWIN_13_64          VmGuestOs = "DARWIN_13_64"
	OS_DARWIN_14_64          VmGuestOs = "DARWIN_14_64"
	OS_DARWIN_15_64          VmGuestOs = "DARWIN_15_64"
	OS_DARWIN_16_64          VmGuestOs = "DARWIN_16_64"
	OS_DARWIN_17_64          VmGuestOs = "DARWIN_17_64"
	OS_DARWIN_18_64          VmGuestOs = "DARWIN_18_64"
	OS_DARWIN_19_64          VmGuestOs = "DARWIN_19_64"
	OS_DARWIN_20_64          VmGuestOs = "DARWIN_20_64"
	OS_DARWIN_21_64          VmGuestOs = "DARWIN_21_64"
	OS_VMKERNEL              VmGuestOs = "VMKERNEL"
	OS_VMKERNEL_5            VmGuestOs = "VMKERNEL_5"
	OS_VMKERNEL_6            VmGuestOs = "VMKERNEL_6"
	OS_VMKERNEL_65           VmGuestOs = "VMKERNEL_65"
	OS_VMKERNEL_7            VmGuestOs = "VMKERNEL_7"
	OS_AMAZONLINUX2_64       VmGuestOs = "AMAZONLINUX2_64"
	OS_AMAZONLINUX3_64       VmGuestOs = "AMAZONLINUX3_64"
	OS_CRXPOD_1              VmGuestOs = "CRXPOD_1"
	OS_OTHER                 VmGuestOs = "OTHER"
	OS_OTHER_64              VmGuestOs = "OTHER_64"
)

type Template struct {
	Identifier    Identifier              `json:"vm_template"`
	Memory        Memory                  `json:"memory"`
	Disks         map[string]TemplateDisk `json:"disks"`
	Nics          map[string]TemplateNic  `json:"nics"`
	Cpu           TemplateCpu             `json:"cpu"`
	VmHomeStorage TemplateHomeStorage     `json:"vm_home_storage"`
	GuestOS       VmGuestOs               `json:"guest_OS"`
}

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
	GuestOS   VmGuestOs               `json:"guest_OS"`
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

type CpuUpdate struct {
	NumCoresPerSocket int `json:"num_cores_per_socket"`
	NumCpus           int `json:"num_cpus"`
}

type MemoryUpdate struct {
	Memory int `json:"memory"`
}

type HCNicUpdateSpec struct {
	Identifier Identifier `json:"network"`
}

type DiskUpdateSpec struct {
	Capacity int64 `json:"capacity"`
}

type HardwareCustomization struct {
	CpuUpdate     CpuUpdate                  `json:"cpu_update"`
	DisksToRemove []string                   `json:"disks_to_remove"`
	DisksToUpdate map[string]DiskUpdateSpec  `json:"disks_to_update"`
	MemoryUpdate  MemoryUpdate               `json:"memory_update"`
	Nics          map[string]HCNicUpdateSpec `json:"nics"`
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

type DeployGuestCustomization struct {
	Name string `json:"name"`
}

type DeploySpecDiskStoragePolicyType string

const (
	STORAGE_POLICY_USE_SPECIFIED_POLICY DeploySpecDiskStoragePolicyType = "USE_SPECIFIED_POLICY"
	STORAGE_POLICY_USE_SOURCE_POLICY    DeploySpecDiskStoragePolicyType = "USE_SOURCE_POLICY"
)

type DeploySpecDiskStoragePolicy struct {
	Policy null.String                     `json:"policy"`
	Type   DeploySpecDiskStoragePolicyType `json:"type"`
}

type DeploySpecDiskStorage struct {
	Datastore     Identifier                  `json:"datastore"`
	StoragePolicy DeploySpecDiskStoragePolicy `json:"storage_policy"`
}

type DeployTemplateSpec struct {
	Description           string                           `json:"description"`
	DiskStorage           TemplateDiskStorage              `json:"disk_storage"`
	DiskStorageOverrides  map[string]DeploySpecDiskStorage `json:"disk_storage_overrides"`
	GuestCustomization    DeployGuestCustomization         `json:"guest_customization"`
	HardwareCustomization HardwareCustomization            `json:"hardware_customization"`
	Name                  string                           `json:"name"`
	Placement             DeployPlacement                  `json:"placement"`
	PoweredOn             bool                             `json:"powered_on"`
	VmHomeStorage         DeployHomeStorage                `json:"vm_home_storage"`
}

type DeployTemplatePayload struct {
	Spec DeployTemplateSpec `json:"spec"`
}

type VMPowerStateResponse struct {
	State PowerState `json:"state"`
}

type GuestHostnameGeneratorType string

const (
	HOSTNAME_FIXED               GuestHostnameGeneratorType = "FIXED"
	HOSTNAME_PREFIX              GuestHostnameGeneratorType = "PREFIX"
	HOSTNAME_VIRTUAL_MACHINE     GuestHostnameGeneratorType = "VIRTUAL_MACHINE"
	HOSTNAME_USER_INPUT_REQUIRED GuestHostnameGeneratorType = "USER_INPUT_REQUIRED"
)

type GuestHostnameGenerator struct {
	FixedName *string                    `json:"fixed_name"`
	Prefix    *string                    `json:"prefix"`
	Type      GuestHostnameGeneratorType `json:"type"`
}

type GuestLinuxConfiguration struct {
	Domain     string                 `json:"domain"`
	Hostname   GuestHostnameGenerator `json:"hostname"`
	ScriptText *string                `json:"script_text"`
	TimeZone   *string                `json:"time_zone"`
}

type GuestWindowsConfigurationRebootOption string

const (
	WINDOWS_REBOOT    GuestWindowsConfigurationRebootOption = "REBOOT"
	WINDOWS_NO_REBOOT GuestWindowsConfigurationRebootOption = "NO_REBOOT"
	WINDOWS_SHUTDOWN  GuestWindowsConfigurationRebootOption = "SHUTDOWN"
)

type GuestDomainType string

const (
	DOMAIN_WORKGROUP GuestDomainType = "WORKGROUP"
	DOMAIN_DOMAIN    GuestDomainType = "DOMAIN"
)

type GuestDomain struct {
	Domain         *string         `json:"domain"`
	DomainPassword *string         `json:"domain_password"`
	DomainUsername *string         `json:"domain_username"`
	Type           GuestDomainType `json:"type"`
	Workgroup      *string         `json:"workgroup"`
}

type GuestGuiUnattended struct {
	AutoLogon      bool    `json:"auto_logon"`
	AutoLogonCount int64   `json:"auto_logon_count"`
	Password       *string `json:"password"`
	TimeZone       int64   `json:"time_zone"`
}

type GuestUserData struct {
	ComputerName GuestHostnameGenerator `json:"computer_name"`
	FullName     string                 `json:"full_name"`
	Organization string                 `json:"organization"`
	ProductKey   string                 `json:"product_key"`
}

type GuestWindowsSysprep struct {
	Domain             *GuestDomain       `json:"domain"`
	GuiRunOnceCommands *[]string          `json:"gui_run_once_commands"`
	GuiUnattended      GuestGuiUnattended `json:"gui_unattended"`
	UserData           GuestUserData      `json:"user_data"`
}

type GuestWindowsConfiguration struct {
	Reboot     *GuestWindowsConfigurationRebootOption `json:"reboot"`
	Sysprep    *GuestWindowsSysprep                   `json:"sysprep"`
	SysprepXml *string                                `json:"sysprep_xml"`
}

type GuestConfigurationSpec struct {
	LinuxConfig   *GuestLinuxConfiguration   `json:"linux_config"`
	WindowsConfig *GuestWindowsConfiguration `json:"windows_config"`
}

type GuestGlobalDNSSettings struct {
	DnsServers    *[]string `json:"dns_servers"`
	DnsSuffixList *[]string `json:"dns_suffix_list"`
}

type GuestIpv4Type string

const (
	IPV4_DHCP                GuestIpv4Type = "DHCP"
	IPV4_STATIC              GuestIpv4Type = "STATIC"
	IPV4_USER_INPUT_REQUIRED GuestIpv4Type = "USER_INPUT_REQUIRED"
)

type GuestIpv4 struct {
	Gateways  *[]string     `json:"gateways"`
	IpAddress *string       `json:"ip_address"`
	Prefix    *int64        `json:"prefix"`
	Type      GuestIpv4Type `json:"type"`
}

type GuestIpv6Address struct {
	IpAddress string `json:"ip_address"`
	Prefix    int64  `json:"prefix"`
}

type GuestIpv6Type string

const (
	IPV6_DHCP                GuestIpv6Type = "DHCP"
	IPV6_STATIC              GuestIpv6Type = "STATIC"
	IPV6_USER_INPUT_REQUIRED GuestIpv6Type = "USER_INPUT_REQUIRED"
)

type GuestIpv6 struct {
	Gateways *[]string           `json:"gateways"`
	IPv6     *[]GuestIpv6Address `json:"ipv6"`
	Type     GuestIpv6Type       `json:"type"`
}

type GuestWindowsNetworkAdapterSettingsNetBIOSMode string

const (
	NET_BIOS_MODE_USE_DHCP GuestWindowsNetworkAdapterSettingsNetBIOSMode = "USE_DHCP"
	NET_BIOS_MODE_ENABLE   GuestWindowsNetworkAdapterSettingsNetBIOSMode = "ENABLE"
	NET_BIOS_MODE_DISABLE  GuestWindowsNetworkAdapterSettingsNetBIOSMode = "DISABLE"
)

type GuestWindowsNetworkAdapterSettings struct {
	DnsDomain   *string                                        `json:"dns_domain"`
	DnsServers  *[]string                                      `json:"dns_servers"`
	NetBIOSMode *GuestWindowsNetworkAdapterSettingsNetBIOSMode `json:"net_BIOS_mode"`
	WinsServers *[]string                                      `json:"wins_servers"`
}

type GuestIPSettings struct {
	IPv4    *GuestIpv4                          `json:"ipv4"`
	IPv6    *GuestIpv6                          `json:"ipv6"`
	Windows *GuestWindowsNetworkAdapterSettings `json:"windows"`
}

type GuestAdapterMapping struct {
	Adapter    GuestIPSettings `json:"adapter"`
	MacAddress *string         `json:"mac_address"`
}

type GuestCustomizationSpec struct {
	ConfigurationSpec GuestConfigurationSpec `json:"configuration_spec"`
	GlobalDNSSettings GuestGlobalDNSSettings `json:"global_DNS_settings"`
	Interfaces        []GuestAdapterMapping  `json:"interfaces"`
}

type CreateGuestCustomizationSpec struct {
	Description string                 `json:"description"`
	Name        string                 `json:"name"`
	Spec        GuestCustomizationSpec `json:"spec"`
}

const (
	VSPHERE_SESSION_CACHE_FILE = ".vcenter_session.cache"
)

func (vs *VSphere) authorize() (sessionToken string, err error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return
	}
	cachePath := path.Join(homePath, VSPHERE_SESSION_CACHE_FILE)
	if _, fileErr := os.Stat(cachePath); fileErr == nil {
		// Session cache exists, so check if it's valid
		sessionTokenBytes, err := ioutil.ReadFile(cachePath)
		if err != nil {
			return "", err
		}
		sessionToken = string(sessionTokenBytes)
		checkSessionRequest, err := http.NewRequest(http.MethodGet, (vs.BaseUrl + "/api/session"), nil)
		if err != nil {
			return "", err
		}
		checkSessionRequest.Header.Add("vmware-api-session-id", sessionToken)
		checkSessionResponse, err := vs.HttpClient.Do(checkSessionRequest)
		if err != nil {
			return "", err
		}
		// Return the session token if it's still valid
		if checkSessionResponse.StatusCode == http.StatusOK {
			return sessionToken, nil
		}
	}
	// Either we never had a session or it's invalid
	authRequest, err := http.NewRequest(http.MethodPost, (vs.BaseUrl + "/api/session"), nil)
	if err != nil {
		return
	}
	authRequest.SetBasicAuth(vs.Username, vs.Password)
	authResponse, err := vs.HttpClient.Do(authRequest)
	if err != nil {
		return
	}
	if authResponse.StatusCode != 201 {
		err = errors.New("recieved status " + authResponse.Status)
		return
	}

	defer authResponse.Body.Close()

	err = json.NewDecoder(authResponse.Body).Decode(&sessionToken)
	if err != nil {
		return
	}
	// Write that session token to a cache file
	err = ioutil.WriteFile(cachePath, []byte(sessionToken), 0644)
	if err != nil {
		return
	}
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
	request.Header.Set("User-Agent", "LaForge/3.0.1")
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
	request.Header.Set("User-Agent", "LaForge/3.0.1")
	request.Header.Add("vmware-api-session-id", sessionToken)
	request.Header.Add("Content-Type", "application/json")
	return
}

func (vs *VSphere) ListVms() (vms []VirtualMachine, err error) {
	log.Debug("vSphere | ListVms")
	request, err := vs.generateAuthorizedRequest("GET", "/api/vcenter/vm")
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-Okay Response from ListVms")
		err = errors.New("recieved status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&vms)
	if err != nil {
		return
	}
	return
}

func (vs *VSphere) GetVm(name string) (vm VirtualMachine, err error) {
	log.WithFields(log.Fields{
		"name": name,
	}).Debug("vSphere | GetVm")
	request, err := vs.generateAuthorizedRequest("GET", "/api/vcenter/vm?names="+name)
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-Okay Response from GetVm")
		err = errors.New("recieved status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	var vmList []VirtualMachine
	err = json.NewDecoder(response.Body).Decode(&vmList)
	if err != nil {
		return
	}
	if len(vmList) <= 0 {
		err = fmt.Errorf("no vm's were found for the name \"%s\"", name)
		return
	}
	vm = vmList[0]
	return
}

func (vs *VSphere) ListDatastores() (datastores []Datastore, err error) {
	log.Debug("vSphere | ListDatastores")
	request, err := vs.generateAuthorizedRequest("GET", "/api/vcenter/datastore")
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-Okay Response from ListDatastores")
		err = errors.New("received status " + response.Status + " from VSphere while listing datastores")
		return
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&datastores)
	if err != nil {
		return
	}
	return
}

func (vs *VSphere) GetDatastoreByName(name string) (datastore Datastore, err error) {
	log.WithFields(log.Fields{
		"name": name,
	}).Debug("vSphere | GetDatastoreByName")
	request, err := vs.generateAuthorizedRequest("GET", "/api/vcenter/datastore?names="+name)
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-Okay Response from GetDatastoreByName")
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	var datastoreList []Datastore
	err = json.NewDecoder(response.Body).Decode(&datastoreList)
	if err != nil {
		return
	}
	if len(datastoreList) < 1 {
		err = errors.New("no datastore found for the name\"" + name + "\"")
		return
	}
	if len(datastoreList) > 1 {
		err = errors.New("more than one (" + fmt.Sprint(len(datastoreList)) + ") datastore found for the name\"" + name + "\"")
		return
	}
	datastore = datastoreList[0]
	return
}

func (vs *VSphere) ListFolders() (folders []Folder, err error) {
	log.Debug("vSphere | ListFolders")
	request, err := vs.generateAuthorizedRequest("GET", "/api/vcenter/folder")
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-Okay Response from ListFolders")
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&folders)
	if err != nil {
		return
	}
	return
}

func (vs *VSphere) GetFolderByName(name string) (folder Folder, err error) {
	log.WithFields(log.Fields{
		"name": name,
	}).Debug("vSphere | GetFolderByName")
	request, err := vs.generateAuthorizedRequest("GET", "/api/vcenter/folder?names="+name)
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-Okay Response from GetFolderByName")
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	var folderList []Folder
	err = json.NewDecoder(response.Body).Decode(&folderList)
	if err != nil {
		return
	}
	if len(folderList) < 1 {
		err = errors.New("no folder found for the name\"" + name + "\"")
		return
	}
	if len(folderList) > 1 {
		err = errors.New("more than one (" + fmt.Sprint(len(folderList)) + ") folder found for the name\"" + name + "\"")
		return
	}
	folder = folderList[0]
	return
}

func (vs *VSphere) ListResourcePools() (resourcePools []ResourcePool, err error) {
	log.Debug("vSphere | ListResourcePools")
	request, err := vs.generateAuthorizedRequest("GET", "/api/vcenter/resource-pool")
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-Okay Response from ListResourcePools")
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&resourcePools)
	if err != nil {
		return
	}
	return
}

func (vs *VSphere) GetResourcePoolByName(name string) (resourcePool ResourcePool, err error) {
	log.WithFields(log.Fields{
		"name": name,
	}).Debug("vSphere | GetResourcePoolByName")
	request, err := vs.generateAuthorizedRequest("GET", "/api/vcenter/resource-pool?names="+name)
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-Okay Response from GetResourcePoolByName")
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	var resourcePoolList []ResourcePool
	err = json.NewDecoder(response.Body).Decode(&resourcePoolList)
	if err != nil {
		return
	}
	if len(resourcePoolList) < 1 {
		err = errors.New("no resource pool found for the name\"" + name + "\"")
		return
	}
	if len(resourcePoolList) > 1 {
		err = errors.New("more than one (" + fmt.Sprint(len(resourcePoolList)) + ") resource pool found for the name\"" + name + "\"")
		return
	}
	resourcePool = resourcePoolList[0]
	return
}

func (vs *VSphere) ListNetworks() (networks []Network, err error) {
	log.Debug("vSphere | ListNetworks")
	request, err := vs.generateAuthorizedRequest("GET", "/api/vcenter/network")
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-Okay Response from ListNetworks")
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&networks)
	if err != nil {
		return
	}
	return
}

func (vs *VSphere) GetNetworkByName(name string) (network Network, err error) {
	log.WithFields(log.Fields{
		"name": name,
	}).Debug("vSphere | GetNetworkByName")
	request, err := vs.generateAuthorizedRequest("GET", "/api/vcenter/network?names="+name)
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-Okay Response from GetNetworkByName")
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	var networkList []Network
	err = json.NewDecoder(response.Body).Decode(&networkList)
	if err != nil {
		return
	}
	if len(networkList) < 1 {
		err = errors.New("no network found for the name\"" + name + "\"")
		return
	}
	if len(networkList) > 1 {
		err = errors.New("more than one (" + fmt.Sprint(len(networkList)) + ") network found for the name\"" + name + "\"")
		return
	}
	network = networkList[0]
	return
}
func (vs *VSphere) GetTemplateIDByName(contentLibraryName string, templateName string) (templateId string, err error) {
	log.WithFields(log.Fields{
		"contentLibraryName": contentLibraryName,
		"templateName":       templateName,
	}).Debug("vSphere | GetTemplateIDByName")
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
	log.WithFields(log.Fields{
		"templateId": templateId,
	}).Debug("vSphere | GetTemplate")
	request, err := vs.generateAuthorizedRequest("GET", "/api/vcenter/vm-template/library-items/"+templateId)
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-Okay Response from GetTemplate")
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&template)
	if err != nil {
		return
	}
	return
}

func (vs *VSphere) CreateVM(vmSpec VirtualMachineSpec) (err error) {
	log.WithFields(log.Fields{
		"vmSpec.name": vmSpec.Name,
	}).Debug("vSphere | CreateVM")
	requestData := CreateVirtualMachineData{
		Spec: vmSpec,
	}
	requestDataString, err := json.Marshal(requestData)
	if err != nil {
		return
	}
	request, err := vs.generateAuthorizedRequestWithData("POST", "/api/vcenter/vm", bytes.NewBuffer(requestDataString))
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-Okay Response from CreateVM")
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	// TODO: Actually parse the response
	return
}

func (vs *VSphere) GetVMPowerState(vmIdentifier Identifier) (powerState PowerState, err error) {
	log.WithFields(log.Fields{
		"vmIdentifier": vmIdentifier,
	}).Debug("vSphere | GetVMPowerState")
	request, err := vs.generateAuthorizedRequest(http.MethodGet, ("/api/vcenter/vm/" + string(vmIdentifier) + "/power"))
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-Okay Response from ListResourcePools")
		err = errors.New("recieved status " + response.Status + " from VSphere")
		return
	}
	var responseBody VMPowerStateResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return
	}
	return responseBody.State, nil
}

func (vs *VSphere) ShutdownVM(vmIdentifier Identifier) (err error) {
	log.WithFields(log.Fields{
		"vmIdentifier": vmIdentifier,
	}).Debug("vSphere | ShutdownVM")
	request, err := vs.generateAuthorizedRequest(http.MethodPost, ("/api/vcenter/vm/" + string(vmIdentifier) + "/power?action=stop"))
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusNoContent {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-NoContent Response from ShutdownVM")
		return fmt.Errorf("error while trying to shutdown VM %s: %s", vmIdentifier, response.Status)
	}
	return
}

func (vs *VSphere) DeleteVM(name string) (err error) {
	log.WithFields(log.Fields{
		"name": name,
	}).Debug("vSphere | DeleteVM")
	vm, err := vs.GetVm(name)
	if err != nil {
		return
	}

	// Check if the VM is turned on
	powerState, err := vs.GetVMPowerState(vm.Identifier)
	if err != nil {
		return
	}
	// Shutdown the VM prior to deleting it
	if powerState == POWER_STATE_ON || powerState == POWER_STATE_SUSPENDED {
		log.WithFields(log.Fields{
			"vmIdentifier": vm.Identifier,
		}).Debug("Shutting down VM prior to deletion...")
		err = vs.ShutdownVM(vm.Identifier)
		if err != nil {
			return err
		}
	}

	request, err := vs.generateAuthorizedRequest(http.MethodDelete, ("/api/vcenter/vm/" + string(vm.Identifier)))
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusNoContent {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-NoContent Response from DeleteVM")
		return errors.New("recieved status " + response.Status + " from VSphere while trying to delete the VM: " + name)
	}
	return
}

func (vs *VSphere) DeployTemplate(templateId string, spec DeployTemplateSpec) (err error) {
	log.WithFields(log.Fields{
		"templateId": templateId,
		"spec.Name":  spec.Name,
	}).Debug("vSphere | DeployTemplate")
	requestDataString, err := json.Marshal(spec)
	if err != nil {
		return
	}
	request, err := vs.generateAuthorizedRequestWithData("POST", "/api/vcenter/vm-template/library-items/"+templateId+"?action=deploy", bytes.NewBuffer(requestDataString))
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	// DEBUG: output the deployment id for debug purposes
	defer response.Body.Close()
	deploymentBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	if response.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"status": response.Status,
		}).Warn("vSphere | Non-Okay Response from DeployTemplate")
		log.Warnf("%s\n", deploymentBytes)
		err = errors.New("received status " + response.Status + " from VSphere")
		return
	}

	log.WithFields(log.Fields{
		"name":       spec.Name,
		"identifier": string(deploymentBytes),
	}).Debug("Deployed VM from Template")

	return
}

func (vs *VSphere) GenerateGuestCustomization(ctx context.Context, template Template, provisionedHost *ent.ProvisionedHost) (spec *GuestCustomizationSpec, err error) {
	log.WithFields(log.Fields{
		"template.Identifier": template.Identifier,
		"provisionedHost.ID":  provisionedHost.ID,
	}).Debug("vSphere | GenerateGuestCustomization")
	host, err := provisionedHost.QueryProvisionedHostToHost().Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while querying host from provisioned host: %v", err)
	}
	env, err := provisionedHost.QueryProvisionedHostToProvisionedNetwork().QueryProvisionedNetworkToTeam().QueryTeamToBuild().QueryBuildToEnvironment().Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while querying environment from provisioned host: %v", err)
	}
	// team, err := provisionedHost.QueryProvisionedHostToProvisionedNetwork().QueryProvisionedNetworkToTeam().Only(ctx)
	// if err != nil {
	// 	return nil, fmt.Errorf("error while querying team from provisioned host: %v", err)
	// }
	provisionedNetwork, err := provisionedHost.QueryProvisionedHostToProvisionedNetwork().Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("err while querying provisioned network from provisioned host: %v", err)
	}
	network, err := provisionedNetwork.QueryProvisionedNetworkToNetwork().Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("error while querying network from provisioned network: %v", err)
	}
	var linuxConfig *GuestLinuxConfiguration = nil
	var windowsConfig *GuestWindowsConfiguration = nil
	var windowsAdapterSettings *GuestWindowsNetworkAdapterSettings = nil
	var dnsServers []string = []string{}

	// Check if the network has an authoritative DNS server
	authoritativeDns, exists := network.Vars["authoritative_dns_ip"]
	if exists {
		dnsServers = append(dnsServers, authoritativeDns)
	}
	// Check for master DNS config entry
	masterDnsServer, exists := env.Config["master_dns_server"]
	if exists {
		dnsServers = append(dnsServers, masterDnsServer)
	}
	// Back up with Google DNS in case no dns server is set
	dnsServers = append(dnsServers, "8.8.8.8")

	guestHostname := GuestHostnameGenerator{
		FixedName: &host.Hostname,
		Type:      HOSTNAME_FIXED,
	}

	agentUrl := fmt.Sprintf("%s/agents/%s", vs.ServerUrl, provisionedHost.ID)
	// If Windows
	if strings.Contains(string(template.GuestOS), "WIN") {
		windowsConfig = &GuestWindowsConfiguration{
			Sysprep: &GuestWindowsSysprep{
				GuiRunOnceCommands: &[]string{
					fmt.Sprintf("Invoke-WebRequest \"%s\" -OutFile \"C:\\laforge.exe\"", agentUrl),
					"New-Service -name \"laforge\" -binaryPathName \"C:\\laforge.exe\" -displayName \"CPTC LaForge Agent\" -startupType Automatic",
					"Start-Service -name \"laforge\"",
				},
				UserData: GuestUserData{
					ComputerName: guestHostname,
					FullName:     "",
					Organization: "",
					ProductKey:   "",
				},
			},
			SysprepXml: new(string),
		}
		windowsAdapterSettings = &GuestWindowsNetworkAdapterSettings{
			DnsServers: &dnsServers,
		}
	} else {
		// Otherwise must be Unix
		configScript := fmt.Sprintf(`# create the systemd service file for the Laforge
		curl -sL -o /laforge.bin %s
		chmod +x /laforge.bin
		
		cat << EOF > /etc/systemd/system/laforge.service
		# /etc/systemd/system/laforge.service
		# version 0.1
		[Unit]
		Description=Laforge Agent Service
		After=network.target
		[Service]
		Environment=HOME=/
		Environment=USER=root
		User=root
		Group=root
		ExecStart=/laforge.bin
		[Install]
		WantedBy=multi-user.target
		EOF
		
		systemctl daemon-reload
		systemctl enable laforge
		systemctl start laforge`, agentUrl)
		linuxConfig = &GuestLinuxConfiguration{
			Domain:     "",
			Hostname:   guestHostname,
			ScriptText: &configScript,
		}
	}

	networkAddrParts := strings.Split(provisionedNetwork.Cidr, "/")
	networkAddr := networkAddrParts[0]
	networkOctets := strings.Split(networkAddr, ".")
	hostAddress := strings.Join(append(networkOctets[:3], fmt.Sprint(host.LastOctet)), ".")
	gatewayAddress := strings.Join(append(networkOctets[:3], "254"), ".")
	cidrPrefix, err := strconv.ParseInt(networkAddrParts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error while parsing CIDR Prefix: %v", err)
	}

	spec = &GuestCustomizationSpec{
		ConfigurationSpec: GuestConfigurationSpec{
			LinuxConfig:   linuxConfig,
			WindowsConfig: windowsConfig,
		},
		GlobalDNSSettings: GuestGlobalDNSSettings{
			DnsServers: &dnsServers,
		},
		Interfaces: []GuestAdapterMapping{
			{
				Adapter: GuestIPSettings{
					IPv4: &GuestIpv4{
						Gateways: &[]string{
							gatewayAddress,
						},
						IpAddress: &hostAddress,
						Prefix:    &cidrPrefix,
						Type:      IPV4_STATIC,
					},
					Windows: windowsAdapterSettings,
				},
				MacAddress: new(string),
			},
		},
	}
	return spec, nil
}

func (vs *VSphere) CreateGuestCustomization(name string, spec GuestCustomizationSpec) (err error) {
	log.WithFields(log.Fields{
		"name": name,
	}).Debug("vSphere | CreateGuestCustomization")
	payload := CreateGuestCustomizationSpec{
		Description: "Created by LaForge",
		Name:        name,
		Spec:        spec,
	}
	jsonString, err := json.Marshal(payload)
	if err != nil {
		return
	}
	request, err := vs.generateAuthorizedRequestWithData(http.MethodPost, "/api/vcenter/guest/customization-specs", bytes.NewBuffer(jsonString))
	if err != nil {
		return
	}
	response, err := vs.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != 201 {
		return fmt.Errorf("unkown error %s while creating guest customization spec %s", response.Status, name)
	}
	return
}
