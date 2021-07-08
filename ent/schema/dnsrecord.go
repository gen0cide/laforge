package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// DNSRecord holds the schema definition for the DNSRecord entity.
type DNSRecord struct {
	ent.Schema
}

// Fields of the DNSRecord.
func (DNSRecord) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("hcl_id").
			StructTag(`hcl:"id,label"`),
		field.String("name").
			StructTag(`hcl:"name,attr"`),
		field.JSON("values", []string{}).
			StructTag(`hcl:"values,optional"`),
		field.String("type").
			StructTag(`hcl:"type,attr"`),
		field.String("zone").
			StructTag(`hcl:"zone,attr" `),
		field.JSON("vars", map[string]string{}).
			StructTag(`hcl:"vars,optional"`),
		field.Bool("disabled").
			StructTag(`hcl:"disabled,optional"`),
		field.JSON("tags", map[string]string{}).
			StructTag(`hcl:"tags,optional"`),
	}
}

// Edges of the DNSRecord.
func (DNSRecord) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("DNSRecordToEnvironment", Environment.Type).
			Ref("EnvironmentToDNSRecord").
			Unique(),
	}
}
