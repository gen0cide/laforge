// Package nsxt is for interfacing with the NSX-T REST API
package nsxt

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type NSXTClient struct {
	BaseUrl    string
	HttpClient http.Client
	IpPoolName string
}

type NSXTResourceType string

const (
	NSXT_Infra                     NSXTResourceType = "Infra"
	NSXT_ChildSegment              NSXTResourceType = "ChildSegment"
	NSXT_ChildTier1                NSXTResourceType = "ChildTier1"
	NSXT_ChildLocaleServices       NSXTResourceType = "ChildLocaleServices"
	NSXT_Segment                   NSXTResourceType = "Segment"
	NSXT_Tier0                     NSXTResourceType = "Tier0"
	NSXT_Tier1                     NSXTResourceType = "Tier1"
	NSXT_LocaleServices            NSXTResourceType = "LocaleServices"
	NSXT_IpAddressPoolStaticSubnet NSXTResourceType = "IpAddressPoolStaticSubnet"
)

type NSXTSubnet struct {
	GatewayAddress string `json:"gateway_address"`
}

type NSXTSegment struct {
	ResourceType     NSXTResourceType `json:"resource_type"`
	ConnectivityPath string           `json:"connectivity_path"`
	ID               string           `json:"id"`
	Subnets          []NSXTSubnet     `json:"subnets"`
}

type NSXTChildSegment struct {
	ResourceType NSXTResourceType `json:"resource_type"`
	Segment      NSXTSegment      `json:"Segment"`
}

type NSXTTier1Type string

const (
	NSXTTier1_NATTED   NSXTTier1Type = "NATTED"
	NSXTTier1_ROUTED   NSXTTier1Type = "ROUTED"
	NSXTTier1_ISOLATED NSXTTier1Type = "ISOLATED"
)

type NSXTRouteAdvertisementType string

const (
	NSXT_Tier1_RA_STATIC_ROUTES        NSXTRouteAdvertisementType = "TIER1_STATIC_ROUTES"
	NSXT_Tier1_RA_CONNECTED            NSXTRouteAdvertisementType = "TIER1_CONNECTED"
	NSXT_Tier1_RA_NAT                  NSXTRouteAdvertisementType = "TIER1_NAT"
	NSXT_Tier1_RA_LB_VIP               NSXTRouteAdvertisementType = "TIER1_LB_VIP"
	NSXT_Tier1_RA_LB_SNAT              NSXTRouteAdvertisementType = "TIER1_LB_SNAT"
	NSXT_Tier1_RA_DNS_FORWARDER_IP     NSXTRouteAdvertisementType = "TIER1_DNS_FORWARDER_IP"
	NSXT_Tier1_RA_IPSEC_LOCAL_ENDPOINT NSXTRouteAdvertisementType = "TIER1_IPSEC_LOCAL_ENDPOINT"
)

type NSXTLocaleServices struct {
	ResourceType    NSXTResourceType `json:"resource_type"`
	EdgeClusterPath string           `json:"edge_cluster_path"`
	ID              string           `json:"id"`
}

type NSXTChildLocaleServices struct {
	ResourceType   NSXTResourceType   `json:"resource_type"`
	LocaleServices NSXTLocaleServices `json:"LocaleServices"`
}

type NSXTTier1 struct {
	ResourceType            NSXTResourceType             `json:"resource_type"`
	ID                      string                       `json:"id"`
	RouteAdvertisementTypes []NSXTRouteAdvertisementType `json:"route_advertisement_types"`
	Tier0Path               string                       `json:"tier0_path"`
	Children                []NSXTChildLocaleServices    `json:"children"`
}

type NSXTChildTier1 struct {
	ResourceType NSXTResourceType `json:"resource_type"`
	Tier1        NSXTTier1        `json:"Tier1"`
}

type NSXTCreateTier1Payload struct {
	ResourceType NSXTResourceType `json:"resource_type"`
	Children     []NSXTChildTier1 `json:"children"`
}

