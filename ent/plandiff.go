// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/gen0cide/laforge/ent/buildcommit"
	"github.com/gen0cide/laforge/ent/plan"
	"github.com/gen0cide/laforge/ent/plandiff"
	"github.com/google/uuid"
)

// PlanDiff is the model entity for the PlanDiff schema.
type PlanDiff struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// Revision holds the value of the "revision" field.
	Revision int `json:"revision,omitempty"`
	// NewState holds the value of the "new_state" field.
	NewState plandiff.NewState `json:"new_state,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the PlanDiffQuery when eager-loading is set.
	Edges PlanDiffEdges `json:"edges"`

	// Edges put into the main struct to be loaded via hcl
	// PlanDiffToBuildCommit holds the value of the PlanDiffToBuildCommit edge.
	HCLPlanDiffToBuildCommit *BuildCommit `json:"PlanDiffToBuildCommit,omitempty"`
	// PlanDiffToPlan holds the value of the PlanDiffToPlan edge.
	HCLPlanDiffToPlan *Plan `json:"PlanDiffToPlan,omitempty"`
	//
	plan_diff_plan_diff_to_build_commit *uuid.UUID
	plan_diff_plan_diff_to_plan         *uuid.UUID
}

// PlanDiffEdges holds the relations/edges for other nodes in the graph.
type PlanDiffEdges struct {
	// PlanDiffToBuildCommit holds the value of the PlanDiffToBuildCommit edge.
	PlanDiffToBuildCommit *BuildCommit `json:"PlanDiffToBuildCommit,omitempty"`
	// PlanDiffToPlan holds the value of the PlanDiffToPlan edge.
	PlanDiffToPlan *Plan `json:"PlanDiffToPlan,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// PlanDiffToBuildCommitOrErr returns the PlanDiffToBuildCommit value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PlanDiffEdges) PlanDiffToBuildCommitOrErr() (*BuildCommit, error) {
	if e.loadedTypes[0] {
		if e.PlanDiffToBuildCommit == nil {
			// The edge PlanDiffToBuildCommit was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: buildcommit.Label}
		}
		return e.PlanDiffToBuildCommit, nil
	}
	return nil, &NotLoadedError{edge: "PlanDiffToBuildCommit"}
}

// PlanDiffToPlanOrErr returns the PlanDiffToPlan value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e PlanDiffEdges) PlanDiffToPlanOrErr() (*Plan, error) {
	if e.loadedTypes[1] {
		if e.PlanDiffToPlan == nil {
			// The edge PlanDiffToPlan was loaded in eager-loading,
			// but was not found.
			return nil, &NotFoundError{label: plan.Label}
		}
		return e.PlanDiffToPlan, nil
	}
	return nil, &NotLoadedError{edge: "PlanDiffToPlan"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*PlanDiff) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case plandiff.FieldRevision:
			values[i] = new(sql.NullInt64)
		case plandiff.FieldNewState:
			values[i] = new(sql.NullString)
		case plandiff.FieldID:
			values[i] = new(uuid.UUID)
		case plandiff.ForeignKeys[0]: // plan_diff_plan_diff_to_build_commit
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		case plandiff.ForeignKeys[1]: // plan_diff_plan_diff_to_plan
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			return nil, fmt.Errorf("unexpected column %q for type PlanDiff", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the PlanDiff fields.
func (pd *PlanDiff) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case plandiff.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				pd.ID = *value
			}
		case plandiff.FieldRevision:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field revision", values[i])
			} else if value.Valid {
				pd.Revision = int(value.Int64)
			}
		case plandiff.FieldNewState:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field new_state", values[i])
			} else if value.Valid {
				pd.NewState = plandiff.NewState(value.String)
			}
		case plandiff.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field plan_diff_plan_diff_to_build_commit", values[i])
			} else if value.Valid {
				pd.plan_diff_plan_diff_to_build_commit = new(uuid.UUID)
				*pd.plan_diff_plan_diff_to_build_commit = *value.S.(*uuid.UUID)
			}
		case plandiff.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field plan_diff_plan_diff_to_plan", values[i])
			} else if value.Valid {
				pd.plan_diff_plan_diff_to_plan = new(uuid.UUID)
				*pd.plan_diff_plan_diff_to_plan = *value.S.(*uuid.UUID)
			}
		}
	}
	return nil
}

// QueryPlanDiffToBuildCommit queries the "PlanDiffToBuildCommit" edge of the PlanDiff entity.
func (pd *PlanDiff) QueryPlanDiffToBuildCommit() *BuildCommitQuery {
	return (&PlanDiffClient{config: pd.config}).QueryPlanDiffToBuildCommit(pd)
}

// QueryPlanDiffToPlan queries the "PlanDiffToPlan" edge of the PlanDiff entity.
func (pd *PlanDiff) QueryPlanDiffToPlan() *PlanQuery {
	return (&PlanDiffClient{config: pd.config}).QueryPlanDiffToPlan(pd)
}

// Update returns a builder for updating this PlanDiff.
// Note that you need to call PlanDiff.Unwrap() before calling this method if this PlanDiff
// was returned from a transaction, and the transaction was committed or rolled back.
func (pd *PlanDiff) Update() *PlanDiffUpdateOne {
	return (&PlanDiffClient{config: pd.config}).UpdateOne(pd)
}

// Unwrap unwraps the PlanDiff entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pd *PlanDiff) Unwrap() *PlanDiff {
	tx, ok := pd.config.driver.(*txDriver)
	if !ok {
		panic("ent: PlanDiff is not a transactional entity")
	}
	pd.config.driver = tx.drv
	return pd
}

// String implements the fmt.Stringer.
func (pd *PlanDiff) String() string {
	var builder strings.Builder
	builder.WriteString("PlanDiff(")
	builder.WriteString(fmt.Sprintf("id=%v", pd.ID))
	builder.WriteString(", revision=")
	builder.WriteString(fmt.Sprintf("%v", pd.Revision))
	builder.WriteString(", new_state=")
	builder.WriteString(fmt.Sprintf("%v", pd.NewState))
	builder.WriteByte(')')
	return builder.String()
}

// PlanDiffs is a parsable slice of PlanDiff.
type PlanDiffs []*PlanDiff

func (pd PlanDiffs) config(cfg config) {
	for _i := range pd {
		pd[_i].config = cfg
	}
}