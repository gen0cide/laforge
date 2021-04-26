package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Environment holds the schema definition for the Environment entity.
type Environment struct {
	ent.Schema
}

// Annotations of the Environment.
func (Environment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		edge.Annotation{
			StructTag: `json:"edges"`,
		},
	}
}

// Fields of the Environment.
func (Environment) Fields() []ent.Field {
	return []ent.Field{
		field.String("hcl_id").
			StructTag(`hcl:"id,label"`),
		field.String("competition_id").
			StructTag(`hcl:"competition_id,attr"`),
		field.String("name").
			StructTag(`hcl:"name,attr"`),
		field.String("description").
			StructTag(`hcl:"description,attr"`),
		field.String("builder").
			StructTag(`hcl:"builder,attr"`),
		field.Int("team_count").
			StructTag(`hcl:"team_count,attr"`),
		field.Int("revision").
			StructTag(`hcl:"revision,optional"`),
		field.JSON("admin_cidrs", []string{}).
			StructTag(`hcl:"admin_ranges,attr"`),
		field.JSON("exposed_vdi_ports", []string{}).
			StructTag(`hcl:"vdi_allowed_tcp_ports"`),
		field.JSON("config", map[string]string{}).
			StructTag(`hcl:"config,optional"`),
		field.JSON("tags", map[string]string{}).
			StructTag(`hcl:"tags,optional"`),
	}
}

// Edges of the Environment.
func (Environment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("EnvironmentToUser", User.Type).
			StructTag(`hcl:"maintainer,block"`),
		edge.To("EnvironmentToHost", Host.Type),
		edge.To("EnvironmentToCompetition", Competition.Type),
		edge.To("EnvironmentToIdentity", Identity.Type),
		edge.To("EnvironmentToCommand", Command.Type),
		edge.To("EnvironmentToScript", Script.Type),
		edge.To("EnvironmentToFileDownload", FileDownload.Type),
		edge.To("EnvironmentToFileDelete", FileDelete.Type),
		edge.To("EnvironmentToFileExtract", FileExtract.Type),
		edge.To("EnvironmentToIncludedNetwork", IncludedNetwork.Type).
			StructTag(`hcl:"included_network,block"`),
		edge.To("EnvironmentToFinding", Finding.Type),
		edge.To("EnvironmentToDNSRecord", DNSRecord.Type),
		edge.To("EnvironmentToDNS", DNS.Type),
		edge.To("EnvironmentToNetwork", Network.Type),
		edge.To("EnvironmentToHostDependency", HostDependency.Type),
		edge.From("EnvironmentToBuild", Build.Type).
			Ref("BuildToEnvironment"),
	}
}
