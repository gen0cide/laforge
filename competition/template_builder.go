package competition

type TemplateBuilder struct {
	Competition *Competition
	Environment *Environment
	Pod         *Pod
	Network     *Network
	Host        *Host
}

type TFTemplate struct {
	Name     string
	Builder  TemplateBuilder
	Template string
	Rendered string
}