type NSXTCreateSegmentPayload struct {
	ResourceType NSXTResourceType   `json:"resource_type"`
	Children     []NSXTChildSegment `json:"children"`
}

type NSXpostDHCPpayload struct { //interface that defines the JSON body for adding a DHCP profile

	Display_name_post string `json:"display_name"`
	Edge_cluster_id   string `json:"edge_cluster_id"`
}

type NSXTAddSubnetPayload struct {
	Subnets          []NSXTSubnet `json:"subnets"`
	ConnectivityPath string       `json:"connectivity_path"`
}

type NSXTTier0AdvancedConfig struct {
	ForwardingUpTimer int    `json:"forwarding_up_timer"`
	Connectivity      string `json:"connectivity"`
}

type NSXTTier0 struct {
	TransitSubnets         []string                `json:"transit_subnets"`
	InternalTransitSubnets []string                `json:"internal_transit_subnets"`
	HaMode                 string                  `json:"ha_mode"`
	FailoverMode           string                  `json:"failover_mode"`
	Ipv6ProfilePaths       []string                `json:"ipv6_profile_paths"`
	ForceWhitelisting      bool                    `json:"force_whitelisting"`
	DefaultRuleLogging     bool                    `json:"default_rule_logging"`
	DisableFirewall        bool                    `json:"disable_firewall"`
	AdvancedConfig         NSXTTier0AdvancedConfig `json:"advanced_config"`
	ResourceType           NSXTResourceType        `json:"resource_type"`
	ID                     string                  `json:"id"`
	DisplayName            string                  `json:"display_name"`
	Path                   string                  `json:"path"`
	RelativePath           string                  `json:"relative_path"`
	ParentPath             string                  `json:"parent_path"`
	UniqueId               string                  `json:"unique_id"`
	MarkedForDelete        bool                    `json:"marked_for_delete"`
	Overridden             bool                    `json:"overridden"`
	CreateUser             string                  `json:"_create_user"`
	CreateTime             uint                    `json:"_create_time"`
	LastModifiedUser       string                  `json:"_last_modified_user"`
	LastModifiedTime       uint                    `json:"_last_modified_time"`
	SystemOwned            bool                    `json:"_system_owned"`
	Protection             string                  `json:"_protection"`
	Revision               int                     `json:"_revision"`
}

type NSXTListTier0Result struct {
	Results       []NSXTTier0 `json:"results"`
	ResultCount   int         `json:"result_count"`
	SortBy        string      `json:"sort_by"`
	SortAscending bool        `json:"sort_ascending"`
}

type NSXTErrorCode int

const (
	NSXT_Tier1_Has_Children NSXTErrorCode = 500030
	NSXT_Segment_Has_VMs    NSXTErrorCode = 503040
)

type NSXTErrorResponse struct {
	HttpStatus string        `json:"httpStatus"`
	ErrorCode  NSXTErrorCode `json:"error_code"`
	ModuleName string        `json:"module_name"`
	Message    string        `json:"error_message"`
}

type NSXTIpAddress string

type NSXTAllocationRange struct {
	Start NSXTIpAddress `json:"start"`
	End   NSXTIpAddress `json:"end"`
}

type NSXTIpSubnet struct {
	Cidr             NSXTIpAddress         `json:"cidr"`
	GatewayIp        NSXTIpAddress         `json:"gateway_ip"`
	DnsNameservers   []NSXTIpAddress       `json:"dns_nameservers"`
	DnsSuffix        string                `json:"dns_suffix"`
	AllocationRanges []NSXTAllocationRange `json:"allocation_ranges"`
	ResourceType     NSXTResourceType      `json:"resource_type"`
	Id               string                `json:"id"`
	DisplayName      string                `json:"display_name"`
	Path             string                `json:"path"`
	RelativePath     string                `json:"relative_path"`
	ParentPath       string                `json:"parent_path"`
	UniqueId         string                `json:"unique_id"`
	MarkedForDelete  bool                  `json:"marked_for_delete"`
	Overridden       bool                  `json:"overridden"`
	CreateUser       string                `json:"_create_user"`
	CreateTime       uint                  `json:"_create_time"`
	LastModifiedUser string                `json:"_last_modified_user"`
	LastModifiedTime uint                  `json:"_last_modified_time"`
	SystemOwned      bool                  `json:"_system_owned"`
	Protection       string                `json:"_protection"`
	Revision         int                   `json:"_revision"`
}

