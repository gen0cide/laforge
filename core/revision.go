package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"time"

	"path/filepath"
)

const (
	// RevStatusUnknown is when the resource can be assumed to be destoryed or MIA
	RevStatusUnknown RevStatus = `UNKNOWN`

	// RevStatusStale is when the resource is no longer used by outside systems (non-existent terraform host)
	RevStatusStale RevStatus = `STALE`

	// RevStatusActive is when the resource has actively been created successfully
	RevStatusActive RevStatus = `ACTIVE`

	// RevStatusPlanned is when the resource doesnt require a config change, but has not been created yet.
	RevStatusPlanned RevStatus = `PLANNED`

	// RevStatusFailed is when the resource was attempted to have been created, but failed in the process.
	RevStatusFailed RevStatus = `FAILED`

	// RevModTouch describes a revision that needs to be updated on disk
	RevModTouch RevMod = `TOUCH`

	// RevModDelete describes a revision that needs to be deleted off disk
	RevModDelete RevMod = `DELETE`

	// RevModCreate describes a revision that needs to be created on disk
	RevModCreate RevMod = `CREATE`

	// RevModRebuild describes a revision that needs to be rebuilt due to human declaration
	RevModRebuild RevMod = `REBUILD`
)

// RevStatus is a type used to describe the current state of the revision
type RevStatus string

// RevMod is an internal type alias to label needs of objects within an environments deployment
type RevMod string

// Revision is used to describe a small .lfrevision file placed in the root of each path
//easyjson:json
type Revision struct {
	ID         string            `json:"id"`
	Type       LFType            `json:"type"`
	Status     RevStatus         `json:"status"`
	Checksum   uint64            `json:"checksum"`
	Timestamp  time.Time         `json:"timestamp"`
	ExternalID string            `json:"external_id"`
	Vars       map[string]string `json:"vars"`
}

// Touch sets the current timestamp and status to active for use within templating engines
func (r *Revision) Touch() *Revision {
	r.Status = RevStatusActive
	r.Timestamp = time.Now()
	return r
}

// TouchWithID touches the revision and updates it's External ID resource
func (r *Revision) TouchWithID(s string) *Revision {
	r.Touch()
	r.ExternalID = s
	return r
}

// Taint changes the revision to one that is a stale state
func (r *Revision) Taint() *Revision {
	r.Status = RevStatusStale
	r.Timestamp = time.Now()
	r.Checksum = 666
	return r
}

// AbsPath returns a joined file path for build types and below
func (r *Revision) AbsPath(basedir string) string {
	if r.Type == LFTypeEnvironment {
		return fmt.Sprintf(".env.lfrevision")
	}
	if r.Type == LFTypeConnection {
		return filepath.Join(basedir, filepath.Dir(r.ID), r.Filename())
	}
	return filepath.Join(basedir, r.ID, r.Filename())
}

// Filename returns the base filename of the revision file
func (r *Revision) Filename() string {
	if r.Type == LFTypeProvisioningStep {
		return fmt.Sprintf(".%s.pstep.lfrevision", path.Base(r.ID))
	}
	return fmt.Sprintf(".%s.lfrevision", string(r.Type))
}

// ParseRevisionFile attempts to parse a revision file at the given location
func ParseRevisionFile(fpath string) (*Revision, error) {
	//nolint:gosec
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	var rev Revision
	err = json.Unmarshal(data, &rev)
	if err != nil {
		return nil, err
	}

	return &rev, nil
}

// ToJSONString converts the revision to a JSON string
func (r *Revision) ToJSONString() string {
	data, _ := json.Marshal(r)
	return string(data)
}
