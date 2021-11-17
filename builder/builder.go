package builder

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gen0cide/laforge/builder/vspherensxt"
	"github.com/gen0cide/laforge/builder/vspherensxt/nsxt"
	"github.com/gen0cide/laforge/builder/vspherensxt/vsphere"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/logging"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

type Builder interface {
	ID() string
	Name() string
	Description() string
	Author() string
	Version() string
	DeployHost(ctx context.Context, provisionedHost *ent.ProvisionedHost) (err error)
	DeployNetwork(ctx context.Context, provisionedNetwork *ent.ProvisionedNetwork) (err error)
	TeardownHost(ctx context.Context, provisionedHost *ent.ProvisionedHost) (err error)
	TeardownNetwork(ctx context.Context, provisionedNetwork *ent.ProvisionedNetwork) (err error)
}

func BuilderFromEnvironment(environment *ent.Environment, logger *logging.Logger) (genericBuilder Builder, err error) {
	switch environment.Builder {
	case "vsphere-nsxt":
		genericBuilder, err = NewVSphereNSXTBuilder(environment, logger)
		if err != nil {
			logrus.Errorf("Failed to make vSphere NSX-T builder. Err: %v", err)
			return
		}
		return
	}
	err = fmt.Errorf("error: builder not found")
	logrus.Error(err)
	return
}

// NewVSphereNSXTBuilder creates a builder instance to deploy environments to VSphere and NSX-T
func NewVSphereNSXTBuilder(env *ent.Environment, logger *logging.Logger) (builder vspherensxt.VSphereNSXTBuilder, err error) {
	laforgeServerUrl, exists := env.Config["laforge_server_url"]
	if !exists {
		err = errors.New("laforge_server_url doesn't exist in the environment configuration")
		return
	}
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
	nsxtCertPath, exists := env.Config["nsxt_cert_path"]
	if !exists {
		err = errors.New("nsxt_cert_path doesn't exist in the environment configuration")
		return
	}
	nsxtCACertPath, exists := env.Config["nsxt_ca_cert_path"]
	if !exists {
		err = errors.New("nsxt_ca_cert_path doesn't exist in the environment configuration")
		return
	}
	nsxtKeyPath, exists := env.Config["nsxt_key_path"]
	if !exists {
		err = errors.New("nsxt_key_path doesn't exist in the environment configuration")
		return
	}
	nsxtBaseUrl, exists := env.Config["nsxt_base_url"]
	if !exists {
		err = errors.New("nsxt_base_url doesn't exist in the environment configuration")
		return
	}
	nsxtIpPoolName, exists := env.Config["nsxt_ip_pool_name"]
	if !exists {
		err = errors.New("nsxt_ip_pool_name doesn't exist in the environment configuration")
		return
	}
	nsxtEdgeClusterPath, exists := env.Config["nsxt_edge_cluster_path"]
	if !exists {
		err = errors.New("nsxt_edge_cluster_path doesn't exist in the environment configuration")
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
	maxBuildWorkersString, exists := env.Config["vsphere_max_build_workers"]
	if !exists {
		maxBuildWorkersString = "8"
	}
	maxBuildThreads, err := strconv.Atoi(maxBuildWorkersString)
	if err != nil {
		maxBuildThreads = 8
	}
	maxTeardownWorkersString, exists := env.Config["vsphere_max_teardown_workers"]
	if !exists {
		maxTeardownWorkersString = "16"
	}
	maxTeardownThreads, err := strconv.Atoi(maxTeardownWorkersString)
	if err != nil {
		maxTeardownThreads = 16
	}

	httpClient := http.Client{
		Timeout: 5 * time.Minute,
	}

	nsxtHttpClient, err := nsxt.NewPrincipalIdentityClient(nsxtCertPath, nsxtKeyPath, nsxtCACertPath)
	if err != nil {
		return
	}

	nsxtClient := nsxt.NSXTClient{
		HttpClient:      nsxtHttpClient,
		BaseUrl:         nsxtBaseUrl,
		IpPoolName:      nsxtIpPoolName,
		EdgeClusterPath: nsxtEdgeClusterPath,
		MaxRetries:      10,
		Logger:          logger,
	}

	vsphereClient := vsphere.VSphere{
		HttpClient: httpClient,
		ServerUrl:  laforgeServerUrl,
		BaseUrl:    vsphereBaseUrl,
		Username:   vsphereUsername,
		Password:   vspherePassword,
		MaxRetries: 10,
		Logger:     logger,
	}

	vsphere.InitializeGovmomi(&vsphereClient, vsphereBaseUrl, vsphereUsername, vspherePassword)

	ctx := context.Background()

	datastore, exists, err := vsphereClient.GetDatastoreSummaryByName(ctx, datastoreName)
	if err != nil {
		return
	}
	if !exists {
		err = fmt.Errorf("error datastore \"%s\" doesn't exist", datastoreName)
		logrus.Error(err)
		return
	}

	folder, err := vsphereClient.GetFolderSummaryByName(ctx, folderName)
	if err != nil {
		err = fmt.Errorf("error finding folder: %v", err)
		logrus.Error(err)
		return
	}

	resourcePool, err := vsphereClient.Finder.ResourcePool(ctx, resourcePoolName)
	if err != nil {
		err = fmt.Errorf("error finding resource pool: %v", err)
		logrus.Error(err)
		return
	}

	deployWorkerPool := semaphore.NewWeighted(int64(maxBuildThreads))
	teardownWorkerPool := semaphore.NewWeighted(int64(maxTeardownThreads))

	builder = vspherensxt.VSphereNSXTBuilder{
		HttpClient:                httpClient,
		Username:                  vsphereUsername,
		Password:                  vspherePassword,
		NsxtClient:                nsxtClient,
		TemplatePrefix:            templatePrefix,
		VSphereClient:             vsphereClient,
		VSphereContentLibraryName: contentLibraryName,
		VSphereDatastore:          datastore,
		VSphereResourcePool:       resourcePool,
		VSphereFolder:             folder,
		Logger:                    logger,
		MaxWorkers:                maxBuildThreads,
		DeployWorkerPool:          deployWorkerPool,
		TeardownWorkerPool:        teardownWorkerPool,
	}
	return
}