type NSXTListIpSubnetsResponse struct {
	Results       []NSXTIpSubnet `json:"results"`
	ResultCount   int            `json:"result_count"`
	SortBy        string         `json:"sort_by"`
	SortAscending bool           `json:"sort_ascending"`
}

type NSXTNATAction string

const (
	NSXT_NAT_SNAT NSXTNATAction = "SNAT"
)

type NSXTIPElementList string

type NSXTNATRule struct {
	Description       string            `json:"description"`
	Action            NSXTNATAction     `json:"action"`
	Id                string            `json:"id"`
	SourceNetwork     NSXTIPElementList `json:"source_network"`
	TranslatedNetwork NSXTIPElementList `json:"translated_network"`
}

func NewPrincipalIdentityClient(certPath, keyPath, caCertPath string) (client http.Client, err error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return
	}
	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return
	}
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		caCertPool = x509.NewCertPool()
	}
	ok := caCertPool.AppendCertsFromPEM(caCert)
	if !ok {
		err = errors.New("failed to add Root CA to Certificate Pool")
		return
	}

	tlsConfig := tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		// Force the http client to send the cert every time, regardless of TLS compatibility
		GetClientCertificate: func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return &cert, nil
		},
	}
	transport := http.Transport{
		TLSClientConfig: &tlsConfig,
	}
	client = http.Client{
		Transport: &transport,
		Timeout:   2 * time.Minute,
	}
	return
}

func (nsxt *NSXTClient) generateAuthorizedRequest(method, path string) (request *http.Request, err error) {
	request, err = http.NewRequest(method, (nsxt.BaseUrl + path), nil)
	if err != nil {
		return
	}
	request.Header.Set("User-Agent", "LaForge/3.0.1")
	return
}

func (nsxt *NSXTClient) generateAuthorizedRequestWithData(method string, path string, data *bytes.Buffer) (request *http.Request, err error) {
	request, err = http.NewRequest(method, (nsxt.BaseUrl + path), data)
	if err != nil {
		return
	}
	request.Header.Set("User-Agent", "LaForge/3.0.1")
	request.Header.Add("Content-Type", "application/json")
	return
}

func (nsxt *NSXTClient) CreateTier1(name string, tier0Path string, edgeClusterPath string) (err error) {
	payload := NSXTCreateTier1Payload{
		ResourceType: NSXT_Infra,
		Children: []NSXTChildTier1{
			{
				ResourceType: NSXT_ChildTier1,
				Tier1: NSXTTier1{
					ResourceType: NSXT_Tier1,
					ID:           name,
					RouteAdvertisementTypes: []NSXTRouteAdvertisementType{
						NSXT_Tier1_RA_NAT,
					},
					Tier0Path: tier0Path,
					Children: []NSXTChildLocaleServices{
						{
							ResourceType: NSXT_ChildLocaleServices,
							LocaleServices: NSXTLocaleServices{
								ResourceType:    NSXT_LocaleServices,
								EdgeClusterPath: edgeClusterPath,
								ID:              name + "-Edge-Routing",
							},
						},
					},
				},
			},
		},
	}

	jsonString, err := json.Marshal(payload)
	if err != nil {
		return
	}
	request, err := nsxt.generateAuthorizedRequestWithData(http.MethodPatch, "/policy/api/v1/infra", bytes.NewBuffer(jsonString))
	if err != nil {
		return
	}
	response, err := nsxt.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != 200 && response.Body != nil {
		var nsxtError NSXTErrorResponse
		err = json.NewDecoder(response.Body).Decode(&nsxtError)
		if err != nil {
			return
		}
		if nsxtError.ErrorCode == 500087 || nsxtError.Message == "The object was modified by somebody else. Please retry." {
			return
		}
		log.Errorf("error while creating tier-1: %v", nsxtError)
		err = errors.New("recieved status " + response.Status + " from NSX-T while adding Tier 1 Router " + name)
		return
	}
	return
}

