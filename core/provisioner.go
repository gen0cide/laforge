package core

// Provisioner is a meta interface to provide provisioning steps to the Builder
type Provisioner interface {
	Hasher
	Pather

	// Kind denotes the type of Provisioner this is
	Kind() string
}
