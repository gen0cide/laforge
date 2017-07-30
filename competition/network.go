package competition

type Network struct {
	CIDR     string `yaml:"cidr"`
	Name     string `yaml:"name"`
	DNSZone  string `yaml:"subdomain"`
	Provider string `yaml:"provider"`
	Vars     `yaml:"variables"`
	Hosts    []Host `yaml:"hosts"`
	Pod      `yaml:"-"`
}
