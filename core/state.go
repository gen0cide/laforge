package core

import (
	"fmt"

	"github.com/cespare/xxhash"
)

// Remote defines a configuration object that keeps terraform and remote files synchronized
//easyjson:json
type Remote struct {
	ID     string            `hcl:"id,label" json:"id,omitempty"`
	Type   string            `hcl:"type,attr" json:"type,omitempty"`
	Config map[string]string `hcl:"config,optional" json:"config,omitempty"`
}

// Hash implements the Hasher interface
func (r *Remote) Hash() uint64 {
	return xxhash.Sum64String(
		fmt.Sprintf(
			"id=%v type=%v config=%v",
			r.ID,
			r.Type,
			r.Config,
		),
	)
}
