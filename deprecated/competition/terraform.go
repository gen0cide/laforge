package competition

import "fmt"

type TFObject interface {
	Mix(tpl *TemplateBuilder)
}

type TFVar struct {
	Name        string
	Description string
	Value       string
}

func (t *TFVar) Mix(tpl *TemplateBuilder) {
	return
}

type TFAWSProvider struct {
	Name      string
	AccessKey string
	SecretKey string
	Region    string
	Profile   string
}

func (t *TFAWSProvider) Mix(tpl *TemplateBuilder) {
	t.AccessKey = tpl.Competition.AWS.APIKey
	t.SecretKey = tpl.Competition.AWS.APISecret
	t.Region = tpl.Competition.AWS.Region
}

type TFAWSKeyPair struct {
	Name       string
	KeyName    string
	PublicKey  string
	PrivateKey string
}

func (t *TFAWSKeyPair) Mix(tpl *TemplateBuilder) {
	keyName := fmt.Sprintf("sshkey_%s", tpl.Environment.Name)
	t.Name = keyName
	t.KeyName = keyName
	t.PublicKey = tpl.Competition.SSHPublicKey()
	t.PrivateKey = tpl.Competition.SSHPrivateKey()
}

type TFAWSVirtualPrivateCloud struct {
	Name               string
	CIDR               string
	EnableDNSHostnames bool
	Tags               []string
}

func (t *TFAWSVirtualPrivateCloud) Mix(tpl *TemplateBuilder) {
	t.Name = tpl.EnvItemName(t)
	t.CIDR = tpl.Environment.DefaultCIDR()
	t.EnableDNSHostnames = true
}

type TFAWSNATGateway struct {
	Name      string
	VPC       TFAWSVirtualPrivateCloud
	ElasticIP TFAWSElasticIP
	Subnet    TFAWSSubnet
}

func (t *TFAWSNATGateway) Mix(tpl *TemplateBuilder) {
	t.Name = tpl.NetItemName(t)
}

type TFAWSSubnet struct {
	Name                string
	VPC                 TFAWSVirtualPrivateCloud
	CIDR                string
	AvailabilityZone    string
	MapPublicIPOnLaunch bool
	DependsOn           []string
	Tags                []string
}

func (t *TFAWSSubnet) Mix(tpl *TemplateBuilder) {
	t.Name = tpl.NetItemName(t)
	t.CIDR = tpl.Network.CIDR
	t.MapPublicIPOnLaunch = false
}

type TFAWSRouteTable struct {
	Name   string
	VPC    TFAWSVirtualPrivateCloud
	Routes []TFAWSRoute
	Tags   []string
}

func (t *TFAWSRouteTable) Mix(tpl *TemplateBuilder) {
	t.Name = tpl.NetItemName(t)
}

type TFAWSRoute struct {
	Name       string
	CIDR       string
	NATGateway TFAWSNATGateway
}

func (t *TFAWSRoute) Mix(tpl *TemplateBuilder) {
	return
}

type TFAWSRouteTableAssociation struct {
	Name       string
	Subnet     TFAWSSubnet
	RouteTable TFAWSRouteTable
	Tags       []string
}

func (t *TFAWSRouteTableAssociation) Mix(tpl *TemplateBuilder) {
	t.Name = tpl.NetItemName(t)
}

type TFAWSDHCPOptions struct {
	Name       string
	DomainName string
	DNSServers []string
	Tags       []string
}

func (t *TFAWSDHCPOptions) Mix(tpl *TemplateBuilder) {
	t.Name = tpl.EnvItemName(t)
	t.DomainName = tpl.Environment.Domain
}

type TFAWSDHCPOptionsAssociation struct {
	Name        string
	VPC         TFAWSVirtualPrivateCloud
	DHCPOptions TFAWSDHCPOptions
}

func (t *TFAWSDHCPOptionsAssociation) Mix(tpl *TemplateBuilder) {
	t.Name = tpl.EnvItemName(t)
}

type TFAWSIngressRule struct {
	FromPort string
	ToPort   string
	Protocol string
	CIDRs    []string
}

func (t *TFAWSIngressRule) Mix(tpl *TemplateBuilder) {
	return
}

type TFAWSEgressRule struct {
	FromPort string
	ToPort   string
	Protocol string
	CIDRs    []string
}

func (t *TFAWSEgressRule) Mix(tpl *TemplateBuilder) {
	return
}

type TFAWSSecurityGroup struct {
	Name         string
	Description  string
	VPC          TFAWSVirtualPrivateCloud
	IngressRules []TFAWSIngressRule
	EgressRules  []TFAWSEgressRule
	Tags         []string
}

func (t *TFAWSSecurityGroup) Mix(tpl *TemplateBuilder) {
	return
}

type TFAWSElasticIP struct {
	Name     string
	VPC      bool
	Attach   bool
	Instance TFAWSInstance
}

func (t *TFAWSElasticIP) Mix(tpl *TemplateBuilder) {
	return
}

type TFOutput struct {
	Name  string
	Value string
}

func (t *TFOutput) Mix(tpl *TemplateBuilder) {
	return
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
}

func (t *TFAWSInstance) Mix(tpl *TemplateBuilder) {
	return
}

type TFUserDataScript struct {
	Name string
}

func (t *TFUserDataScript) Mix(tpl *TemplateBuilder) {
	return
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

func (t *TFInstanceConnection) Mix(tpl *TemplateBuilder) {
	return
}

type TFFileProvisioner struct {
	Name        string
	Source      string
	Content     string
	Destination string
}

func (t *TFFileProvisioner) Mix(tpl *TemplateBuilder) {
	return
}

type TFInlineProvisioner struct {
	Name     string
	Commands []string
}

func (t *TFInlineProvisioner) Mix(tpl *TemplateBuilder) {
	return
}

type TFScriptProvisioner struct {
	Scripts []TFTemplate
}

func (t *TFScriptProvisioner) Mix(tpl *TemplateBuilder) {
	return
}

type TFAWSRoute53ARecord struct {
	Name     string
	ZoneID   string
	Hostname string
	TTL      int
	Records  []string
}

func (t *TFAWSRoute53ARecord) Mix(tpl *TemplateBuilder) {
	return
}

type TFGCPProvider struct {
	Name           string
	CredentialFile string
	Project        string
	Region         string
}

func (t *TFGCPProvider) Mix(tpl *TemplateBuilder) {
	return
}

type TFGCPNetwork struct {
	Name              string
	AutoCreateSubnets bool
}

func (t *TFGCPNetwork) Mix(tpl *TemplateBuilder) {
	return
}

type TFGCPSubnet struct {
	Name    string
	CIDR    string
	Network TFGCPNetwork
}

func (t *TFGCPSubnet) Mix(tpl *TemplateBuilder) {
	return
}

type TFGCPFirewall struct {
	Name         string
	Network      TFGCPNetwork
	SourceRanges []string
	SourceTags   []string
	TargetTags   []string
	Rules        []TFGCPFirewallAllowRule
}

func (t *TFGCPFirewall) Mix(tpl *TemplateBuilder) {
	return
}

type TFGCPFirewallAllowRule struct {
	Protocol string
	Ports    []int
}

func (t *TFGCPFirewallAllowRule) Mix(tpl *TemplateBuilder) {
	return
}

type TFGCPElasticIP struct {
	Name string
}

func (t *TFGCPElasticIP) Mix(tpl *TemplateBuilder) {
	return
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
}

func (t *TFGCPInstance) Mix(tpl *TemplateBuilder) {
	return
}
