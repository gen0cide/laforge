// Code generated by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/gen0cide/laforge/ent/fileextract"
)

// FileExtract is the model entity for the FileExtract schema.
type FileExtract struct {
	config ` json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// HclID holds the value of the "hcl_id" field.
	HclID string `json:"hcl_id,omitempty" hcl:"id,label"`
	// Source holds the value of the "source" field.
	Source string `json:"source,omitempty" hcl:"source,attr"`
	// Destination holds the value of the "destination" field.
	Destination string `json:"destination,omitempty" hcl:"destination,attr"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty" hcl:"type,attr"`
	// Tags holds the value of the "tags" field.
	Tags map[string]string `json:"tags,omitempty" hcl:"tags,optional"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the FileExtractQuery when eager-loading is set.
	Edges FileExtractEdges `json:"edges"`

	// Edges put into the main struct to be loaded via hcl
	// FileExtractToTag holds the value of the FileExtractToTag edge.
	HCLFileExtractToTag []*Tag `json:"FileExtractToTag,omitempty"`
	// FileExtractToEnvironment holds the value of the FileExtractToEnvironment edge.
	HCLFileExtractToEnvironment []*Environment `json:"FileExtractToEnvironment,omitempty"`
	//
	provisioning_step_provisioning_step_to_file_extract *int
}

// FileExtractEdges holds the relations/edges for other nodes in the graph.
type FileExtractEdges struct {
	// FileExtractToTag holds the value of the FileExtractToTag edge.
	FileExtractToTag []*Tag `json:"FileExtractToTag,omitempty"`
	// FileExtractToEnvironment holds the value of the FileExtractToEnvironment edge.
	FileExtractToEnvironment []*Environment `json:"FileExtractToEnvironment,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [2]bool
}

// FileExtractToTagOrErr returns the FileExtractToTag value or an error if the edge
// was not loaded in eager-loading.
func (e FileExtractEdges) FileExtractToTagOrErr() ([]*Tag, error) {
	if e.loadedTypes[0] {
		return e.FileExtractToTag, nil
	}
	return nil, &NotLoadedError{edge: "FileExtractToTag"}
}

// FileExtractToEnvironmentOrErr returns the FileExtractToEnvironment value or an error if the edge
// was not loaded in eager-loading.
func (e FileExtractEdges) FileExtractToEnvironmentOrErr() ([]*Environment, error) {
	if e.loadedTypes[1] {
		return e.FileExtractToEnvironment, nil
	}
	return nil, &NotLoadedError{edge: "FileExtractToEnvironment"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*FileExtract) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case fileextract.FieldTags:
			values[i] = &[]byte{}
		case fileextract.FieldID:
			values[i] = &sql.NullInt64{}
		case fileextract.FieldHclID, fileextract.FieldSource, fileextract.FieldDestination, fileextract.FieldType:
			values[i] = &sql.NullString{}
		case fileextract.ForeignKeys[0]: // provisioning_step_provisioning_step_to_file_extract
			values[i] = &sql.NullInt64{}
		default:
			return nil, fmt.Errorf("unexpected column %q for type FileExtract", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the FileExtract fields.
func (fe *FileExtract) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case fileextract.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			fe.ID = int(value.Int64)
		case fileextract.FieldHclID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field hcl_id", values[i])
			} else if value.Valid {
				fe.HclID = value.String
			}
		case fileextract.FieldSource:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field source", values[i])
			} else if value.Valid {
				fe.Source = value.String
			}
		case fileextract.FieldDestination:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field destination", values[i])
			} else if value.Valid {
				fe.Destination = value.String
			}
		case fileextract.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				fe.Type = value.String
			}
		case fileextract.FieldTags:

			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field tags", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &fe.Tags); err != nil {
					return fmt.Errorf("unmarshal field tags: %v", err)
				}
			}
		case fileextract.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field provisioning_step_provisioning_step_to_file_extract", value)
			} else if value.Valid {
				fe.provisioning_step_provisioning_step_to_file_extract = new(int)
				*fe.provisioning_step_provisioning_step_to_file_extract = int(value.Int64)
			}
		}
	}
	return nil
}

// QueryFileExtractToTag queries the "FileExtractToTag" edge of the FileExtract entity.
func (fe *FileExtract) QueryFileExtractToTag() *TagQuery {
	return (&FileExtractClient{config: fe.config}).QueryFileExtractToTag(fe)
}

// QueryFileExtractToEnvironment queries the "FileExtractToEnvironment" edge of the FileExtract entity.
func (fe *FileExtract) QueryFileExtractToEnvironment() *EnvironmentQuery {
	return (&FileExtractClient{config: fe.config}).QueryFileExtractToEnvironment(fe)
}

// Update returns a builder for updating this FileExtract.
// Note that you need to call FileExtract.Unwrap() before calling this method if this FileExtract
// was returned from a transaction, and the transaction was committed or rolled back.
func (fe *FileExtract) Update() *FileExtractUpdateOne {
	return (&FileExtractClient{config: fe.config}).UpdateOne(fe)
}

// Unwrap unwraps the FileExtract entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (fe *FileExtract) Unwrap() *FileExtract {
	tx, ok := fe.config.driver.(*txDriver)
	if !ok {
		panic("ent: FileExtract is not a transactional entity")
	}
	fe.config.driver = tx.drv
	return fe
}

// String implements the fmt.Stringer.
func (fe *FileExtract) String() string {
	var builder strings.Builder
	builder.WriteString("FileExtract(")
	builder.WriteString(fmt.Sprintf("id=%v", fe.ID))
	builder.WriteString(", hcl_id=")
	builder.WriteString(fe.HclID)
	builder.WriteString(", source=")
	builder.WriteString(fe.Source)
	builder.WriteString(", destination=")
	builder.WriteString(fe.Destination)
	builder.WriteString(", type=")
	builder.WriteString(fe.Type)
	builder.WriteString(", tags=")
	builder.WriteString(fmt.Sprintf("%v", fe.Tags))
	builder.WriteByte(')')
	return builder.String()
}

// FileExtracts is a parsable slice of FileExtract.
type FileExtracts []*FileExtract

func (fe FileExtracts) config(cfg config) {
	for _i := range fe {
		fe[_i].config = cfg
	}
}