// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"github.com/facebook/ent/dialect/sql"
	"github.com/gen0cide/laforge/ent/finding"
)

// Finding is the model entity for the Finding schema.
type Finding struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Severity holds the value of the "severity" field.
	Severity finding.Severity `json:"severity,omitempty"`
	// Difficulty holds the value of the "difficulty" field.
	Difficulty finding.Difficulty `json:"difficulty,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the FindingQuery when eager-loading is set.
	Edges FindingEdges `json:"edges"`
}

// FindingEdges holds the relations/edges for other nodes in the graph.
type FindingEdges struct {
	// FindingToUser holds the value of the FindingToUser edge.
	FindingToUser []*User
	// FindingToTag holds the value of the FindingToTag edge.
	FindingToTag []*Tag
	// FindingToHost holds the value of the FindingToHost edge.
	FindingToHost []*Host
	// FindingToScript holds the value of the FindingToScript edge.
	FindingToScript []*Script
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// FindingToUserOrErr returns the FindingToUser value or an error if the edge
// was not loaded in eager-loading.
func (e FindingEdges) FindingToUserOrErr() ([]*User, error) {
	if e.loadedTypes[0] {
		return e.FindingToUser, nil
	}
	return nil, &NotLoadedError{edge: "FindingToUser"}
}

// FindingToTagOrErr returns the FindingToTag value or an error if the edge
// was not loaded in eager-loading.
func (e FindingEdges) FindingToTagOrErr() ([]*Tag, error) {
	if e.loadedTypes[1] {
		return e.FindingToTag, nil
	}
	return nil, &NotLoadedError{edge: "FindingToTag"}
}

// FindingToHostOrErr returns the FindingToHost value or an error if the edge
// was not loaded in eager-loading.
func (e FindingEdges) FindingToHostOrErr() ([]*Host, error) {
	if e.loadedTypes[2] {
		return e.FindingToHost, nil
	}
	return nil, &NotLoadedError{edge: "FindingToHost"}
}

// FindingToScriptOrErr returns the FindingToScript value or an error if the edge
// was not loaded in eager-loading.
func (e FindingEdges) FindingToScriptOrErr() ([]*Script, error) {
	if e.loadedTypes[3] {
		return e.FindingToScript, nil
	}
	return nil, &NotLoadedError{edge: "FindingToScript"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Finding) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullString{}, // name
		&sql.NullString{}, // description
		&sql.NullString{}, // severity
		&sql.NullString{}, // difficulty
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Finding fields.
func (f *Finding) assignValues(values ...interface{}) error {
	if m, n := len(values), len(finding.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	f.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[0])
	} else if value.Valid {
		f.Name = value.String
	}
	if value, ok := values[1].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field description", values[1])
	} else if value.Valid {
		f.Description = value.String
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field severity", values[2])
	} else if value.Valid {
		f.Severity = finding.Severity(value.String)
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field difficulty", values[3])
	} else if value.Valid {
		f.Difficulty = finding.Difficulty(value.String)
	}
	return nil
}

// QueryFindingToUser queries the FindingToUser edge of the Finding.
func (f *Finding) QueryFindingToUser() *UserQuery {
	return (&FindingClient{config: f.config}).QueryFindingToUser(f)
}

// QueryFindingToTag queries the FindingToTag edge of the Finding.
func (f *Finding) QueryFindingToTag() *TagQuery {
	return (&FindingClient{config: f.config}).QueryFindingToTag(f)
}

// QueryFindingToHost queries the FindingToHost edge of the Finding.
func (f *Finding) QueryFindingToHost() *HostQuery {
	return (&FindingClient{config: f.config}).QueryFindingToHost(f)
}

// QueryFindingToScript queries the FindingToScript edge of the Finding.
func (f *Finding) QueryFindingToScript() *ScriptQuery {
	return (&FindingClient{config: f.config}).QueryFindingToScript(f)
}

// Update returns a builder for updating this Finding.
// Note that, you need to call Finding.Unwrap() before calling this method, if this Finding
// was returned from a transaction, and the transaction was committed or rolled back.
func (f *Finding) Update() *FindingUpdateOne {
	return (&FindingClient{config: f.config}).UpdateOne(f)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (f *Finding) Unwrap() *Finding {
	tx, ok := f.config.driver.(*txDriver)
	if !ok {
		panic("ent: Finding is not a transactional entity")
	}
	f.config.driver = tx.drv
	return f
}

// String implements the fmt.Stringer.
func (f *Finding) String() string {
	var builder strings.Builder
	builder.WriteString("Finding(")
	builder.WriteString(fmt.Sprintf("id=%v", f.ID))
	builder.WriteString(", name=")
	builder.WriteString(f.Name)
	builder.WriteString(", description=")
	builder.WriteString(f.Description)
	builder.WriteString(", severity=")
	builder.WriteString(fmt.Sprintf("%v", f.Severity))
	builder.WriteString(", difficulty=")
	builder.WriteString(fmt.Sprintf("%v", f.Difficulty))
	builder.WriteByte(')')
	return builder.String()
}

// Findings is a parsable slice of Finding.
type Findings []*Finding

func (f Findings) config(cfg config) {
	for _i := range f {
		f[_i].config = cfg
	}
}
