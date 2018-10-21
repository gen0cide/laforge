package core

import (
	"encoding/binary"
	"sort"

	"github.com/cespare/xxhash"
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
	buf := make([]byte, binary.Size(c)*64)
	for _, x := range c {
		binary.PutUvarint(buf, x)
	}
	return xxhash.Sum64(buf)
}
