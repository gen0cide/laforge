package competition

type TFVar struct {
	Name        string
	Description string
	Value       string
}

type TFAWSProvider struct {
	Name      string
	AccessKey string
	SecretKey string
	Region    string
	Profile   string
}

type TFAWSKeyPair struct {
	Name       string
	KeyName    string
	PublicKey  string
	PrivateKey string
}

type TFAWSVirtualPrivateCloud struct {
	Name               string
	CIDR               string
	EnableDNSHostnames bool
	Tags               []string
}

type TFAWSInternetGateway struct {
	Name  string
	VPCID string
}

type TFAWSSubnet struct {
	Name                string
	VPCID               string
	CIDR                string
	AvailabilityZone    string
	MapPublicIPOnLaunch bool
	DependsOn           []string
	Tags                []string
}

type TFAWSRouteTable struct {
	Name   string
	VPCID  string
	Routes []TFAWSRoute
	Tags   []string
}

type TFAWSRoute struct {
	Name      string
	CIDR      string
	GatewayID string
}

type TFAWSRouteTableAssociation struct {
	Name         string
	SubnetID     string
	RouteTableID string
	Tags         []string
}

type TFAWSDHCPOptions struct {
	Name       string
	DomainName string
	DNSServers []string
	Tags       []string
}

type TFAWSDHCPOptionsAssociation struct {
	Name          string
	VPCID         string
	DHCPOptionsID string
}

type TFAWSIngressRule struct {
	FromPort string
	ToPort   string
	Protocol string
	CIDRs    []string
}

type TFAWSEgressRule struct {
	FromPort string
	ToPort   string
	Protocol string
	CIDRs    []string
}

type TFAWSSecurityGroup struct {
	Name         string
	Description  string
	VPCID        string
	IngressRules []TFAWSIngressRule
	EgressRules  []TFAWSEgressRule
	Tags         []string
}

type TFAWSElasticIP struct {
	Name     string
	VPC      bool
	Instance string
}

type TFOutput struct {
	Name  string
	Value string
}

type TFAWSInstance struct {
	Name                     string
	AMI                      string
	InstanceType             string
	SubnetID                 string
	KeyName                  string
	PrivateIP                string
	AssociatePublicIPAddress bool
	SecurityGroups           []TFAWSSecurityGroup
	UserDataScript           TFUserDataScript
	Connection               TFInstanceConnection
	FileProvisioners         []TFFileProvisioner
	InlineProvisioners       []TFInlineProvisioner
	ScriptProvisioners       []TFScriptProvisioner
	Tags                     []string
}

type TFUserDataScript struct {
	Name     string
	Builder  TemplateBuilder
	Template string
	Rendered string
}

type TFInstanceConnection struct {
	ConnType          string
	User              string
	Password          string
	PrivateKey        string
	BastionHost       string
	BastionUser       string
	BastionPrivateKey string
	WinRMInsecure     bool
}

type TFFileProvisioner struct {
	Name        string
	Source      string
	Content     string
	Destination string
	TFTemplate
}

type TFInlineProvisioner struct {
	Name     string
	Commands []string
	TFTemplate
}

type TFScriptProvisioner struct {
	Scripts []TFTemplate
}

type TFAWSRoute53ARecord struct {
	Name       string
	ZoneID     string
	Hostname   string
	RecordType string
	TTL        int
	Records    []string
}

type TFGCPProvider struct {
	Name           string
	CredentialFile string
	Project        string
	Region         string
}

type TFGCPNetwork struct {
	Name              string
	AutoCreateSubnets bool
}

type TFGCPSubnet struct {
	Name    string
	CIDIR   string
	Network string
}

type TFGCPFirewall struct {
	Name         string
	Network      string
	SourceRanges []string
	SourceTags   []string
	TargetTags   []string
	Rules        []TFGCPFirewallAllowRule
}

type TFGCPFirewallAllowRule struct {
	Protocol string
	Ports    []int
}

type TFGCPElasticIP struct {
	Name string
}

type TFGCPInstance struct {
	Name               string
	InstanceType       string
	Zone               string
	Tags               []string
	Image              string
	Subnet             string
	PrivateIP          string
	ElasticIP          string
	SSHPublicKey       string
	UserDataScript     TFUserDataScript
	Connection         TFInstanceConnection
	FileProvisioners   []TFFileProvisioner
	InlineProvisioners []TFInlineProvisioner
	ScriptProvisioners []TFScriptProvisioner
}