func (nsxt *NSXTClient) DeleteTier1(name string) (nsxtError *NSXTErrorResponse, err error) {
	edgeRoutingRequest, err := nsxt.generateAuthorizedRequest(http.MethodDelete, ("/policy/api/v1/infra/tier-1s/" + name + "/locale-services/" + name + "-Edge-Routing"))
	if err != nil {
		return nil, fmt.Errorf("error while making the DELETE request for Edge-Routing: %v", err)
	}
	edgeRoutingResponse, err := nsxt.HttpClient.Do(edgeRoutingRequest)
	if err != nil {
		return nil, fmt.Errorf("unknown error while deleting Edge-Routing for %s: %v", name, err)
	}
	if edgeRoutingResponse.StatusCode != 200 {
		var nsxtError NSXTErrorResponse
		json.NewDecoder(edgeRoutingResponse.Body).Decode(&nsxtError)
		return &nsxtError, nil
	}
	tier1Request, err := nsxt.generateAuthorizedRequest(http.MethodDelete, ("/policy/api/v1/infra/tier-1s/" + name))
	if err != nil {
		return nil, fmt.Errorf("error while making DELETE request for Tier-1: %v", err)
	}
	tier1Response, err := nsxt.HttpClient.Do(tier1Request)
	if err != nil {
		return nil, fmt.Errorf("unkown error while deleting Tier-1 %s: %v", name, err)
	}
	if tier1Response.StatusCode != 200 {
		var nsxtError NSXTErrorResponse
		json.NewDecoder(tier1Response.Body).Decode(&nsxtError)
		return &nsxtError, nil
	}
	return
}

func (nsxt *NSXTClient) CreateSegment(name string, tier1path string, gatewayAddress string) (err error) {
	payload := NSXTCreateSegmentPayload{
		ResourceType: "Infra",
		Children: []NSXTChildSegment{
			{
				ResourceType: "ChildSegment",
				Segment: NSXTSegment{
					ResourceType:     "Segment",
					ID:               name,
					ConnectivityPath: tier1path,
					Subnets: []NSXTSubnet{
						{
							GatewayAddress: gatewayAddress,
						},
					},
				},
			},
		},
	}

	jsonString, err := json.Marshal(payload)
	if err != nil {
		return
	}
	request, err := nsxt.generateAuthorizedRequestWithData(http.MethodPatch, "/policy/api/v1/infra", bytes.NewBuffer(jsonString))
	if err != nil {
		return
	}
	response, err := nsxt.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != 200 {
		err = errors.New("recieved status " + fmt.Sprint(response.StatusCode) + " from NSX-T while adding segment " + name)
		return
	}
	return
}

func (nsxt *NSXTClient) DeleteSegment(name string) (nsxtError *NSXTErrorResponse, err error) {
	segmentRequest, err := nsxt.generateAuthorizedRequest(http.MethodDelete, ("/policy/api/v1/infra/segments/" + name))
	if err != nil {
		return nil, fmt.Errorf("error while making the DELETE request for the segment %s: %v", name, err)
	}
	segmentResponse, err := nsxt.HttpClient.Do(segmentRequest)
	if err != nil {
		return nil, fmt.Errorf("unkown error while deleting Segment %s: %v", name, err)
	}
	if segmentResponse.StatusCode != 200 && segmentRequest.Body != nil {
		var nsxtError NSXTErrorResponse
		err := json.NewDecoder(segmentRequest.Body).Decode(&nsxtError)
		if err != nil {
			return nil, err
		}
		return &nsxtError, nil
	}
	return
}

