package core

import "github.com/tidwall/buntdb"

// State is the primary object used to interface with the build's on disk state table
type State struct {
	DB        *buntdb.DB
	Current   *Snapshot
	Persisted *Snapshot
	Tainted   map[string]bool
}

// Assumptions:
// 1) The State will have been previously "built" and will not require checking
// existing assets
//
// 2) The state

// Calculate Delta
// Determine what if terraform has run
