package competition

type TFVar struct {
	Name        string
	Description string
	Value       string
	Template    TFTemplate
}

type TFAWSProvider struct {
	Name      string
	AccessKey string
	SecretKey string
	Region    string
	Profile   string
	Template  TFTemplate
}

type TFAWSKeyPair struct {
	Name       string
	KeyName    string
	PublicKey  string
	PrivateKey string
	Template   TFTemplate
}

type TFAWSVirtualPrivateCloud struct {
	Name               string
	CIDR               string
	EnableDNSHostnames bool
	Tags               []string
	Template           TFTemplate
}

type TFAWSNATGateway struct {
	Name      string
	VPC       TFAWSVirtualPrivateCloud
	ElasticIP TFAWSElasticIP
	Subnet    TFAWSSubnet
	Template  TFTemplate
}

type TFAWSSubnet struct {
	Name                string
	VPC                 TFAWSVirtualPrivateCloud
	CIDR                string
	AvailabilityZone    string
	MapPublicIPOnLaunch bool
	DependsOn           []string
	Tags                []string
	Template            TFTemplate
}

type TFAWSRouteTable struct {
	Name     string
	VPC      TFAWSVirtualPrivateCloud
	Routes   []TFAWSRoute
	Tags     []string
	Template TFTemplate
}

type TFAWSRoute struct {
	Name       string
	CIDR       string
	NATGateway TFAWSNATGateway
	Template   TFTemplate
}

type TFAWSRouteTableAssociation struct {
	Name       string
	Subnet     TFAWSSubnet
	RouteTable TFAWSRouteTable
	Tags       []string
	Template   TFTemplate
}

type TFAWSDHCPOptions struct {
	Name       string
	DomainName string
	DNSServers []string
	Tags       []string
	Template   TFTemplate
}

type TFAWSDHCPOptionsAssociation struct {
	Name        string
	VPC         TFAWSVirtualPrivateCloud
	DHCPOptions TFAWSDHCPOptions
	Template    TFTemplate
}

type TFAWSIngressRule struct {
	FromPort string
	ToPort   string
	Protocol string
	CIDRs    []string
	Template TFTemplate
}

type TFAWSEgressRule struct {
	FromPort string
	ToPort   string
	Protocol string
	CIDRs    []string
	Template TFTemplate
}

type TFAWSSecurityGroup struct {
	Name         string
	Description  string
	VPC          TFAWSVirtualPrivateCloud
	IngressRules []TFAWSIngressRule
	EgressRules  []TFAWSEgressRule
	Tags         []string
	Template     TFTemplate
}

type TFAWSElasticIP struct {
	Name     string
	VPC      bool
	Attach   bool
	Instance TFAWSInstance
	Template TFTemplate
}

type TFOutput struct {
	Name     string
	Value    string
	Template TFTemplate
}

type TFAWSInstance struct {
	Name                     string
	AMI                      string
	InstanceType             string
	Subnet                   TFAWSSubnet
	KeyPair                  TFAWSKeyPair
	PrivateIP                string
	AssociatePublicIPAddress bool
	SecurityGroups           []TFAWSSecurityGroup
	UserDataScript           TFUserDataScript
	Connection               TFInstanceConnection
	FileProvisioners         []TFFileProvisioner
	InlineProvisioners       []TFInlineProvisioner
	ScriptProvisioners       []TFScriptProvisioner
	Tags                     []string
	Template                 TFTemplate
}

type TFUserDataScript struct {
	Name     string
	Template TFTemplate
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
	Template          TFTemplate
}

type TFFileProvisioner struct {
	Name        string
	Source      string
	Content     string
	Destination string
	Template    TFTemplate
}

type TFInlineProvisioner struct {
	Name     string
	Commands []string
	Template TFTemplate
}

type TFScriptProvisioner struct {
	Scripts  []TFTemplate
	Template TFTemplate
}

type TFAWSRoute53ARecord struct {
	Name     string
	ZoneID   string
	Hostname string
	TTL      int
	Records  []string
	Template TFTemplate
}

type TFGCPProvider struct {
	Name           string
	CredentialFile string
	Project        string
	Region         string
	Template       TFTemplate
}

type TFGCPNetwork struct {
	Name              string
	AutoCreateSubnets bool
	Template          TFTemplate
}

type TFGCPSubnet struct {
	Name     string
	CIDR     string
	Network  TFGCPNetwork
	Template TFTemplate
}

type TFGCPFirewall struct {
	Name         string
	Network      TFGCPNetwork
	SourceRanges []string
	SourceTags   []string
	TargetTags   []string
	Rules        []TFGCPFirewallAllowRule
	Template     TFTemplate
}

type TFGCPFirewallAllowRule struct {
	Protocol string
	Ports    []int
	Template TFTemplate
}

type TFGCPElasticIP struct {
	Name     string
	Template TFTemplate
}

type TFGCPInstance struct {
	Name               string
	HostName           string
	InstanceType       string
	AvailabilityZone   string
	Tags               []string
	AMI                string
	Subnet             TFGCPSubnet
	PrivateIP          string
	ElasticIP          TFGCPElasticIP
	SSHPublicKey       string
	UserDataScript     TFUserDataScript
	Connection         TFInstanceConnection
	FileProvisioners   []TFFileProvisioner
	InlineProvisioners []TFInlineProvisioner
	ScriptProvisioners []TFScriptProvisioner
	Template           TFTemplate
}
