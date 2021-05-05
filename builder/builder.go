package builder

type Spec struct {
	Template    string
	Hostname    string
	Description string
	DiskSize    int
}

type Builder interface {
	Init(url string) (err error)
	DeploySpec(spec Spec) (err error)
}