func (nsxt *NSXTClient) CheckExistsTier1(name string) (exists bool, err error) {
	request, err := nsxt.generateAuthorizedRequest(http.MethodGet, ("/policy/api/v1/infra/tier-1s/" + name))
	if err != nil {
		return
	}
	response, err := nsxt.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode == 200 {
		exists = true
	} else if response.StatusCode == 404 {
		exists = false
	} else {
		err = errors.New("Unknown error occurred while checking if Tier 1 " + name + " exists; Status = " + response.Status)
		return false, err
	}
	return
}

func (nsxt *NSXTClient) GetTier0s() (tier0s []NSXTTier0, err error) {
	request, err := nsxt.generateAuthorizedRequest(http.MethodGet, "/policy/api/v1/infra/tier-0s/")
	if err != nil {
		return
	}
	response, err := nsxt.HttpClient.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()
	var tier0ListResult NSXTListTier0Result
	err = json.NewDecoder(response.Body).Decode(&tier0ListResult)
	if err != nil {
		return
	}
	tier0s = tier0ListResult.Results
	return
}

func (nsxt *NSXTClient) GetIpPoolSubnets(ipPoolName string) (ipSubnets []NSXTIpSubnet, err error) {
	request, err := nsxt.generateAuthorizedRequest(http.MethodGet, ("/policy/api/v1/infra/ip-pools/" + ipPoolName + "/ip-subnets"))
	if err != nil {
		return
	}
	response, err := nsxt.HttpClient.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()
	var ipSubnetsResponse NSXTListIpSubnetsResponse
	err = json.NewDecoder(response.Body).Decode(&ipSubnetsResponse)
	if err != nil {
		return
	}
	ipSubnets = ipSubnetsResponse.Results
	return
}

func (nsxt *NSXTClient) CreateNATRule(tier1Name string, sourceNetwork NSXTIPElementList, translatedNetwork NSXTIPElementList) (err error) {
	payload := NSXTNATRule{
		Description:       "NAT for CPTC Competition",
		Action:            NSXT_NAT_SNAT,
		Id:                (tier1Name + "-NAT"),
		SourceNetwork:     sourceNetwork,
		TranslatedNetwork: translatedNetwork,
	}

	jsonString, err := json.Marshal(payload)
	if err != nil {
		return
	}
	request, err := nsxt.generateAuthorizedRequestWithData(http.MethodPatch, ("/policy/api/v1/infra/tier-1s/" + tier1Name + "/nat/USER/nat-rules/" + tier1Name + "-NAT"), bytes.NewBuffer(jsonString))
	if err != nil {
		return
	}
	response, err := nsxt.HttpClient.Do(request)
	if err != nil {
		return
	}
	if response.StatusCode != 200 {
		if response.StatusCode == 409 {
			log.Infof("NAT Rule for %s already exists", tier1Name)
			return nil
		}
		err = fmt.Errorf("error %s while making NAT rules for \"%s\"", response.Status, tier1Name)
		return
	}
	return
}

func (nsxt *NSXTClient) DeleteNATRule(tier1Name string) (nsxtError *NSXTErrorResponse, err error) {
	request, err := nsxt.generateAuthorizedRequest(http.MethodDelete, ("/policy/api/v1/infra/tier-1s/" + tier1Name + "/nat/USER/nat-rules/" + tier1Name + "-NAT"))
	if err != nil {
		return nil, fmt.Errorf("error while making the DELETE request for the NAT Rule for %s: %v", tier1Name, err)
	}
	response, err := nsxt.HttpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("unkown error while deleting NAT Rule for %s: %v", tier1Name, err)
	}
	if response.StatusCode != 200 {
		var nsxtError NSXTErrorResponse
		json.NewDecoder(response.Body).Decode(&nsxtError)
		return &nsxtError, nil
	}
	return
}
