# LaForge vSphere + NSX-T Builder

> *NOTE: This builder is still in development and is only verified to work with VMware vCenter Server 7.0*

## Setting up vSphere for use with LaForge

There is a bit of setup required in order to get LaForge to run with vSphere and NSX-T.

### VM Templates and Content Libraries

The LaForge vSphere + NSX-T builder takes advantage of vSphere's [Content Libraries](https://docs.vmware.com/en/VMware-vSphere/7.0/com.vmware.vsphere.vm_admin.doc/GUID-254B2CE8-20A8-43F0-90E8-3F6776C2C896.html) to deploy it's hosts.

#### Windows Templates

Follow the guide [here](https://docs.microsoft.com/en-us/azure/cloud-adoption-framework/manage/hybrid/server/best-practices/vmware-windows-template?/azure/cloud-adoption-framework/_bread/toc.json&toc=/azure/cloud-adoption-framework/scenarios/hybrid/toc.json)

**NOTE: Stop at the "Convert to template" Step, we use a different method to place the template in the content library**

#### Linux Templates

```shell
$ sudo [apt-get/yum] update
$ sudo [apt-get/yum] upgrade -y
$ sudo [apt-get/yum] install -y open-vm-tools cloud-init
$ sudo dpkg-reconfigure cloud-init # NoCloud, OVF, None
$ sudo nano /etc/cloud/cloud.cfg # add `disable_vmware_customization: false` to EOF
$ sudo sed -i 's/preserve_hostname: false/preserve_hostname: true/g' /etc/cloud/cloud.cfg
$ sudo truncate -s0 /etc/hostname
$ sudo hostnamectl set-hostname localhost
$ sudo vmware-toolbox-cmd config set deployPkg enable-custom-scripts true
$ cat /dev/null > ~/.bash_history && history -c
$ sudo shutdown now
```

## LaForge environment config variables

### Specifying this builder in an environment

In order to use this builder for an environment, you have to configure the `builder` field in the environment laforge config.

```terraform
environment "/envs/xxxxx" {
  // ...
  builder = "vsphere-nsxt"
  // ...
}
```

### New LaForge Specific Config Vars

In order to deploy the agents to all the hosts, the builder needs to know the Fully Qualified Domain Name (FQDN) of the LaForge server.

```terraform
environment "/envs/xxxxx" {
  // ...
  config = {
    // ...
    laforge_server_url      = "<LAFORGE_URL>"          // The FQDN of the main LaForge server
    // ...
  }
  // ...
}
```

### vSphere Specific Config Vars

The builder requires some vSphere/vCenter specific config variables to be able to interact with the vCenter Automation API.

```terraform
environment "/envs/xxxxx" {
  // ...
  config = {
    // ...
    vsphere_username        = "<USERNAME>"             // Username of the vSphere user used for LaForge
    vsphere_password        = "<PASSWORD>"             // Password of the LaForge user
    vsphere_base_url        = "<URL>"                  // The URL of the vSphere server
    vsphere_content_library = "<CONTENT_LIBRARY_NAME>" // The name of the Content Library containing the VM Templates for LaForge
    vsphere_datastore       = "<DATASTORE_NAME>"       // The name of the Datastore to place the VMs on
    vsphere_resource_pool   = "<RESOURCE_POOL_NAME>"   // The name of the Resource Pool to assign the VMs to
    vsphere_folder          = "<FOLDER_NAME>"          // The name of the Folder to put the VMs in
    vsphere_template_prefix = "<PREFIX_NAME>"          // The prefix given to each of the LaForge VM templates
    // ...
  }
  // ...
}
```

### NSX-T Specific Config Vars

The builder requires some NSX-T specific config variables to be able to interact with the NSX-T Datacenter API.

```terraform
environment "/envs/xxxxx" {
  // ...
  config = {
    // ...
    nsxt_cert_path          = "<CERT_PATH>"            // The absolute path to the certificate of the Principal Identity User for NSX-T
    nsxt_ca_cert_path       = "<CA_CERT_PATH>"         // The absolute path to the CA Cert for the Pricipal Identity User
    nsxt_key_path           = "<PRIVATE_KEY_PATH>"     // The absolute path to the Private Key for the Principal Identity User
    nsxt_base_url           = "<URL>"                  // The URL of the NSX-T server
    nsxt_ip_pool_name       = "<IP_POOL_NAME>"         // The name of the IP Pool to be used for NAT
    // ...
  }
  // ...
}
```