package builder

import (
	"net/http"

	"github.com/gen0cide/laforge/builder/vsphere"
)

type CyberRangeBuilder struct {
	HttpClient                http.Client
	Username                  string
	Password                  string
	NsxtUrl                   string
	VSphereClient             vsphere.VSphere
	VSphereContentLibraryName string
	VSphereDatastoreName      string
	VShpereResourcePoolName   string
	VSphereFolderName         string
	VSphereDatastore          vsphere.Datastore
	VShpereResourcePool       vsphere.ResourcePool
	VSphereFolder             vsphere.Folder
}

func (builder CyberRangeBuilder) Init(url string) (err error) {
	builder.VSphereClient = vsphere.VSphere{
		Client:   builder.HttpClient,
		BaseUrl:  url,
		Username: builder.Username,
		Password: builder.Password,
	}
	return nil
}

func (builder CyberRangeBuilder) DeploySpec(spec Spec) (err error) {
	// builder.VSphereClient.DeployTemplate(spec.Template)
	return nil
}
