package core

import (
	"encoding/json"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
)

const (
	// DBKeySnapshot is the database key for persisting the snapshot in our local filesystem
	DBKeySnapshot = `/snapshot`
)

var (
	// ErrSnapshotsMatch is thrown when two snapshots are functionally the same during delta calculation
	ErrSnapshotsMatch = errors.New("snapshots reflected the same state")
)

// State is the primary object used to interface with the build's on disk state table
type State struct {
	DB        *buntdb.DB
	Current   *Snapshot
	Persisted *Snapshot
}

// Plan is a type that describes how to get from one state to the next
//easyjson:json
type Plan struct {
	Snapshot *Snapshot       `json:"target,omitempty"`
	Tasks    map[string]Doer `json:"-"`
	// Graph          *graph.Mutable  `json:"-"`
	OrderedTaskIDs []int           `json:"-"`
	Tasklist       []interface{}   `json:"tasks"`
	Tainted        map[string]bool `json:"tainted"`
}

// BuildTasks returns a list of tasks that need to happen
func (p *Plan) BuildTasks() error {
	// enumerate the snapshot, finding tasks and creating them
	return nil
}

// CalculateDelta attempts to determine what needs to be done to bring a base in line with target
func CalculateDelta(base *Snapshot, target *Snapshot) ([]string, error) {
	if base.Hash() == target.Hash() {
		return []string{}, ErrSnapshotsMatch
	}

	// baseGraph := build.Specific(base.Graph)
	// targetGraph := build.Specific(target.Graph)

	return []string{}, nil

}

// Open attempts to create a DB connector for the state given a local file path
func (s *State) Open(dbfile string) error {
	db, err := buntdb.Open(dbfile)
	if err != nil {
		return err
	}
	s.DB = db
	return nil
}

// LoadSnapshotFromDB attempts to load the last Snapshot object from the DB, assigning it to *State.Persisted and returning it if it was successful.
func (s *State) LoadSnapshotFromDB() (*Snapshot, error) {
	return nil, nil
}

// CreateDBSchema attempts to create the database indexes appropriately
func (s *State) CreateDBSchema() error {
	return nil
}

// PersistSnapshot will save the provided snapshot into the current snapshot entry of the database, overwriting any existing snapshot.
func (s *State) PersistSnapshot(snap *Snapshot) error {
	jsonData, err := json.Marshal(snap)
	if err != nil {
		return err
	}
	err = s.DB.Update(func(tx *buntdb.Tx) error {
		_, overwritten, err := tx.Set(DBKeySnapshot, string(jsonData), nil)
		if err != nil {
			return err
		}
		if overwritten {
			cli.Logger.Infof("Persistent Snapshot overwritten in state DB")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// TODO: Need to detect when a laforge build is required
// At least worn and allow to continue with a flag
// Things it might need to do:
// - Upload a file (RemoteFile)
// OWNERS:
//	 /envs/*/*/teams/*/networks/*/hosts/*/steps/*
// Execute a command (not logged) (SecretCommand)
// OWNERS:
//	 /envs/*/*/teams/*/networks/*/hosts/*/conn
// Run a script (logged) (upload + execute + delete) (Script)
// OWNERS:
//   /envs/*/*/teams/*/networks/*/hosts/*/steps/*
// Run terraform command sequence (terraform init, terraform refresh, terraform taint {id}, terraform apply -auto-approve -parallelism=10, terraform destroy -force -parallelism=10)
// OWNERS:
//   /envs/*/*/teams/*
//
//

// Assumptions:
// 1) The State will have been previously "built" and will not require checking
// existing assets
//
// 2) The state

// Calculate Delta
// Determine what if terraform has run
