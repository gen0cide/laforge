package core

import (
	"encoding/json"
	"sort"

	"github.com/cespare/xxhash"
	"github.com/gen0cide/laforge/core/cli"
)

// ChecksumList is a type alias for a set of computed hashes
//easyjson:json
type ChecksumList []uint64

// Len implements the sort interface
func (c ChecksumList) Len() int { return len(c) }

// Swap implements the sort interface
func (c ChecksumList) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Less implements the sort interface
func (c ChecksumList) Less(i, j int) bool { return c[i] < c[j] }

// Hash implements the hasher interface
func (c ChecksumList) Hash() uint64 {
	sort.Sort(c)
	d, err := json.Marshal(c)
	if err != nil {
		cli.Logger.Errorf("unable to generate hash for a checksum list")
		return uint64(666)
	}
	return xxhash.Sum64(d)
}
